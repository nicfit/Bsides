package bsides

import (
	"container/list"
	"errors"
	"fmt"
	"sync"
)

// NewPlaylist creates a thread-safe, doubly linked list implementation of the Playlist interface.
// Along with the list is a O(1) index into the list by ID.a
func NewPlaylist() Playlist {
	p := &playlist{
		items: list.New(),
		index: make(map[ID]*list.Element),
	}
	return p
}

// playlist is a thread-safe doubly linked list implementation of the Playlist interface.
// Along with the list is a O(1) index into the list by ID.
type playlist struct {
	mu    sync.RWMutex
	items *list.List           // Main doubly linked list
	index map[ID]*list.Element // O(1) lookup by Id
}

func (p *playlist) Len() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.items.Len()
}

func (p *playlist) Add(item Item) (ID, error) {
	if item == nil {
		return ID{}, errors.New("item is nil")
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	id := NewID()
	elem := p.items.PushBack(&Entry{
		ID:   id,
		Item: item,
	})
	p.index[id] = elem

	return id, nil
}

func (p *playlist) Remove(id ID) bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	var (
		exists bool
		elem   *list.Element
	)
	if elem, exists = p.index[id]; exists {
		p.items.Remove(elem)
		delete(p.index, id)
	}

	return exists
}

func (p *playlist) Insert(item Item, id ID) (ID, bool, error) {
	if item == nil {
		return NilID, false, errors.New("item is nil")
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if elem, exists := p.index[id]; exists {
		new_id := NewID()
		new_elem := p.items.InsertBefore(&Entry{ID: new_id, Item: item}, elem)
		p.index[new_id] = new_elem

		return new_id, true, nil
	}

	return NilID, false, nil
}

func (p *playlist) Move(src, dst ID) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	src_elem, src_exists := p.index[src]
	dst_elem, dst_exists := p.index[dst]
	if !src_exists {
		return fmt.Errorf("src ID does not exist: %v", src)
	} else if !dst_exists {
		return fmt.Errorf("dest ID does not exist: %v", dst)
	}

	if p.index[dst] == p.items.Back() {
		p.items.MoveAfter(src_elem, dst_elem)
	} else {
		p.items.MoveBefore(src_elem, dst_elem)
	}
	return nil
}

func (p *playlist) Get(id ID) Entry {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if elem, exists := p.index[id]; exists {
		return *elem.Value.(*Entry)
	}
	return Entry{}
}

func (p *playlist) Entries() []Entry {
	p.mu.RLock()
	defer p.mu.RUnlock()

	entries := make([]Entry, 0, p.items.Len())
	for elem := p.items.Front(); elem != nil; elem = elem.Next() {
		entry := elem.Value.(*Entry)
		entries = append(entries, *entry)
	}

	return entries
}

func (p *playlist) Front() ID {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if p.items.Len() > 0 {
		return p.items.Front().Value.(*Entry).ID
	}
	return NilID
}

func (p *playlist) Back() ID {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if p.items.Len() > 0 {
		return p.items.Back().Value.(*Entry).ID
	}
	return NilID
}

func (p *playlist) Next(id ID) ID {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.findNeighbor(id, true)
}

func (p *playlist) Prev(id ID) ID {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.findNeighbor(id, false)
}

func (p *playlist) findNeighbor(id ID, forward bool) ID {
	if _, exists := p.index[id]; exists {
		for elem := p.items.Front(); elem != nil; elem = elem.Next() {
			if elem.Value.(*Entry).ID == id {
				if forward && elem.Next() != nil {
					return elem.Next().Value.(*Entry).ID
				} else if !forward && elem.Prev() != nil {
					return elem.Prev().Value.(*Entry).ID
				}
			}
		}
	}
	return NilID
}

func (p *playlist) ForEach(f func(ID, Item) bool) {
	for elem := p.items.Front(); elem != nil; elem = elem.Next() {
		if !f(elem.Value.(*Entry).ID, elem.Value.(*Entry).Item) {
			break
		}
	}
}
