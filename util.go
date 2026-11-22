package bsides

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func MakeTestPlaylist(t *testing.T, n int) Playlist {
	p := NewPlaylist()
	for i := range n {
		_, err := p.Add(NewItem(fmt.Sprintf("Track %d", i), fmt.Sprintf("file%d.mp3", i)))
		require.NoError(t, err)
	}
	return p
}
