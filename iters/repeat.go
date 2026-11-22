package iters

import "github.com/nicfit/bsides"

type RepeatMode int

const (
	RepeatPlaylist = iota
	RepeatItem
)

func NewRepeatingIterator(pl bsides.Playlist, mode RepeatMode) bsides.Iterator {
	iter := &RepeatIterator{
		SequentialIterator: bsides.NewSequentialIterator(pl).(*bsides.SequentialIterator),
		playlist:           pl,
		mode:               mode,
	}
	iter.Reset()
	return iter
}

type RepeatIterator struct {
	*bsides.SequentialIterator

	mode     RepeatMode
	playlist bsides.Playlist
}

func (it *RepeatIterator) HasNext() bool {
	return it.playlist.Len() > 0
}

func (it *RepeatIterator) HasPrev() bool {
	return it.playlist.Len() > 0
}

func (it *RepeatIterator) Next() bsides.Entry {
	if it.mode == RepeatItem {
		if it.Current().IsZero() {
			it.Select(it.playlist.Front())
		}
		return it.Current()
	}

	next := it.SequentialIterator.Next()
	if !next.IsZero() {
		return next
	}

	it.SequentialIterator.Reset()
	return it.SequentialIterator.Next()
}

func (it *RepeatIterator) Prev() bsides.Entry {
	if it.mode == RepeatItem {
		if it.Current().IsZero() {
			it.Select(it.playlist.Back())
		}
		return it.Current()
	}

	prev := it.SequentialIterator.Prev()
	if !prev.IsZero() {
		return prev
	}

	back := it.playlist.Back()
	it.SequentialIterator.Select(back)
	return it.playlist.Get(back)
}
