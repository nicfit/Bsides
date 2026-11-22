package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/nicfit/bsides"
	"github.com/nicfit/bsides/m3u"
)

func LoadPlaylist(path string) (bsides.Playlist, error) {
	var (
		pl  bsides.Playlist
		err error
	)

	ext := filepath.Ext(path)
	switch ext {
	case ".m3u", ".m3u8":
		pl, err = m3u.Load(os.Args[1])
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

func Main(args []string, stdout, stderr io.Writer) error {
	if len(args) != 1 {
		return errors.New("usage: bsides <playlist>")
	}

	playlist, err := LoadPlaylist(args[0])
	if err != nil {
		return err
	}

	fmt.Fprintln(stdout, playlist.(m3u.Playlist).Export())

	return nil
}

func main() {
	if err := Main(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		log.Fatal(err)
	}
}
