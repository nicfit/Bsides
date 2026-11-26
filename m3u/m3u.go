package m3u

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

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
	Import(path string, opts ...ImportOption) error
}

func (m *m3u) Metadata() bsides.PlaylistMetadata {
	return bsides.PlaylistMetadata{
		TypeID:     TYPE_ID,
		Extensions: []string{DOT_M3U, DOT_M3U8},
		Mimetypes:  []string{"application/mpegurl", "application/x-mpegurl", "audio/mpegurl", "audio/x-mpeg"},
	}
}

func NewPlaylist() Playlist {
	m3u := &m3u{
		Playlist: bsides.NewPlaylist(),
	}
	return m3u
}

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
