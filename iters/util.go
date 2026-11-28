package iters

import "github.com/nicfit/bsides"

type iteratorWrapper struct {
	base bsides.Iterator
}

func (i *iteratorWrapper) Playlist() bsides.Playlist {
	return i.base.Playlist()
}

func (i *iteratorWrapper) Reset() {
	i.base.Reset()
}

func (i *iteratorWrapper) Current() bsides.Entry {
	return i.base.Current()
}

func (i *iteratorWrapper) Next() bsides.Entry {
	return i.base.Next()
}

func (i *iteratorWrapper) Prev() bsides.Entry {
	return i.base.Prev()
}

func (i *iteratorWrapper) Select(id bsides.ID) bool {
	return i.base.Select(id)
}

func (i *iteratorWrapper) HasNext() bool {
	return i.base.HasNext()
}

func (i *iteratorWrapper) HasPrev() bool {
	return i.base.HasPrev()
}
