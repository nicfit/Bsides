package iters

import (
	"sort"
	"testing"

	"github.com/nicfit/bsides"
	"github.com/stretchr/testify/assert"
)

func TestNewShuffleIterator(t *testing.T) {
	pl := bsides.MakeTestPlaylist(t, 87)

	ordered := make([]string, 0, pl.Len())
	pl.ForEach(func(id bsides.ID, item bsides.Item) bool {
		ordered = append(ordered, id.String())
		return true
	})
	assert.Equal(t, pl.Len(), len(ordered))

	iter := NewShuffleIterator(pl).(*ShuffleIterator)
	assert.Equal(t, pl, iter.Playlist())
	assert.Zero(t, iter.Current())
	assert.Equal(t, pl.Len(), len(iter.shuffled))

	shuffled := make([]string, 0, pl.Len())
	for _, id := range iter.shuffled {
		shuffled = append(shuffled, id.String())
	}
	assert.NotEqual(t, ordered, shuffled)
	// Using uuidV7 means a string short can get back to the original order
	sort.Strings(shuffled)
	assert.Equal(t, ordered, shuffled)
}

func TestNewShuffleIteratorHasNextPrev(t *testing.T) {
	pl := bsides.MakeTestPlaylist(t, 8)
	iter := NewShuffleIterator(pl).(*ShuffleIterator)
	assert.Equal(t, pl.Len(), len(iter.shuffled))
	assert.Zero(t, iter.Current())
	assert.True(t, iter.HasNext())
	assert.True(t, iter.HasPrev())
}

func TestNewShuffleIteratorNext(t *testing.T) {
	pl := bsides.MakeTestPlaylist(t, 92)
	iter := NewShuffleIterator(pl).(*ShuffleIterator)
	assert.Equal(t, pl.Len(), len(iter.shuffled))
	assert.Zero(t, iter.Current())
	assert.True(t, iter.HasNext())

	var next bsides.Entry
	count := 0
	for iter.HasNext() {
		next = iter.Next()
		count++
		assert.LessOrEqual(t, count, pl.Len())
	}
	assert.Equal(t, pl.Len(), count)
	// The loop never resulted in a call to Next when !HasNext, so Current is valid
	assert.NotZero(t, iter.Current())
	assert.Equal(t, next, iter.Current())
	assert.False(t, iter.HasNext())
	assert.Zero(t, iter.Next())
	assert.Zero(t, iter.Current())
}

func TestNewShuffleIteratorPrev(t *testing.T) {
	pl := bsides.MakeTestPlaylist(t, 41)
	iter := NewShuffleIterator(pl).(*ShuffleIterator)
	assert.Equal(t, pl.Len(), len(iter.shuffled))
	assert.Zero(t, iter.Current())
	assert.True(t, iter.HasPrev())

	var prev bsides.Entry
	count := 0
	for iter.HasPrev() {
		prev = iter.Prev()
		count++
		assert.LessOrEqual(t, count, pl.Len())
	}
	assert.Equal(t, pl.Len(), count)
	// The loop never resulted in a call to Prev when !Hasprev, so Current is valid
	assert.NotZero(t, iter.Current())
	assert.Equal(t, prev, iter.Current())
	assert.False(t, iter.HasPrev())
	assert.Zero(t, iter.Prev())
	assert.Zero(t, iter.Current())
}

func TestNewShuffleIteratorEmpty(t *testing.T) {
	pl := bsides.MakeTestPlaylist(t, 0)
	iter := NewShuffleIterator(pl).(*ShuffleIterator)
	assert.False(t, iter.HasPrev())
	assert.False(t, iter.HasNext())
	assert.Zero(t, iter.Current())
	assert.Zero(t, iter.Next())
}

func TestNewShuffleIteratorSelect(t *testing.T) {
	pl := bsides.MakeTestPlaylist(t, 4)
	third := pl.Entries()[2].ID
	iter := NewShuffleIterator(pl).(*ShuffleIterator)

	assert.True(t, iter.Select(third))
	assert.Equal(t, pl.Get(third), iter.Current())
	assert.NotContains(t, iter.shuffled, third)

	assert.False(t, iter.Select(bsides.NewID()))

	assert.True(t, iter.Select(pl.Entries()[3].ID))
	assert.True(t, iter.Select(pl.Entries()[0].ID))
	assert.True(t, iter.HasNext())

	// Last item...
	assert.True(t, iter.Select(pl.Entries()[1].ID))
	assert.False(t, iter.HasNext())

	iter.Reset()
	assert.True(t, iter.HasNext())
}

func TestShuffleWithRepeating(t *testing.T) {
	pl := bsides.MakeTestPlaylist(t, 634)
	shuffle := NewShuffleIterator(pl)
	repeating := NewRepeatingIterator(shuffle, 5)
	for range 3 {
		assert.NotZero(t, repeating.Next())
	}
	for range 5 {
		assert.NotZero(t, repeating.Prev())
	}
}
