package factory

import (
	"fmt"
	"path/filepath"

	"github.com/nicfit/bsides"
	"github.com/nicfit/bsides/m3u"
)

func Load(path string) (bsides.Playlist, error) {
	var (
		pl  bsides.Playlist
		err error
	)

	ext := filepath.Ext(path)
	switch ext {
	case m3u.DOT_M3U, m3u.DOT_M3U8:
		pl, err = m3u.Load(path)
	case ".pls":
		return nil, fmt.Errorf("unsupported file extension: %s", ext)
	case ".xspf", ".xml":
		return nil, fmt.Errorf("unsupported file extension: %s", ext)
	case ".jspf", ".json":
		return nil, fmt.Errorf("unsupported file extension: %s", ext)
	default:
		return nil, fmt.Errorf("invalid file extension: %s", ext)
	}

	return pl, err
}
