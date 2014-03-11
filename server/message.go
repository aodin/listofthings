package server

import (
	"fmt"
)

// Messages that will be sent by the server and users
type Message struct {
	Id      int64       `json:"id"`
	Body    string      `json:"body"`
	Content interface{} `json:"content"`
}

// Event messages for create-update-delete events on resources
type ResourceMessage struct {
	Method  string   `json:"method"`
	Content Resource `json:"content"`
}

// Event messages for create-update-delete events specific to Thing
type ThingMessage struct {
	Method string `json:"method"`
	Item   *Thing `json:"content"`
}

func (t ThingMessage) String() string {
	return fmt.Sprintf("%s: %s", t.Method, t.Item)
}
