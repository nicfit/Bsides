package m3u

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/nicfit/bsides"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	m3u, err := Load("sample.m3u")
	assert.NoError(t, err)
	assert.NotNil(t, m3u)

	m3uExt, err := Load("sample-extended.m3u")
	assert.NoError(t, err)
	assert.NotNil(t, m3uExt)

	m3u8, err := Load("sample.m3u8")
	assert.NoError(t, err)
	assert.NotNil(t, m3u8)

	dne, err := Load("Nirvana-Beeswax.m3u")
	assert.Error(t, err)
	assert.Nil(t, dne)

	huh, err := Load("m3u.go")
	assert.Error(t, err)
	assert.Nil(t, huh)
}

func Test_NewTrack(t *testing.T) {
	track, err := NewTrack("zombie dance.mp3")
	require.NoError(t, err)
	require.NotNil(t, track)
	assert.Equal(t, "zombie dance.mp3", track.URI())
}

func Test_NewTrack_noURI(t *testing.T) {
	track, err := NewTrack("")
	require.Error(t, err)
	require.Nil(t, track)
}

func TestLoadSample(t *testing.T) {
	m3u, err := Load("sample.m3u")
	require.NoError(t, err)
	require.NotNil(t, m3u)
	assert.Equal(t, "sample.m3u", m3u.Path())
	assert.Equal(t, "", m3u.Name())
	assert.False(t, m3u.IsExtended())
	assert.Equal(t, 5, m3u.Len())

	for _, entry := range m3u.Entries() {
		assert.NotEmpty(t, entry.ID)
		track, ok := entry.Item.(Track)
		assert.True(t, ok)

		// NO extra meta data
		assert.Empty(t, track.Title())
		assert.Empty(t, track.Duration())
		assert.Empty(t, track.Artist())
		assert.Empty(t, track.Album())
		assert.Empty(t, track.AlbumArtURL())
		assert.Empty(t, track.Group())
		assert.Zero(t, track.Size())
		assert.Empty(t, track.Genre())
	}
}

func TestLoadSampleExtended(t *testing.T) {
	m3u, err := Load("sample-extended.m3u")
	require.NotNil(t, m3u)
	require.NoError(t, err)
	assert.Equal(t, "sample-extended.m3u", m3u.Path())
	assert.Equal(t, "Sample Playlist", m3u.Name())
	assert.True(t, m3u.IsExtended())
	assert.Equal(t, 5, m3u.Len())

	// unique ids
	ids := make([]bsides.ID, 0, 16)
	for _, entry := range m3u.Entries() {
		assert.NotContains(t, ids, entry.ID)
		ids = append(ids, entry.ID)
	}

	for _, entry := range m3u.Entries() {
		track, ok := entry.Item.(Track)
		require.True(t, ok)

		if slices.Contains([]string{"file1.mp3", "file2.mp3"}, track.URI()) {
			assert.Equal(t, "Bangers", track.Group())
		} else {
			assert.Empty(t, track.Group())
		}

		if "file1.mp3" == track.URI() {
			assert.Equal(t, "Turning Point", track.Artist())
			assert.Equal(t, "Demo", track.Album())
			assert.Equal(t, "Hardcore", track.Genre())
			assert.Equal(t, "https://tinyurl.com/3apbd3ne", track.AlbumArtURL())
		} else {
			assert.Empty(t, track.Artist())
			assert.Empty(t, track.Album())
			assert.Empty(t, track.Genre())
			assert.Empty(t, track.AlbumArtURL())
		}

		if "file3.mp3" == track.URI() {
			assert.Equal(t, int64(123456), track.Size())
		} else {
			assert.Zero(t, track.Size())
		}

		if "file5.mp3" == track.URI() {
			assert.Equal(t, "The Brian Jonestown Massacre: Who?", track.Title())
			assert.Equal(t, 172*time.Second, track.Duration())
		} else {
			assert.Equal(t, "", track.Title())
			assert.Equal(t, 0*time.Second, track.Duration())
		}
	}
}

func TestLoad404(t *testing.T) {
	m3u, err := Load("404.m3u")
	assert.Nil(t, m3u)
	assert.Error(t, err)
}

