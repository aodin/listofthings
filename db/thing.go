package db

import (
	"fmt"

	"github.com/aodin/listofthings/db/fields"
)

const MaxNameLength = 256

// Thing is a thing with a name
type Thing struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	fields.Timestamp
}

func (t Thing) String() string {
	return t.Name
}

func (t Thing) Error() error {
	if t.Name == "" {
		return fmt.Errorf("Names cannot be blank")
	}
	if len(t.Name) > MaxNameLength {
		return fmt.Errorf(
			"Names cannot be longer than %d characters",
			MaxNameLength,
		)
	}
	return nil
}

func NewThing(name string) Thing {
	return Thing{Name: name}
}
