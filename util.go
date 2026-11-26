package bsides

import (
	"fmt"
	"path/filepath"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
)

var AudioFileExtensions = []string{
	".mp3", ".m4a",
	".aac",
	".aiff",
	".alac",
	".amr", ".awb",
	".ape",
	".au",
	".cda",
	".flac",
	".mid", ".midi",
	".oga", ".ogg",
	".opus",
	".wav", ".wave",
	".webm",
	".wma",
	".wv",
}

func IsAudioFile(path string) bool {
	return slices.Contains(AudioFileExtensions, filepath.Ext(path))
}

func MakeTestPlaylist(t *testing.T, n int) Playlist {
	p := NewPlaylist()
	for i := range n {
		_, err := p.Add(NewItem(fmt.Sprintf("Track %d", i), fmt.Sprintf("file%d.mp3", i)))
		require.NoError(t, err)
	}
	return p
}
