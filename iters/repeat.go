package iters

import "github.com/nicfit/bsides"

func NewRepeatingIterator(iter bsides.Iterator) bsides.Iterator {
	r := &RepeatIterator{
		iter: iter,
	}
	r.Reset()
	return r
}

type RepeatIterator struct {
	iter bsides.Iterator
}

func (it *RepeatIterator) Reset() {
	it.iter.Reset()
}

func (it *RepeatIterator) Current() bsides.Entry {
	return it.iter.Current()
}

func (it *RepeatIterator) Select(id bsides.ID) bool {
	return it.iter.Select(id)
}

func (it *RepeatIterator) Playlist() bsides.Playlist {
	return it.iter.Playlist()
}

func (it *RepeatIterator) HasNext() bool {
	if it.iter.HasNext() {
		return true
	}
	return it.iter.Playlist().Len() > 0
}

func (it *RepeatIterator) HasPrev() bool {
	if it.iter.HasPrev() {
		return true
	}
	// FIXME: This is not true when wrapping Shuffle, which forces Prev=false
	return it.iter.Playlist().Len() > 0
}

func (it *RepeatIterator) Next() bsides.Entry {
	if !it.iter.HasNext() {
		it.iter.Reset()
	}
	return it.iter.Next()
}

func (it *RepeatIterator) Prev() bsides.Entry {
	prev := it.iter.Prev()
	if !prev.IsZero() {
		return prev
	}

	back := it.iter.Playlist().Back()
	it.iter.Select(back)
	return it.iter.Playlist().Get(back)
}
