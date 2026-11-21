package m3u

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/nicfit/bsides"
)

// Load a .m3u or .m3u8 playlist from path. Per the "standard", .m3u uses the system's local encoding,
// and .m3u8 is utf-8.
func Load(path string) (bsides.Playlist, error) {
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

	// FIXME: WIP
	var playlist bsides.Playlist
	_ = reader
	return playlist, nil
}

const (
	DOT_M3U  = ".m3u"
	DOT_M3U8 = ".m3u8"
)
