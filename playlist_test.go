package bsides

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPlaylist(t *testing.T) {
	pl := NewPlaylist()
	require.NotNil(t, pl)
	assert.Zero(t, pl.Len())
	assert.NotNil(t, pl.(*playlist).items)
	assert.NotNil(t, pl.(*playlist).index)
}

func TestNewPlaylistAddRemoveMove(t *testing.T) {
	ids := make([]ID, 0)

	pl := NewPlaylist()
	require.NotNil(t, pl)
	assert.Equal(t, 0, pl.Len())

	first := NewItem("1", "")
	first_id, err := pl.Add(first)
	require.NoError(t, err)
	// 1
	assert.Equal(t, 1, pl.Len())
	assert.NotEmpty(t, first_id)
	ids = append(ids, first_id)
	testEntriesEqual(t, []Entry{{ID: first_id, Item: first}}, pl.Entries())

	second := NewItem("2", "0")
	second_id, err := pl.Add(second)
	require.NoError(t, err)
	// 1, 2
	assert.Equal(t, 2, pl.Len())
	assert.NotEmpty(t, second_id)
	assert.NotContains(t, ids, second_id)
	ids = append(ids, second_id)

	// insert
	third := NewItem("3", "0")
	third_id, ok, err := pl.Insert(third, first_id)
	require.NoError(t, err)
	// 3, 1, 2
	assert.True(t, ok)
	assert.Equal(t, 3, pl.Len())
	assert.NotEmpty(t, third_id)
	assert.NotContains(t, ids, third_id)
	ids = append(ids, third_id)
	testEntriesEqual(t, []Entry{
		{ID: third_id, Item: third},
		{ID: first_id, Item: first},
		{ID: second_id, Item: second},
	}, pl.Entries())

	fourth := NewItem("4", "0")
	fourth_id, ok, err := pl.Insert(fourth, second_id)
	require.NoError(t, err)
	// 3, 1, 4, 2
	assert.True(t, ok)
	assert.Equal(t, 4, pl.Len())
	assert.NotEmpty(t, fourth_id)
	assert.NotContains(t, ids, fourth_id)
	ids = append(ids, fourth_id)
	testEntriesEqual(t, []Entry{
		{ID: third_id, Item: third},
		{ID: first_id, Item: first},
		{ID: fourth_id, Item: fourth},
		{ID: second_id, Item: second},
	}, pl.Entries())

	fifth := NewItem("5", "0")
	fifth_id, ok, err := pl.Insert(fifth, fourth_id)
	require.NoError(t, err)
	// 3, 1, 5, 4, 2
	assert.True(t, ok)
	assert.Equal(t, 5, pl.Len())
	assert.NotEmpty(t, fifth_id)
	assert.NotContains(t, ids, fifth_id)
	ids = append(ids, fifth_id)
	testEntriesEqual(t, []Entry{
		{ID: third_id, Item: third},
		{ID: first_id, Item: first},
		{ID: fifth_id, Item: fifth},
		{ID: fourth_id, Item: fourth},
		{ID: second_id, Item: second},
	}, pl.Entries())

	pl.Remove(second_id)
	assert.Equal(t, 4, pl.Len())

	pl.Remove(first_id)
	assert.Equal(t, 3, pl.Len())

	// 3, 5, 4
	testEntriesEqual(t, []Entry{
		{ID: third_id, Item: third},
		{ID: fifth_id, Item: fifth},
		{ID: fourth_id, Item: fourth},
	}, pl.Entries())

	err = pl.Move(pl.Entries()[1].ID, pl.Entries()[0].ID)
	require.NoError(t, err)
	// 5, 3, 4
	err = pl.Move(pl.Entries()[1].ID, pl.Entries()[2].ID)
	require.NoError(t, err)
	// 5, 4, 3

	testEntriesEqual(t, []Entry{
		{ID: fifth_id, Item: fifth},
		{ID: fourth_id, Item: fourth},
		{ID: third_id, Item: third},
	}, pl.Entries())
}

func TestPlaylistInsertNonExistent(t *testing.T) {
	pl := NewPlaylist()
	item := NewItem("file.mp3", "file.mp3")
	id, ok, err := pl.Insert(item, NewID())
	require.NoError(t, err)
	assert.False(t, ok)
	assert.Empty(t, id)
	assert.Equal(t, 0, pl.Len())
}

