package bsides

import (
	"sort"
	"testing"
	"uuid"

	"github.com/stretchr/testify/assert"
)

func TestID(t *testing.T) {
	assert.Equal(t, ID(uuid.Nil()), NilID)

	id := ID{}
	assert.Zero(t, id)
	assert.True(t, id.IsNil())
	assert.Equal(t, "00000000-0000-0000-0000-000000000000", id.String())
	assert.Equal(t, id, NilID)
}

func TestNewID(t *testing.T) {
	id := NewID()
	assert.False(t, id.IsNil())
	assert.NotEqual(t, NilID, id)
	assert.NotEqual(t, NewID().String(), NewID().String())

	id2 := NewID()
	id3 := NewID()
	series := []string{id2.String(), id3.String(), id.String(), NilID.String()}
	sort.Strings(series)
	assert.Equal(t, series[0], NilID.String())
	assert.Equal(t, series[1], id.String())
	assert.Equal(t, series[2], id2.String())
	assert.Equal(t, series[3], id3.String())
}

func TestEntry(t *testing.T) {
	e := Entry{}
	assert.True(t, e.IsZero())
	assert.Equal(t, NilID, e.ID)
	assert.Nil(t, e.Item)
}
