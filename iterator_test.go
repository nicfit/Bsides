package bsides

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStandardIteratorNext(t *testing.T) {
	p := MakeTestPlaylist(t, 99)
	it := NewSequentialIterator(p)

	assert.Zero(t, it.Current())

	// Two passes through the iterator to ensure Reset() works correctly
	for range 2 {
		count := 0
		for entry := it.Next(); !entry.IsZero(); entry = it.Next() {
			require.NotNil(t, entry.Item)
			require.Equal(t, fmt.Sprintf("Track %d", count), entry.Item.DisplayName())
			count++
		}
		assert.Equal(t, 99, count)

		// No more calls to Next() should return nil
		assert.True(t, it.HasPrev())
		assert.False(t, it.HasNext())
		for range 10 {
			assert.Zero(t, it.Next())
		}

		it.Reset()
		assert.Zero(t, it.Current())
	}
}

func TestStandardIteratorPrev(t *testing.T) {
	p := MakeTestPlaylist(t, 99)
	it := NewSequentialIterator(p)

	assert.Zero(t, it.Current())
	assert.Zero(t, it.Prev())
	// Set the iterator to the end of the playlist
	assert.False(t, it.Select(NewID()))
	assert.True(t, it.Select(p.Entries()[p.Len()-1].ID))
	assert.Equal(t, "Track 98", it.Current().Item.DisplayName())

	// Two passes through the iterator to ensure Reset() works correctly
	for range 2 {
		count := 97
		for entry := it.Prev(); !entry.IsZero(); entry = it.Prev() {
			require.NotNil(t, entry.Item)
			require.Equal(t, fmt.Sprintf("Track %d", count), entry.Item.DisplayName())
			count--
		}
		assert.Equal(t, -1, count)
		assert.NotZero(t, it.Current())
		assert.Equal(t, fmt.Sprintf("Track %d", 0), it.Current().Item.DisplayName())

		// No more calls to Prev() should return nil
		assert.True(t, it.HasNext())
		assert.False(t, it.HasPrev())
		for range 10 {
			assert.Zero(t, it.Prev())
		}

		it.Reset()
		assert.Zero(t, it.Current())
		assert.True(t, it.Select(p.Entries()[p.Len()-1].ID))
	}
}

func TestStandardIteratorEmptyPlaylist(t *testing.T) {
	p := MakeTestPlaylist(t, 0)
	it := NewSequentialIterator(p)
	assert.Zero(t, it.Current())
}
