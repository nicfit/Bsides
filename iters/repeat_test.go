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
	for _, mode := range []RepeatMode{RepeatItem, RepeatPlaylist} {
		iter = NewRepeatingIterator(empty, mode)
		assert.False(t, iter.HasNext())
		assert.False(t, iter.HasPrev())
		assert.Zero(t, iter.Next())
		assert.Zero(t, iter.Prev())
		assert.False(t, iter.Select(bsides.NewID()))
		assert.False(t, iter.Select(empty.Front()))
		assert.False(t, iter.Select(empty.Back()))
	}
}

func TestNewRepeatingIteratorSizeOne(t *testing.T) {
	var iter bsides.Iterator
	pl := bsides.MakeTestPlaylist(t, 1)

	for _, iter = range []bsides.Iterator{
		NewRepeatingIterator(pl, RepeatItem), NewRepeatingIterator(pl, RepeatPlaylist),
	} {
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
}

func TestNewRepeatingIteratorSizeTwo(t *testing.T) {
	var iter bsides.Iterator
	pl := bsides.MakeTestPlaylist(t, 2)

	// RepeatItem
	iter = NewRepeatingIterator(pl, RepeatItem)
	// The iterator is not started, so...
	assert.Zero(t, iter.Current())
	assert.True(t, iter.HasNext())
	assert.True(t, iter.HasPrev())
	next := iter.Next()
	assert.Equal(t, "Track 0", next.Item.DisplayName())
	assert.True(t, iter.HasNext())
	assert.True(t, iter.HasPrev())
	next = iter.Next()
	assert.Equal(t, "Track 0", next.Item.DisplayName())

	// RepeatPlaylist
	iter = NewRepeatingIterator(pl, RepeatPlaylist)
	assert.Zero(t, iter.Current())
	assert.True(t, iter.HasNext())
	assert.True(t, iter.HasPrev())
	next = iter.Next()
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

func TestRepeatingIteratorItem(t *testing.T) {
	pl := bsides.MakeTestPlaylist(t, 99)
	iter := NewRepeatingIterator(pl, RepeatItem)

	assert.True(t, iter.Select(pl.Front()))
	assert.Equal(t, "Track 0", iter.Next().Item.DisplayName())
	for range 25 {
		assert.Equal(t, "Track 0", iter.Next().Item.DisplayName())
	}
	for range 25 {
		assert.Equal(t, "Track 0", iter.Prev().Item.DisplayName())
	}
}

func TestRepeatingIteratorPlaylistForward(t *testing.T) {
	pl := bsides.MakeTestPlaylist(t, 99)

	iter := NewRepeatingIterator(pl, RepeatPlaylist)
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

	iter := NewRepeatingIterator(pl, RepeatPlaylist)
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