func TestMultipleExtPlaylistName(t *testing.T) {
	sampleOk := `#EXTM3U
#PLAYLIST:Sample Playlist
track1.mp3
track2.mp3
#PLAYLIST:Sample Playlist
track3.mp3
`
	_, err := LoadData(strings.NewReader(sampleOk))
	require.NoError(t, err)

	sampleNotOk := `#EXTM3U
#PLAYLIST:Sample Playlist
track1.mp3
track2.mp3
#PLAYLIST:SampleXXX Playlist
track3.mp3
`
	_, err = LoadData(strings.NewReader(sampleNotOk))
	require.Error(t, err)
}

func Test_extSplitNone(t *testing.T) {
	before, after := extSplit("In My Life")
	assert.Equal(t, "In My Life", before)
	assert.Equal(t, "", after)
}

func TestInvalid_EXT_BYTES(t *testing.T) {
	sampleOk := `#EXTM3U
#EXTBYT:Be Song
track1.mp3
`
	_, err := loadData(strings.NewReader(sampleOk))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid size")
}

func Test_EXT_INF_duration(t *testing.T) {

	for _, inf := range []string{"#EXTINF:205s,BJM - Straight Up and Down"} {
		sample := "#EXTM3U\n" + inf + "\ntrack1.mp3\n"
		_, err := loadData(strings.NewReader(sample))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid duration")
	}

	// Okay to omit it
	for _, inf := range []string{"#EXTINF:,BJM - Be Song", "#EXTINF: ,BJM - Be Song", "#EXTINF:0,BJM - Be Song"} {
		sample := "#EXTM3U\n" + inf + "\ntrack1.mp3\n"
		m3u, err := loadData(strings.NewReader(sample))
		require.NoError(t, err)
		assert.Zero(t, m3u.Entries()[0].Item.(Track).Duration())
	}
}

func TestErrCatchLoadPath(t *testing.T) {
	path, err := os.CreateTemp("", "*.m3u")
	require.NoError(t, err)
	require.NotNil(t, path)
	defer func() { _ = os.Remove(path.Name()) }()
	_, err = path.Write([]byte("#EXTM3U\n#EXTINF:666s,Stars and Sons\ntrack.mp3"))
	require.NoError(t, err)
	assert.NoError(t, path.Close())

	_, err = Load(path.Name())
	assert.Error(t, err)
}

func TestNewPlaylist(t *testing.T) {
	m3u := NewPlaylist()
	require.NotNil(t, m3u)
	assert.Empty(t, m3u.Path())
	assert.Empty(t, m3u.Name())
	assert.False(t, m3u.IsExtended())
	assert.Equal(t, 0, m3u.Len())
}

func TestAddInsertWrongType(t *testing.T) {
	m3u, err := Load("sample-extended.m3u")
	require.NoError(t, err)
	require.NotNil(t, m3u)

	pid, err := m3u.Add(bsides.NewItem("The Start of Your Ending", "mobbdeep.mp3"))
	assert.Empty(t, pid)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "requires m3u.Track")

	pid, inserted, err := m3u.Insert(bsides.NewItem("The Start of Your Ending", "tune.ogg"),
		m3u.Entries()[3].ID)
	assert.Empty(t, pid)
	assert.False(t, inserted)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "requires m3u.Track")
}

