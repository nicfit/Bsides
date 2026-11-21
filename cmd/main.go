package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/nicfit/bsides"
)

func Main(args []string, stdout, stderr io.Writer) error {
	if len(args) != 1 {
		return errors.New("usage: bsides <playlist>")
	}

	var playlist bsides.Playlist
	fmt.Fprintln(stdout, playlist)
	return nil
}

func main() {
	if err := Main(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		log.Fatal(err)
	}
}
