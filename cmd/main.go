package main

import (
	"fmt"
	"os"

	"github.com/nicfit/bsides"
)

func Main() error {
	var (
		pl bsides.Playlist
	)
	fmt.Println(pl)

	/*
		if len(os.Args) < 2 {
			return fmt.Errorf("usage: %s <playlist-file>", os.Args[0])
		}
		playlist, err := m3u.Load(os.Args[1])
		if err != nil {
			return err
		}
		fmt.Println(playlist)

		for _, entry := range playlist.Entries() {
			fmt.Println(entry)
		}
		fmt.Println("---------------------------------------------------------")
		playlist.SetName("Extended Exported Playlist")
		playlist.SetExtended(true)
		fmt.Println(playlist.Export())
		fmt.Println("---------------------------------------------------------")
		playlist.SetName("Extended Exported Playlist (select extensions)")
		fmt.Println(playlist.Export(m3u.WithExtension(m3u.EXT_PLAYLIST, m3u.EXT_INF)))
		fmt.Println("---------------------------------------------------------")
		playlist.SetExtended(false)
		fmt.Println(playlist.Export())

		fmt.Println("std iterator: ----------------------------------------")
		it := core.StandardIterator(playlist)
		for entry := it.Next(); !entry.IsZero(); entry = it.Next() {
			fmt.Println(entry)
		}
	*/

	return nil
}

func main() {
	if err := Main(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
