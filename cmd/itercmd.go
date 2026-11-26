package main

import (
	"context"
	"fmt"
	"time"

	"github.com/nicfit/bsides"
	"github.com/nicfit/bsides/iters"
	"github.com/urfave/cli/v3"
)

func newIterCommand() *cli.Command {
	type Flags struct {
		shuffle  bool
		reverse  bool
		repeat   int
		interval time.Duration
	}
	flags := &Flags{}

	root := &cli.Command{
		Name:  "iter",
		Usage: "Iterate a playlist",
		Arguments: []cli.Argument{
			&cli.StringArg{Name: "playlist", Value: "", Required: true},
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "shuffle",
				Usage:       "Iterate in shuffled order",
				Value:       false,
				Destination: &flags.shuffle,
			},
			&cli.BoolFlag{
				Name:        "reverse",
				Usage:       "Iterate in reverse order",
				Value:       false,
				Destination: &flags.reverse,
			},
			&cli.IntFlag{
				Name:        "repeat",
				Usage:       "Repeat N time when the iteration is done",
				Value:       0,
				Destination: &flags.repeat,
			},
			&cli.DurationFlag{
				Name:        "interval",
				Usage:       "Time between iterations",
				Value:       time.Millisecond * 50,
				Destination: &flags.interval,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) (err error) {
			var playlist bsides.Playlist
			if playlist, err = LoadPlaylist(cmd.StringArg("playlist")); err != nil {
				return err
			}

			var (
				iterator bsides.Iterator
				seed     bsides.Entry
			)
			if flags.shuffle {
				iterator = iters.NewShuffleIterator(playlist)
			} else {
				iterator = bsides.NewSequentialIterator(playlist)
				if flags.reverse {
					seed = playlist.Get(playlist.Back())
					iterator.Select(seed.ID)
				}
			}
			if flags.repeat > 0 {
				iterator = iters.NewRepeatingIterator(iterator, flags.repeat)
			}

			_, advance := iterator.HasNext, iterator.Next
			if flags.reverse {
				_, advance = iterator.HasPrev, iterator.Prev
			}

			var entry bsides.Entry
			if seed.IsZero() {
				entry = advance()
			} else {
				entry = seed
			}
			for !entry.IsZero() {
				_, _ = fmt.Fprintln(cmd.Writer, entry.Item.URI())
				time.Sleep(flags.interval)
				entry = advance()
			}

			return nil
		},
	}

	return root
}
