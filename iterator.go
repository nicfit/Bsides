package bsides

// SequentialIterator iterates in sequence. Next will halt at the end, and Prev will halt at the beginning.
func SequentialIterator(p Playlist) Iterator {
	iter := &sequentialIterator{
		playlist: p,
	}
	iter.Reset()
	return iter
}

type sequentialIterator struct {
	playlist Playlist
	curr     ID
}

func (it *sequentialIterator) Reset() {
	it.curr = NilID
}

func (it *sequentialIterator) Current() Entry {
	return it.playlist.Get(it.curr)
}

func (it *sequentialIterator) Next() Entry {
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

func (it *sequentialIterator) Prev() Entry {
	if prev := it.playlist.Prev(it.curr); !prev.IsNil() {
		it.curr = prev
		return it.playlist.Get(prev)
	}

	return Entry{}
}

func (it *sequentialIterator) Select(id ID) bool {
	if it.playlist.Get(id).IsZero() {
		return false
	}
	it.curr = id
	return true
}

func (it *sequentialIterator) HasNext() bool {
	return !it.playlist.Next(it.curr).IsNil()
}

func (it *sequentialIterator) HasPrev() bool {
	return !it.playlist.Prev(it.curr).IsNil()
}
