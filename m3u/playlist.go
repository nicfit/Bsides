package m3u

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/nicfit/bsides"
)

// Playlist extends bsides.Playlist
type Playlist interface {
	bsides.Playlist
	Name() string
	SetName(name string)
	Path() string
	IsExtended() bool
	SetExtended(extended bool)
	Export(opts ...ExportOption) string
}

func NewPlaylist() Playlist {
	m3u := &m3u{
		Playlist: bsides.NewPlaylist(),
	}
	return m3u
}

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
	var (
		b strings.Builder
	)

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
			if track.Duration() > 0 && track.Title() != "" && includeExtension(EXT_INF) {
				b.WriteString(fmt.Sprintf("%s%d,%s\n", EXT_INF, int(track.Duration().Seconds()), track.DisplayName()))
			}
		}

		b.WriteString(track.URI() + "\n")
	}

	return b.String()
}

func WithExtension(e string, exts ...string) ExportOption {
	return func(opts *exportOptions) {
		if opts.includedExtensions == nil {
			opts.includedExtensions = make([]string, 0, len(exts)+1)
		}
		opts.includedExtensions = append(opts.includedExtensions, e)
		if len(exts) > 0 {
			opts.includedExtensions = append(opts.includedExtensions, exts...)
		}
	}
}

type ExportOption func(*exportOptions)

type exportOptions struct {
	includedExtensions []string
}

func newExportOptions() *exportOptions {
	return &exportOptions{
		// nil means all
		includedExtensions: nil,
	}
}