func TestNilAdd(t *testing.T) {
	pl := NewPlaylist()
	id, err := pl.Add(nil)
	assert.Error(t, err)
	assert.Empty(t, id)
	assert.Equal(t, 0, pl.Len())
}

func TestNilInsert(t *testing.T) {
	pl := NewPlaylist()
	pl.Add(NewItem("Resentment", ""))
	assert.Equal(t, 1, pl.Len())

	id, ok, err := pl.Insert(nil, pl.Entries()[0].ID)
	assert.Error(t, err)
	assert.Empty(t, id)
	assert.False(t, ok)
	assert.Equal(t, 1, pl.Len())
}

func TestPlaylistGet(t *testing.T) {
	tracker := map[ID]Item{}

	pl := NewPlaylist()
	for i := range 10 {
		item := NewItem("Hard Stance #"+strconv.Itoa(i), "http://music/file.mp3")
		id, err := pl.Add(item)
		require.NoError(t, err)
		tracker[id] = item
	}

	nope := NewID()
	assert.Zero(t, pl.Get(nope))
	for id, item := range tracker {
		assert.Equal(t, item, pl.Get(id).Item)
	}
}

func TestPlaylistMoveNonExistent(t *testing.T) {
	pl := NewPlaylist()
	id1, err := pl.Add(NewItem("Clash of Views", ""))
	require.NoError(t, err)
	id2, err := pl.Add(NewItem("No One To Blame", ""))
	require.NoError(t, err)
	assert.NotEqual(t, id1, id2)

	assert.NoError(t, pl.Move(id1, id2))
	assert.Error(t, pl.Move(NewID(), id2))
	assert.Error(t, pl.Move(id1, NewID()))
}

func TestPlaylistFrontBack(t *testing.T) {
	empty := NewPlaylist()
	assert.Equal(t, NilID, empty.Front())
	assert.Equal(t, NilID, empty.Back())

	pl := MakeTestPlaylist(t, 13)
	assert.Equal(t, "Track 0", pl.Get(pl.Front()).Item.DisplayName())
	assert.Equal(t, "Track 12", pl.Get(pl.Back()).Item.DisplayName())
}

func TestPlaylistNextPrev(t *testing.T) {

	sz := 13
	pl := MakeTestPlaylist(t, sz)
	entries := pl.Entries()
	for i := range sz {
		var expected_prev, expected_next ID
		switch i {
		case 0:
			expected_prev = NilID
			expected_next = entries[i+1].ID
		case sz - 1:
			expected_prev = entries[i-1].ID
			expected_next = NilID
		default:
			expected_prev = entries[i-1].ID
			expected_next = entries[i+1].ID
		}

		assert.Equal(t, expected_prev, pl.Prev(entries[i].ID))
		assert.Equal(t, expected_next, pl.Next(entries[i].ID))
	}
}

func TestPlaylist_Entries(t *testing.T) {
	// This also verifies makeTestPlaylist
	pl := MakeTestPlaylist(t, 133)
	assert.Equal(t, 133, pl.Len())
	for i, e := range pl.Entries() {
		assert.Equal(t, fmt.Sprintf("Track %d", i), e.Item.DisplayName())
		assert.Equal(t, fmt.Sprintf("file%d.mp3", i), e.Item.URI())
	}
}

func TestPlaylist_ForEach(t *testing.T) {
	pl := MakeTestPlaylist(t, 133)
	entries := pl.Entries()
	count := 0
	pl.ForEach(func(id ID, item Item) bool {
		assert.Equal(t, entries[count].ID, id)
		assert.Equal(t, entries[count].Item, item)
		count++
		return true
	})

	// Interrupt loop
	count = 0
	pl.ForEach(func(id ID, item Item) bool {
		if count == 45 {
			return false
		}
		count++
		return true
	})
	assert.Equal(t, 45, count)
}

///////////////////////////////////////////////////////////////////////

func testEntriesEqual(t *testing.T, expected, actual []Entry) {
	assert.Equal(t, len(expected), len(actual))
	for i := range expected {
		assert.Equal(t, expected[i].ID, actual[i].ID)
		assert.Equal(t, expected[i].Item.DisplayName(), actual[i].Item.DisplayName())
		assert.Equal(t, expected[i].Item.URI(), actual[i].Item.URI())
	}
}
