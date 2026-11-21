package m3u

import (
	"errors"
	"time"

	"github.com/nicfit/bsides"
)

// Track extends bside.Track with M3U details.
type Track interface {
	bsides.Track
	Group() string
	Size() int64
	AlbumArtURL() string
	Genre() string
}

func NewTrack(uri string) (Track, error) {
	if uri == "" {
		return nil, errors.New("empty uri")
	}
	return &entry{
		uri: uri,
	}, nil
}

func (e *entry) URI() string         { return e.uri }
func (e *entry) Group() string       { return e.group }
func (e *entry) Size() int64         { return e.size }
func (e *entry) AlbumArtURL() string { return e.albumArtUrl }
func (e *entry) Genre() string       { return e.genre }

func (e *entry) DisplayName() string     { return e.displayName }
func (e *entry) Title() string           { return e.displayName }
func (e *entry) Duration() time.Duration { return e.duration }
func (e *entry) Artist() string          { return e.artist }
func (e *entry) AlbumArtist() string     { return e.artist }
func (e *entry) Album() string           { return e.album }

type entry struct {
	uri         string
	group       string
	artist      string
	album       string
	genre       string
	albumArtUrl string
	size        int64
	displayName string
	duration    time.Duration
}
