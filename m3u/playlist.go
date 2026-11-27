package m3u

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/nicfit/bsides"
	eyed3 "github.com/nicfit/go-eyeD3"
)

// An M3U implementation of Playlist
type m3u struct {
	bsides.Playlist

	path     string
	name     string
	extended bool
}

func (m *m3u) Name() string     { return m.name }
func (m *m3u) Path() string     { return m.path }
func (m *m3u) IsExtended() bool { return m.extended }

func (m *m3u) SetName(name string) {
	m.name = name
}

func (m *m3u) SetExtended(extended bool) {
	m.extended = extended
}

func (m *m3u) Add(t bsides.Item) (bsides.ID, error) {
	if _, ok := t.(Track); !ok {
		return bsides.NilID, errors.New("m3u.Playlist requires m3u.Track")
	}
	return m.Playlist.Add(t)
}

func (m *m3u) Insert(t bsides.Item, id bsides.ID) (bsides.ID, bool, error) {
	if _, ok := t.(Track); !ok {
		return bsides.NilID, false, errors.New("m3u.Playlist requires m3u.Track")
	}
	return m.Playlist.Insert(t, id)
}

func (m *m3u) Export(opt_args ...ExportOption) string {
	var b strings.Builder

	opts := newExportOptions()
	for _, o := range opt_args {
		o(opts)
	}

	b.Grow(512)

	if m.extended {
		b.WriteString(EXT_M3U + "\n")
		if m.name != "" {
			b.WriteString(EXT_PLAYLIST + m.name + "\n")
		}
	}

	includeExtension := func(ext string) bool {
		if opts.excludedExtensions != nil && slices.Contains(opts.excludedExtensions, ext) {
			return false
		}
		if opts.includedExtensions == nil {
			return true
		} else if slices.Contains(opts.includedExtensions, ext) {
			return true
		}
		return false
	}
	for _, entry := range m.Entries() {
		track := entry.Item.(Track)

		if m.extended {
			if track.Group() != "" && includeExtension(EXT_GROUP) {
				b.WriteString(EXT_GROUP + track.Group() + "\n")
			}
			if track.Artist() != "" && includeExtension(EXT_ARTIST) {
				b.WriteString(EXT_ARTIST + track.Artist() + "\n")
			}
			if track.Album() != "" && includeExtension(EXT_ALBUM) {
				b.WriteString(EXT_ALBUM + track.Album() + "\n")
			}
			if track.AlbumArtURL() != "" && includeExtension(EXT_ALBUM_ART_URL) {
				b.WriteString(EXT_ALBUM_ART_URL + track.AlbumArtURL() + "\n")
			}
			if track.Genre() != "" && includeExtension(EXT_GENRE) {
				b.WriteString(EXT_GENRE + track.Genre() + "\n")
			}
			if track.Size() != 0 && includeExtension(EXT_BYTES) {
				b.WriteString(EXT_BYTES + strconv.FormatInt(track.Size(), 10) + "\n")
			}
			if track.Duration() > 0 && track.DisplayName() != "" && includeExtension(EXT_INF) {
				display_name := track.DisplayName()
				if opts.extInfFormatter != nil {
					display_name = opts.extInfFormatter(track)
				}
				_, _ = fmt.Fprintf(&b, "%s%d,%s\n", EXT_INF, int(track.Duration().Seconds()), display_name)
			}
		}

		b.WriteString(track.URI() + "\n")
	}

	return b.String()
}

func WithExtension(e string, extra ...string) ExportOption {
	return func(opts *exportOptions) {
		if opts.includedExtensions == nil {
			// Nil/all mean include all, so init must happen on first. Not the case for exclude list.
			opts.includedExtensions = make([]string, 0, len(extra)+1)
		}
		opts.includedExtensions = append(opts.includedExtensions, extOptionsPrep(e, extra...)...)
	}
}

func WithoutExtension(e string, extra ...string) ExportOption {
	return func(opts *exportOptions) {
		opts.excludedExtensions = append(opts.excludedExtensions, extOptionsPrep(e, extra...)...)
	}
}

