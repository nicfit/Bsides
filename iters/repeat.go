package iters

import "github.com/nicfit/bsides"

func NewRepeatingIterator(iter bsides.Iterator, times int) bsides.Iterator {
	r := &RepeatIterator{
		iter:  iter,
		times: times,
		count: 0,
	}
	return r
}

type RepeatIterator struct {
	iter         bsides.Iterator
	count, times int
}

func (it *RepeatIterator) Reset() {
	it.count = 0
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
	if it.exhausted() {
		return false
	}
	if it.iter.HasNext() {
		return true
	}
	return it.iter.Playlist().Len() > 0
}

func (it *RepeatIterator) HasPrev() bool {
	if it.exhausted() {
		return false
	}
	if it.iter.HasPrev() {
		return true
	}
	return it.iter.Playlist().Len() > 0
}

func (it *RepeatIterator) Next() bsides.Entry {
	if it.exhausted() {
		return bsides.Entry{}
	}
	if !it.iter.HasNext() {
		it.iter.Reset()
	}
	next := it.iter.Next()
	if !next.IsZero() {
		it.count++
	}
	return next
}

func (it *RepeatIterator) Prev() bsides.Entry {
	if it.exhausted() {
		return bsides.Entry{}
	}
	prev := it.iter.Prev()
	if !prev.IsZero() {
		it.count++
		return prev
	}

	it.iter.Reset()

	var entry bsides.Entry
	if _, ok := it.iter.(*bsides.SequentialIterator); ok {
		// This logic is good for sequential but not shuffled, it results predicable first result
		back := it.iter.Playlist().Back()
		it.iter.Select(back)
		entry = it.iter.Playlist().Get(back)
	} else {
		entry = it.iter.Prev()
	}
	if !entry.IsZero() {
		it.count++
	}
	return entry
}

func (it *RepeatIterator) exhausted() bool {
	return it.count > it.times*it.Playlist().Len()
}
