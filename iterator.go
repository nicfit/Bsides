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
	started  bool
}

func (it *SequentialIterator) Playlist() Playlist {
	return it.playlist
}

func (it *SequentialIterator) Reset() {
	it.curr = NilID
	it.started = false
}

func (it *SequentialIterator) Current() Entry {
	return it.playlist.Get(it.curr)
}

func (it *SequentialIterator) Next() (next Entry) {
	defer func() {
		if !next.IsZero() {
			it.started = true
		}
	}()

	if it.curr == NilID {
		it.curr = it.playlist.Front()
		next = it.playlist.Get(it.curr)
		return
	}

	if n := it.playlist.Next(it.curr); !n.IsNil() {
		it.curr = n
		next = it.playlist.Get(n)
		return
	}

	return
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
	if !it.started && it.playlist.Len() > 0 {
		return true
	}
	return !it.playlist.Next(it.curr).IsNil()
}

func (it *SequentialIterator) HasPrev() bool {
	return !it.playlist.Prev(it.curr).IsNil()
}
