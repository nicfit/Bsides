package iters

import (
	"fmt"
	"testing"

	"github.com/nicfit/bsides"
	"github.com/stretchr/testify/assert"
)

func TestNewRepeatingIteratorSizeZero(t *testing.T) {
	// An empty playlist
	var iter bsides.Iterator
	empty := bsides.NewPlaylist()
	iter = NewRepeatingIterator(bsides.NewSequentialIterator(empty), 3)
	assert.Equal(t, empty, iter.Playlist())
	assert.False(t, iter.HasNext())
	assert.False(t, iter.HasPrev())
	assert.Zero(t, iter.Next())
	assert.Zero(t, iter.Prev())
	assert.False(t, iter.Select(bsides.NewID()))
	assert.False(t, iter.Select(empty.Front()))
	assert.False(t, iter.Select(empty.Back()))
}

func TestNewRepeatingIteratorSizeOne(t *testing.T) {
	pl := bsides.MakeTestPlaylist(t, 1)

	iter := NewRepeatingIterator(bsides.NewSequentialIterator(pl), 3)
	// The iterator is not started, so...
	assert.Zero(t, iter.Current())
	// Size >0 means there is always a next/prev
	assert.True(t, iter.HasNext())
	assert.True(t, iter.HasPrev())

	next := iter.Next()
	assert.Equal(t, "Track 0", next.Item.DisplayName())
	assert.True(t, iter.HasNext())
	assert.True(t, iter.HasPrev())
	next = iter.Next()
	assert.Equal(t, "Track 0", next.Item.DisplayName())
}

func TestNewRepeatingIteratorSizeTwo(t *testing.T) {
	var iter bsides.Iterator
	pl := bsides.MakeTestPlaylist(t, 2)

	iter = NewRepeatingIterator(bsides.NewSequentialIterator(pl), 3)
	assert.Zero(t, iter.Current())
	assert.True(t, iter.HasNext())
	assert.True(t, iter.HasPrev())
	next := iter.Next()
	assert.Equal(t, "Track 0", next.Item.DisplayName())
	next = iter.Next()
	assert.Equal(t, "Track 1", next.Item.DisplayName())
	next = iter.Next()
	assert.Equal(t, "Track 0", next.Item.DisplayName())
	next = iter.Prev()
	assert.Equal(t, "Track 1", next.Item.DisplayName())
	next = iter.Prev()
	assert.Equal(t, "Track 0", next.Item.DisplayName())
}

func TestRepeatingIteratorPlaylistForward(t *testing.T) {
	pl := bsides.MakeTestPlaylist(t, 99)

	iter := NewRepeatingIterator(bsides.NewSequentialIterator(pl), 100)
	count := 0
	for range 1000 {
		assert.True(t, iter.HasNext())
		assert.True(t, iter.HasPrev())
		next := iter.Next()
		assert.Equal(t, fmt.Sprintf("Track %d", count), next.Item.DisplayName())
		count++
		if count == 99 {
			count = 0
		}
	}
}

func TestRepeatingIteratorPlaylistBackward(t *testing.T) {
	pl := bsides.MakeTestPlaylist(t, 99)

	iter := NewRepeatingIterator(bsides.NewSequentialIterator(pl), 100)
	count := 98
	for range 1000 {
		assert.True(t, iter.HasNext())
		assert.True(t, iter.HasPrev())
		next := iter.Prev()
		assert.Equal(t, fmt.Sprintf("Track %d", count), next.Item.DisplayName())
		count--
		if count == -1 {
			count = 98
		}
	}
}

func TestRepeatExhaustionNext(t *testing.T) {
	sz, times := 99, 2
	pl := bsides.MakeTestPlaylist(t, sz)
	iter := NewRepeatingIterator(bsides.NewSequentialIterator(pl), times)
	for i := range sz * (times + 1) {
		next := iter.Next()
		if i < sz*times {
			assert.NotZero(t, next)
		} else {
			assert.Zero(t, next)
			assert.False(t, iter.HasNext())
		}
	}
}

func TestRepeatExhaustionPrev(t *testing.T) {
	sz, times := 99, 2
	pl := bsides.MakeTestPlaylist(t, sz)
	iter := NewRepeatingIterator(bsides.NewSequentialIterator(pl), times)
	for i := range sz * (times + 1) {
		prev := iter.Prev()
		if i < sz*times {
			assert.NotZero(t, prev)
		} else {
			assert.Zero(t, prev)
			assert.False(t, iter.HasPrev())
		}
	}
}
