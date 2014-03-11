package server

import (
	"fmt"
	"html"
	"sync"
)

type Thing struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

func (t *Thing) String() string {
	return fmt.Sprintf("%s", t.Name)
}

type Storage interface {
	List() []*Thing
	Create(*Thing) (*Thing, error)
	Delete(*Thing) (*Thing, error)
	Update(*Thing) (*Thing, error)
}

type Memory struct {
	mutex  sync.Mutex
	things []*Thing
}

func (m *Memory) List() []*Thing {
	// TODO Read-writer mutex

	// Remove nil things
	things := make([]*Thing, 0)
	for _, thing := range m.things {
		if thing == nil {
			continue
		}
		things = append(things, thing)
	}
	return things
}

// Escape all items on input
func (m *Memory) Create(thing *Thing) (*Thing, error) {
	// Lock the memory storage while this occurs
	m.mutex.Lock()
	defer m.mutex.Unlock()

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

	thing.Id = index
	thing.Name = html.EscapeString(thing.Name)
	m.things[index-1] = thing
	return thing, nil
}

func (m *Memory) Delete(thing *Thing) (*Thing, error) {
	// Lock the memory storage while this occurs
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Does the thing exist?
	index := int(thing.Id)
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

func (m *Memory) Update(thing *Thing) (*Thing, error) {
	// Lock the memory storage while this occurs
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Does the index make sense
	index := int(thing.Id)
	if index > len(m.things) || index < 1 {
		return nil, fmt.Errorf("Invalid id")
	}

	// TODO validation
	thing.Name = html.EscapeString(thing.Name)

	// Update it
	m.things[index-1] = thing
	return thing, nil
}

func NewMemoryStore(n int) *Memory {
	return &Memory{things: make([]*Thing, n)}
}
