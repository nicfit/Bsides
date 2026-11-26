package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/nicfit/bsides"
	"github.com/nicfit/bsides/m3u"
	"github.com/urfave/cli/v3"
)

func LoadPlaylist(path string) (bsides.Playlist, error) {
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

func _main(args []string, stdout, stderr io.Writer) error {
	cmd := newMainCommand()
	cmd.Writer, cmd.ErrWriter = stdout, stderr
	return cmd.Run(context.Background(), args)
}

func main() {
	if err := _main(os.Args, os.Stdout, os.Stderr); err != nil {
		log.Fatal(err)
	}
}

func newMainCommand() *cli.Command {
	return &cli.Command{
		Name:    "Bsides",
		Usage:   "Playlist utils...",
		Aliases: []string{"bsides"},
		Version: "0.1.0",
		UsageText: `Bsides [OPTIONS] PLAYLIST
Bsides [OPTIONS] CMD
`,
		HideHelpCommand:            true,
		EnableShellCompletion:      true,
		ShellCompletionCommandName: "completion",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			playlist, err := LoadPlaylist(cmd.StringArg("playlist"))
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.Writer, playlist.(m3u.Playlist).Export())
			return nil
		},
		Arguments: []cli.Argument{
			&cli.StringArg{Name: "playlist", Value: "", Required: true},
		},
		Commands: []*cli.Command{
			{
				Name:  "new",
				Usage: "create new playlists",
				Commands: []*cli.Command{
					newM3uCommand(),
				},
			},
			newIterCommand(),
		},
	}
}
