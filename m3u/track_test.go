package m3u

import (
	"testing"
	"time"

	eyed3 "github.com/nicfit/go-eyeD3"
	"github.com/stretchr/testify/assert"
)

func TestEntry_Init(t *testing.T) {
	eye := eyed3.Track{
		Path: "/Iceburn/Firon/burnII.mp3",
		Info: eyed3.Info{
			TimeSecs:  (time.Minute*7 + time.Second*28).Seconds(),
			SizeBytes: 4352345,
		},
		Title:               "Burn II",
		Artist:              "Iceburn",
		Album:               "Firon",
		AlbumArtist:         "Iceburn",
		AlbumType:           "LP",
		BestReleaseDate:     "1992",
		Genre:               "Post Hardcore",
		OriginalReleaseDate: "1992",
		Publisher:           "Victory",
		ReleaseDate:         "1992",
		TrackNum:            eyed3.Count{Count: 1, Total: 11},
	}

	track := NewTrack("")
	track.Init(eye)
	assert.Equal(t, eye.Path, track.URI())
	assert.Equal(t, eye.Title, track.Title())
	assert.Equal(t, eye.Artist, track.Artist())
}
