package iters

import (
	"math/rand/v2"
	"slices"

	"github.com/nicfit/bsides"
)

func NewShuffleIterator(pl bsides.Playlist) bsides.Iterator {
	iter := &ShuffleIterator{
		playlist: pl,
		shuffled: make([]bsides.ID, 0, pl.Len()),
	}
	iter.Reset()
	return iter
}

type ShuffleIterator struct {
	playlist bsides.Playlist
	shuffled []bsides.ID
	curr     bsides.ID
}

func (s *ShuffleIterator) Playlist() bsides.Playlist {
	return s.playlist
}

func (s *ShuffleIterator) Reset() {
	s.curr = bsides.NilID
	s.playlist.ForEach(func(id bsides.ID, item bsides.Item) bool {
		s.shuffled = append(s.shuffled, id)
		return true
	})
	rand.Shuffle(len(s.shuffled), func(i, j int) {
		s.shuffled[i], s.shuffled[j] = s.shuffled[j], s.shuffled[i]
	})
}

func (s *ShuffleIterator) Current() bsides.Entry {
	if s.curr == bsides.NilID {
		return bsides.Entry{}
	}
	return s.playlist.Get(s.curr)
}

func (s *ShuffleIterator) Next() bsides.Entry {
	entry := s.choose()
	if entry.IsZero() {
		s.curr = bsides.NilID
	}
	return entry
}

func (s *ShuffleIterator) Prev() bsides.Entry {
	entry := s.choose()
	if entry.IsZero() {
		s.curr = bsides.NilID
	}
	return entry
}

func (s *ShuffleIterator) choose() bsides.Entry {
	// Loop to handle items that may have been removed
	for len(s.shuffled) > 0 {
		s.curr = s.shuffled[len(s.shuffled)-1]
		s.shuffled = s.shuffled[:len(s.shuffled)-1]
		if entry := s.playlist.Get(s.curr); !entry.IsZero() {
			return entry
		}
	}
	return bsides.Entry{}
}

func (s *ShuffleIterator) Select(id bsides.ID) bool {
	for i, shuffled_id := range s.shuffled {
		if id == shuffled_id {
			s.curr = id
			s.shuffled = slices.Delete(s.shuffled, i, i+1)
			return true
		}
	}
	return false
}

func (s *ShuffleIterator) HasNext() bool {
	return len(s.shuffled) > 0
}

func (s *ShuffleIterator) HasPrev() bool {
	return len(s.shuffled) > 0
}