func TestPlaylist_AddRemoveInsertMove(t *testing.T) {
	m3u, err := Load("sample-extended.m3u")
	require.NoError(t, err)
	require.NotNil(t, m3u)

	assert.Equal(t, 5, m3u.Len())
	for i := range 5 {
		assert.Equal(t, fmt.Sprintf("file%d.mp3", i+1), m3u.Entries()[i].Item.(Track).URI())
	}

	// Add
	track6, err := NewTrack("file6.mp3")
	require.NoError(t, err)
	pid, err := m3u.Add(track6)
	require.NoError(t, err)
	assert.Equal(t, pid, m3u.Entries()[5].ID)

	track7, err := NewTrack("file7.mp3")
	require.NoError(t, err)
	pid, err = m3u.Add(track7)
	require.NoError(t, err)
	assert.Equal(t, pid, m3u.Entries()[6].ID)

	assert.Equal(t, 7, m3u.Len())

	// Remove
	m3u.Remove(m3u.Entries()[2].ID)
	assert.Equal(t, 6, m3u.Len())

	// Insert

	track8, err := NewTrack("file8.mp3")
	require.NoError(t, err)
	pid, inserted, err := m3u.Insert(track8, m3u.Entries()[0].ID)
	require.NoError(t, err)
	assert.True(t, inserted)
	assert.Equal(t, pid, m3u.Entries()[0].ID)

	track9, err := NewTrack("file9.mp3")
	require.NoError(t, err)
	pid, inserted, err = m3u.Insert(track9, m3u.Entries()[5].ID)
	require.NoError(t, err)
	assert.True(t, inserted)
	assert.Equal(t, pid, m3u.Entries()[5].ID)

	assert.Equal(t, 8, m3u.Len())

	// Test expected order
	for i, d := range []int{8, 1, 2, 4, 5, 9, 6, 7} {
		assert.Equal(t, fmt.Sprintf("file%d.mp3", d), m3u.Entries()[i].Item.(Track).URI())
	}

	src_id := m3u.Entries()[5].ID // 9
	err = m3u.Move(src_id, m3u.Entries()[0].ID)
	require.NoError(t, err)
	for i, d := range []int{9, 8, 1, 2, 4, 5, 6, 7} {
		assert.Equal(t, fmt.Sprintf("file%d.mp3", d), m3u.Entries()[i].Item.(Track).URI())
	}
	src_id = m3u.Entries()[2].ID // 1
	err = m3u.Move(src_id, m3u.Entries()[7].ID)
	require.NoError(t, err)
	for i, d := range []int{9, 8, 2, 4, 5, 6, 7, 1} {
		assert.Equal(t, fmt.Sprintf("file%d.mp3", d), m3u.Entries()[i].Item.(Track).URI())
	}
	src_id = m3u.Entries()[4].ID // 5
	err = m3u.Move(src_id, m3u.Entries()[3].ID)
	require.NoError(t, err)
	for i, d := range []int{9, 8, 2, 5, 4, 6, 7, 1} {
		assert.Equal(t, fmt.Sprintf("file%d.mp3", d), m3u.Entries()[i].Item.(Track).URI())
	}
	// No-op
	err = m3u.Move(m3u.Entries()[5].ID, m3u.Entries()[5].ID)
	require.NoError(t, err)
	for i, d := range []int{9, 8, 2, 5, 4, 6, 7, 1} {
		assert.Equal(t, fmt.Sprintf("file%d.mp3", d), m3u.Entries()[i].Item.(Track).URI())
	}

	// Invalid src
	err = m3u.Move(bsides.NewID(), m3u.Entries()[5].ID)
	require.ErrorContains(t, err, "src ID does not exist:")
	// Invalid dest
	err = m3u.Move(m3u.Entries()[5].ID, bsides.NewID())
	require.ErrorContains(t, err, "dest ID does not exist:")
}

func TestPlaylistExport(t *testing.T) {
	m3u, err := Load("sample.m3u")
	assert.NoError(t, err)

	contents := m3u.Export()
	assert.Equal(t, "file1.mp3\nfile2.mp3\nfile3.mp3\nfile4.mp3\nfile5.mp3\n", contents)
	m3u.SetExtended(true)
	ext_contents := m3u.Export()
	assert.Equal(t, "#EXTM3U\nfile1.mp3\nfile2.mp3\nfile3.mp3\nfile4.mp3\nfile5.mp3\n", ext_contents)

	m3u_ext, err := Load("sample-extended.m3u")
	assert.NoError(t, err)
	assert.True(t, m3u_ext.IsExtended())
	ext_contents = m3u_ext.Export()
	assert.Equal(t, "#EXTM3U\n#PLAYLIST:Sample Playlist\n#EXTGRP:Bangers\n#EXTART:Turning Point\n#EXTALB:Demo\n#EXTALBUMARTURL:https://tinyurl.com/3apbd3ne\n#EXTGENRE:Hardcore\nfile1.mp3\n#EXTGRP:Bangers\nfile2.mp3\n#EXTBYT:123456\nfile3.mp3\nfile4.mp3\n#EXTINF:172,The Brian Jonestown Massacre: Who?\nfile5.mp3\n", ext_contents)
	m3u.SetExtended(false)
	ext_contents = m3u.Export()
	assert.Equal(t, "file1.mp3\nfile2.mp3\nfile3.mp3\nfile4.mp3\nfile5.mp3\n", contents)
}
