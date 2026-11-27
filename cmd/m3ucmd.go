package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/nicfit/bsides/m3u"
	"github.com/urfave/cli/v3"
)

func newM3uCommand() *cli.Command {
	type m3uFlags struct {
		name       string
		extended   bool
		withExt    []string
		withoutExt []string
		noEyeD3    bool
		recursive  bool
	}
	flags := &m3uFlags{}

	new_cmd := &cli.Command{
		Name:  "m3u",
		Usage: "Create a new playlist",
		Description: fmt.Sprintf(`Extended M3U tags:
%s`, strings.Join(m3u.ExtensionTags, ", ")),
		Arguments: []cli.Argument{
			&cli.StringArg{Name: "directory", Value: "", Required: true},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "name",
				Aliases:     []string{"n"},
				Usage:       "Playlist name",
				Destination: &flags.name,
			},
			&cli.BoolFlag{
				Name:        "extended",
				Aliases:     []string{"x"},
				Usage:       "Extended m3u",
				Value:       false,
				Destination: &flags.extended,
			},
			&cli.BoolFlag{
				Name:        "recursive",
				Aliases:     []string{"r"},
				Usage:       "Recurse into all subdirectories",
				Value:       false,
				Destination: &flags.recursive,
			},

			&cli.StringSliceFlag{
				Name:        "with-ext",
				Usage:       "Include extended attribute",
				Value:       nil,
				Destination: &flags.withExt,
			},
			&cli.StringSliceFlag{
				Name:        "without-ext",
				Usage:       "Exclude extended attribute",
				Value:       nil,
				Destination: &flags.withoutExt,
			},
			&cli.BoolFlag{
				Name:        "no-eyeD3",
				Usage:       "Do not use eyedD3 to read metadata",
				Value:       false,
				Destination: &flags.noEyeD3,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			playlist := m3u.NewPlaylist()
			playlist.SetName(flags.name)
			playlist.SetExtended(flags.extended)

			import_opts := make([]m3u.ImportOption, 0)
			if flags.recursive {
				import_opts = append(import_opts, m3u.RecursiveImport())
			}
			if !flags.noEyeD3 {
				import_opts = append(import_opts, m3u.WithEyeD3Metadata())
			}
			if err := playlist.Import(cmd.StringArg("directory"), import_opts...); err != nil {
				return err
			}

			export_opts := make([]m3u.ExportOption, 0)
			for _, ext := range flags.withExt {
				export_opts = append(export_opts, m3u.WithExtension(ext))
			}
			for _, ext := range flags.withoutExt {
				export_opts = append(export_opts, m3u.WithoutExtension(ext))
			}

			_, _ = fmt.Fprintln(cmd.Writer, playlist.Export(export_opts...))

			return nil
		},
	}
	return new_cmd
}
