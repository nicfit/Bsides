package iters

import "github.com/nicfit/bsides"

func NewRepeatingIterator(iter bsides.Iterator, times int) bsides.Iterator {
	r := &RepeatIterator{
		base:  iter,
		times: times,
		count: 0,
	}
	return r
}

type RepeatIterator struct {
	iteratorWrapper
	count, times int
}

func (it *RepeatIterator) Reset() {
	it.count = 0
	it.base.Reset()
}

func (it *RepeatIterator) HasNext() bool {
	if it.exhausted() {
		return false
	}
	if it.base.HasNext() {
		return true
	}
	return it.base.Playlist().Len() > 0
}

func (it *RepeatIterator) HasPrev() bool {
	if it.exhausted() {
		return false
	}
	if it.base.HasPrev() {
		return true
	}
	return it.base.Playlist().Len() > 0
}

func (it *RepeatIterator) Next() bsides.Entry {
	if it.exhausted() {
		return bsides.Entry{}
	}
	if !it.base.HasNext() {
		it.base.Reset()
	}
	next := it.base.Next()
	if !next.IsZero() {
		it.count++
	}
	return next
}

func (it *RepeatIterator) Prev() bsides.Entry {
	if it.exhausted() {
		return bsides.Entry{}
	}
	prev := it.base.Prev()
	if !prev.IsZero() {
		it.count++
		return prev
	}

	it.base.Reset()

	var entry bsides.Entry
	if _, ok := it.base.(*bsides.SequentialIterator); ok {
		// This logic is good for sequential but not shuffled, it results predicable first result
		back := it.base.Playlist().Back()
		it.base.Select(back)
		entry = it.base.Playlist().Get(back)
	} else {
		entry = it.base.Prev()
	}
	if !entry.IsZero() {
		it.count++
	}
	return entry
}

func (it *RepeatIterator) exhausted() bool {
	return it.count >= it.times*it.Playlist().Len()
}
