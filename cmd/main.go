package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/nicfit/bsides/factory"
	"github.com/nicfit/bsides/m3u"
	"github.com/urfave/cli/v3"
)

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
			playlist, err := factory.Load(cmd.StringArg("playlist"))
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
