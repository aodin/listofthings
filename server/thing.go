package server

import (
	"fmt"
	"time"
)

// A basic resource
type Thing struct {
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
	Timestamp time.Time `json:"timestamp"`
}

func (t *Thing) String() string {
	return fmt.Sprintf("%s", t.Name)
}
