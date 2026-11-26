package m3u

import (
	"time"

	"github.com/nicfit/bsides"
	eyed3 "github.com/nicfit/go-eyeD3"
)

// Track extends bside.Track with M3U details.
type Track interface {
	bsides.Track
	Group() string
	Size() int64
	AlbumArtURL() string
	Genre() string

	Init(metadata any)
}

func NewTrack(uri string) Track {
	return &entry{
		uri: uri,
	}
}

func (e *entry) URI() string         { return e.uri }
func (e *entry) Group() string       { return e.group }
func (e *entry) Size() int64         { return e.size }
func (e *entry) AlbumArtURL() string { return e.albumArtUrl }
func (e *entry) Genre() string       { return e.genre }

func (e *entry) DisplayName() string {
	if e.displayName != "" {
		return e.displayName
	}
	var dn string
	if e.title != "" {
		dn = e.title
	}
	if e.artist != "" {
		dn = e.artist + " - " + e.title
	}
	return dn
}

func (e *entry) Title() string           { return e.title }
func (e *entry) Duration() time.Duration { return e.duration }
func (e *entry) Artist() string          { return e.artist }
func (e *entry) AlbumArtist() string     { return e.artist }
func (e *entry) Album() string           { return e.album }

func (e *entry) Init(t any) {
	switch tt := t.(type) {
	case eyed3.Track:
		if tt.Path != "" {
			e.uri = tt.Path
		}
		if tt.Artist != "" {
			e.artist = tt.Artist
		}
		if tt.Album != "" {
			e.album = tt.Album
		}
		if tt.Title != "" {
			e.title = tt.Title
		}
		if tt.Genre != "" {
			e.genre = tt.Genre
		}
		if tt.Info.TimeSecs > 0 {
			e.duration = time.Second * time.Duration(tt.Info.TimeSecs)
		}
		if tt.Info.SizeBytes > 0 {
			e.size = int64(tt.Info.SizeBytes)
		}
	}
}

type entry struct {
	uri         string
	group       string
	artist      string
	album       string
	genre       string
	albumArtUrl string
	size        int64
	title       string
	displayName string
	duration    time.Duration
}
