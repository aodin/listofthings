package server

import (
	"fmt"
	"html"
	"sync"

	"github.com/aodin/listofthings/db"
	"github.com/aodin/listofthings/db/fields"
)

type Store interface {
	List() []*db.Thing
	Create(*db.Thing) (*db.Thing, error)
	Delete(*db.Thing) (*db.Thing, error)
	Update(*db.Thing) (*db.Thing, error)
}

type InMemoryStore struct {
	sync.Mutex
	things []*db.Thing
}

func (m *InMemoryStore) List() []*db.Thing {
	// TODO Read-writer mutex? Do we need to lock on read?

	// Remove nil things
	things := make([]*db.Thing, 0)
	for _, thing := range m.things {
		if thing == nil {
			continue
		}
		things = append(things, thing)
	}
	return things
}

// TODO Should validation be methods of the Resource or the Store?

func (m *InMemoryStore) Create(thing *db.Thing) (*db.Thing, error) {
	// Lock the memory storage while this occurs
	m.Lock()
	defer m.Unlock()

	// Find an open slot
	var index int64
	for i, t := range m.things {
		if t == nil {
			index = int64(i) + 1
			break
		}
	}

	if index == 0 {
		// There are no more open slots
		return nil, fmt.Errorf("Only %d items can be stored.", len(m.things))
	}

	// Escape all items on input
	thing.ID = index
	thing.Name = html.EscapeString(thing.Name)
	thing.Timestamp = fields.NewTimestamp()

	// Add it to the memory store
	m.things[index-1] = thing
	return thing, nil
}

func (m *InMemoryStore) Delete(thing *db.Thing) (*db.Thing, error) {
	// Lock the memory storage while this occurs
	m.Lock()
	defer m.Unlock()

	// Does the thing exist?
	index := int(thing.ID)
	if index > len(m.things) || index < 1 {
		return nil, fmt.Errorf("Invalid id")
	}

	thing = m.things[index-1]
	if thing == nil {
		return nil, fmt.Errorf("Thing does not exist")
	}

	// Delete it
	m.things[index-1] = nil
	return thing, nil
}

func (m *InMemoryStore) Update(thing *db.Thing) (*db.Thing, error) {
	// Lock the memory storage while this occurs
	m.Lock()
	defer m.Unlock()

	// Does the index make sense
	index := int(thing.ID)
	if index > len(m.things) || index < 1 {
		return nil, fmt.Errorf("Invalid id")
	}

	original := m.things[index-1]
	if original == nil {
		return nil, fmt.Errorf("Thing does not exist")
	}

	// Maintain the timestamp and id from the original
	original.Name = html.EscapeString(thing.Name)

	// Update it
	m.things[index-1] = original
	return original, nil
}

func NewMemoryStore(n int) *InMemoryStore {
	return &InMemoryStore{things: make([]*db.Thing, n)}
}
