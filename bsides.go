package bsides

import "uuid"

// Playlist is a sequential collection of Items.
type Playlist interface {
	// Len returns the number of tracks in the playlist.
	Len() int
	// Add adds a track to the playlist and returns its unique identifier. An error is returned when
	// the [Track] is invalid or cannot be added to the playlist.
	Add(Item) (ID, error)
	// Remove removes a track from the playlist by its unique identifier.
	// If the track was found and removed, it returns true; otherwise, it returns false.
	Remove(id ID) bool
	// Insert adds a track to the playlist at the position of the given unique identifier.
	// If the track was successfully inserted, it returns the new unique identifier and true;
	// otherwise, it returns false. An error is returned when the [Track] is invalid or cannot be inserted
	// into the playlist.
	Insert(item Item, id ID) (ID, bool, error)
	// Move moves src to dst. Both Ids must exist or an error is returned.
	Move(src, dst ID) error

	// Get returns the [Entry] by [ID] or a zero value if not found.
	Get(id ID) Entry

	// Entries returns a slice of all entries in the playlist.
	Entries() []Entry

	// Next returns the next ID following the argument `id`, or [NilID] if there is none.
	// These methods are to support [Iterator] implementations, and those should be preferred to traversing playlists.
	Next(id ID) ID
	// Prev returns the ID previous to the argument `id`, or [NilID] if there is none.
	Prev(id ID) ID
	// Front returns the first ID in the Playlist. If the playlist is empty [NilID] is returned.
	Front() ID
	// Back returns the last ID in the Playlist. If the playlist is empty [NilID] is returned.
	Back() ID
}

// Item types are stored in Playlists.
type Item interface {
	URI() string
	DisplayName() string
}

// ID is a unique key for a playlist entry.
type ID uuid.UUID

// NewID returns a time-based random playlist ID (i.e, UUIDv7)
func NewID() ID {
	return ID(uuid.NewV7())
}
func (id ID) String() string {
	return uuid.UUID(id).String()
}
func (id ID) IsNil() bool {
	return id == ID(uuid.Nil())
}

// NilID represent no ID.
var NilID = ID(uuid.Nil())

// Entry represents a single entry in a playlist, consisting of a
// unique identifier and the associated item.
type Entry struct {
	ID   ID
	Item Item
}

// IsZero is a convenience for `entry == Entry{}`, since zero-value
// Entry are used to represent None.
func (entry Entry) IsZero() bool {
	return entry == Entry{}
}

// Iterator types traverse playlists in a particular order.
type Iterator interface {
	// Reset resets the iterator to the beginning of the iteration.
	Reset()
	// Current returns the current entry, or zero-value if none.
	Current() Entry
	// Next returns the next entry, or zero-value if none.
	Next() Entry
	// Prev returns the previous entry, or zero-value if none.
	Prev() Entry
	// Select sets the current entry by playlist ID. It returns true if the
	// entry was found and set; otherwise, it returns false.
	Select(id ID) bool
	// HasNext returns true if there is a next entry; otherwise, it returns false.
	HasNext() bool
	// HasPrev returns true if there is a previous entry; otherwise, it returns false.
	HasPrev() bool
}
