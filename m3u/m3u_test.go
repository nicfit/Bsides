package m3u

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	m3u, err := Load("sample.m3u")
	assert.NoError(t, err)
	assert.Nil(t, m3u)

	m3uExt, err := Load("sample-extended.m3u")
	assert.NoError(t, err)
	assert.Nil(t, m3uExt)

	m3u8, err := Load("sample.m3u8")
	assert.NoError(t, err)
	assert.Nil(t, m3u8)

	dne, err := Load("Nirvana-Beeswax.m3u")
	assert.Error(t, err)
	assert.Nil(t, dne)

	huh, err := Load("m3u.go")
	assert.Error(t, err)
	assert.Nil(t, huh)
}
