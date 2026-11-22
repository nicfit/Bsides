package bsides

// NewSequentialIterator iterates in sequence. Next will halt at the end, and Prev will halt at the beginning.
func NewSequentialIterator(p Playlist) Iterator {
	iter := &SequentialIterator{
		playlist: p,
	}
	iter.Reset()
	return iter
}

type SequentialIterator struct {
	playlist Playlist
	curr     ID
}

func (it *SequentialIterator) Reset() {
	it.curr = NilID
}

func (it *SequentialIterator) Current() Entry {
	return it.playlist.Get(it.curr)
}

func (it *SequentialIterator) Next() Entry {
	if it.curr == NilID {
		it.curr = it.playlist.Front()
		return it.playlist.Get(it.curr)
	}

	if next := it.playlist.Next(it.curr); !next.IsNil() {
		it.curr = next
		return it.playlist.Get(next)
	}
	return Entry{}
}

func (it *SequentialIterator) Prev() Entry {
	if prev := it.playlist.Prev(it.curr); !prev.IsNil() {
		it.curr = prev
		return it.playlist.Get(prev)
	}

	return Entry{}
}

func (it *SequentialIterator) Select(id ID) bool {
	if it.playlist.Get(id).IsZero() {
		return false
	}
	it.curr = id
	return true
}

func (it *SequentialIterator) HasNext() bool {
	return !it.playlist.Next(it.curr).IsNil()
}

func (it *SequentialIterator) HasPrev() bool {
	return !it.playlist.Prev(it.curr).IsNil()
}