func extOptionsPrep(e string, extra ...string) []string {
	extensions := append([]string{e}, extra...)
	for i, ext := range extensions {
		ext = strings.ToUpper(strings.TrimSpace(ext))
		if !strings.HasPrefix(ext, "#") {
			ext = "#" + ext
		}
		if !strings.HasSuffix(ext, ":") {
			ext = ext + ":"
		}
		extensions[i] = ext
	}
	return extensions
}

func WithDisplayNameFormatter(f func(t Track) string) ExportOption {
	return func(opts *exportOptions) {
		opts.extInfFormatter = f
	}
}

func (m *m3u) Import(root string, opt_args ...ImportOption) error {
	opts := newImportOptions()
	for _, o := range opt_args {
		o(opts)
	}

	if err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if !opts.recursive && d.IsDir() && root != path {
			return filepath.SkipDir

		}
		if err != nil {
			return err
		}

		if bsides.IsAudioFile(path) {
			track := NewTrack(path)
			if m.IsExtended() && opts.eyeD3Metadata {
				if metadata, err := eyeD3(path); err != nil {
					return fmt.Errorf("error getting metadata, continuing without: %v", err)
				} else {
					track.Init(metadata)
				}
			}
			if _, err := m.Add(track); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return fmt.Errorf("error walking path: %v", err)
	}

	return nil
}

func RecursiveImport() ImportOption {
	return func(opts *importOptions) {
		opts.recursive = true
	}
}

func WithEyeD3Metadata() ImportOption {
	return func(opts *importOptions) {
		opts.eyeD3Metadata = true
	}
}

func eyeD3(path string) (eyed3.Track, error) {
	tracks, err := eyed3.EyeD3(path)
	if err != nil {
		return eyed3.Track{}, err
	}
	return *tracks[0], nil
}

type ExportOption func(*exportOptions)

type exportOptions struct {
	includedExtensions []string
	excludedExtensions []string
	extInfFormatter    func(t Track) string
}

func newExportOptions() *exportOptions {
	return &exportOptions{
		// nil means all
		includedExtensions: nil,
		excludedExtensions: make([]string, 0),
	}
}

type ImportOption func(*importOptions)

type importOptions struct {
	eyeD3Metadata bool
	recursive     bool
}

func newImportOptions() *importOptions {
	return &importOptions{
		eyeD3Metadata: false,
	}
}

const (
	TYPE_ID           = "M3U"
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

var ExtensionTags = []string{
	EXT_INF, EXT_PLAYLIST, EXT_GROUP, EXT_ALBUM, EXT_ARTIST, EXT_GENRE, EXT_BYTES, EXT_ALBUM_ART_URL,
}

func loadData(fp io.Reader) (_ *m3u, err error) {
	m3u := NewPlaylist().(*m3u)

	var curr *entry = nil
	addEntry := func(uri string) error {
		if uri == "" || strings.HasPrefix(uri, "#") {
			return nil
		}

		curr.uri = uri
		if _, err = m3u.Add(curr); err != nil {
			return err
		}
		curr = nil
		return nil
	}

	line_num := 0
	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		line_num++
		if curr == nil {
			curr = &entry{}
		}

		// Extended M3U check
		if line_num == 1 {
			if line == EXT_M3U {
				m3u.extended = true
			} else {
				if err := addEntry(line); err != nil {
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

		if err := addEntry(line); err != nil {
			return nil, err
		}
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
					var inf string
					d, inf = strings.TrimSpace(d), strings.TrimSpace(n)
					if d != "" {
						var err error
						if entry.duration, err = time.ParseDuration(fmt.Sprintf("%ss", d)); err != nil {
							return fmt.Errorf("invalid duration: %s", d)
						}
					}

					// There is no format for EXT_INF after duration, though "Artist - Title" a common convention.
					// Also, there is not dedicated title EXT tag, so....
					if art, tit, ok := strings.Cut(inf, " - "); ok {
						entry.artist = strings.TrimSpace(art)
						entry.title = strings.TrimSpace(tit)
						entry.displayName = ""
					} else {
						entry.displayName = strings.TrimSpace(inf)
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
