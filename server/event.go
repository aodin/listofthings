package server

import ()

// A event that occurs
type Event struct {
	Method string // Create, Update, Delete
	Id     int64  // zero for item that should be created
}
