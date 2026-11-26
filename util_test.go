package bsides

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsAudioFile(t *testing.T) {
	exts := slices.Clone(AudioFileExtensions)
	for _, ext := range exts {
		assert.True(t, IsAudioFile("file"+ext))
		assert.True(t, IsAudioFile("/dir/file"+ext))
		assert.True(t, IsAudioFile("dir/file"+ext))
		assert.True(t, IsAudioFile("dir/music file.nope"+ext))
	}
	assert.False(t, IsAudioFile("file"))
	assert.False(t, IsAudioFile("cover.jpg"))
	assert.False(t, IsAudioFile("playlist.m3u"))
}
