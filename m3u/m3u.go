package m3u

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Load a .m3u or .m3u8 playlist from path. Per the "standard", .m3u uses the system's local encoding,
// and .m3u8 is utf-8.
func Load(path string) (Playlist, error) {
	fp, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = fp.Close() }()

	var reader io.Reader
	switch filepath.Ext(path) {
	case DOT_M3U:
		enc, err := getLocalEncoding()
		if err != nil {
			return nil, fmt.Errorf("error getting local encoding for m3u: %w", err)
		} else if enc == nil {
			return nil, errors.New("local encoding could not be determined for m3u")
		}
		reader = enc.NewDecoder().Reader(fp)
	case DOT_M3U8:
		reader = fp
	default:
		return nil, fmt.Errorf("unsupported file type: %s", filepath.Ext(path))
	}

	if m3u, err := loadData(reader); err == nil {
		m3u.path = path
		return m3u, nil
	} else {
		return nil, err
	}
}

func LoadData(fp io.Reader) (Playlist, error) {
	return loadData(fp)
}

const (
	DOT_M3U           = ".m3u"
	DOT_M3U8          = ".m3u8"
	EXT_M3U           = "#EXTM3U"
	EXT_INF           = "#EXTINF:"
	EXT_PLAYLIST      = "#PLAYLIST:"
	EXT_GROUP         = "#EXTGRP:"
	EXT_ALBUM         = "#EXTALB:"
	EXT_ARTIST        = "#EXTART:"
	EXT_GENRE         = "#EXTGENRE:"
	EXT_BYTES         = "#EXTBYT:"
	EXT_ALBUM_ART_URL = "#EXTALBUMARTURL:"
	// EXT_M3A is not supported by the current implementation
	EXT_M3A = "#EXTM3A"
	// EXT_BIN is not supported by the current implementation
	EXT_BIN = "#EXTBIN:"
	// EXT_VLC_OPT is not supported by the current implementation
	EXT_VLC_OPT = "#EXTVLCOPT:"
)

func loadData(fp io.Reader) (_ *m3u, err error) {
	m3u := NewPlaylist().(*m3u)

	addEntry := func(e *entry, uri string) error {
		if uri == "" || strings.HasPrefix(uri, "#") {
			return nil
		}

		e.uri = uri
		if _, err = m3u.Add(e); err != nil {
			return err
		}
		return nil
	}
	var curr *entry = nil
	line_num := 0
	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		line_num++

		// Extended M3U check
		if line_num == 1 {
			curr = &entry{}
			if line == EXT_M3U {
				m3u.extended = true
			} else {
				if err := addEntry(curr, line); err != nil {
					return nil, err
				}
			}
			continue
		}

		if m3u.extended && strings.HasPrefix(line, "#") {
			if processed, err := processM3uExtension(line, m3u, curr); err != nil {
				return nil, err
			} else if processed {
				continue
			}
		}

		if err := addEntry(curr, line); err != nil {
			return nil, err
		}
		curr = &entry{}
	}

	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return m3u, nil
}

func extSplit(line string) (string, string) {
	if before, after, ok := strings.Cut(line, ":"); ok {
		return strings.TrimSpace(before), strings.TrimSpace(after)
	}
	return line, ""
}

func processM3uExtension(line string, m3u *m3u, entry *entry) (bool, error) {
	// Playlist name
	if strings.HasPrefix(line, EXT_PLAYLIST) {
		_, name := extSplit(line)
		if m3u.name == "" {
			m3u.name = name
		} else if m3u.name != name {
			return false, fmt.Errorf("multiple playlist names found: %s and %s", m3u.name, name)
		}
		return true, nil
	}

	// The following extensions are collected and associated with the next entry
	type extSetter struct {
		ext string
		set func(string) error
	}
	setters := []extSetter{
		{ext: EXT_GROUP, set: func(val string) error { _, entry.group = extSplit(line); return nil }},
		{ext: EXT_ARTIST, set: func(val string) error { _, entry.artist = extSplit(line); return nil }},
		{ext: EXT_ALBUM, set: func(val string) error { _, entry.album = extSplit(line); return nil }},
		{ext: EXT_GENRE, set: func(val string) error { _, entry.genre = extSplit(line); return nil }},
		{ext: EXT_ALBUM_ART_URL, set: func(val string) error { _, entry.albumArtUrl = extSplit(line); return nil }},
		{ext: EXT_BYTES, set: func(val string) error {
			if _, s := extSplit(line); s != "" {
				var err error
				if entry.size, err = strconv.ParseInt(s, 10, 64); err != nil {
					return fmt.Errorf("invalid size: %s", s)
				}
			}
			return nil
		}},
		{ext: EXT_INF, set: func(val string) error {
			if _, s := extSplit(line); s != "" {
				if d, n, ok := strings.Cut(s, ","); ok {
					d, entry.displayName = strings.TrimSpace(d), strings.TrimSpace(n)
					if d != "" {
						var err error
						if entry.duration, err = time.ParseDuration(fmt.Sprintf("%ss", d)); err != nil {
							return fmt.Errorf("invalid duration: %s", d)
						}
					}
				}
			}
			return nil
		}},
	}

	for _, setter := range setters {
		if strings.HasPrefix(line, setter.ext) {
			if err := setter.set(line); err != nil {
				return false, fmt.Errorf("error setting extension %s: %s", setter.ext, err)
			}
			return true, nil
		}
	}

	return false, nil
}
