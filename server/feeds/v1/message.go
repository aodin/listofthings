package v1

import (
	"encoding/json"
	"fmt"
)

type Message interface {
	String() string
}

type OutgoingMessage struct {
	Resource string      `json:"resource"`
	Event    string      `json:"method"`
	Content  interface{} `json:"content"`
}

func (msg OutgoingMessage) String() string {
	return fmt.Sprintf("%s %s: %+v", msg.Event, msg.Resource, msg.Content)
}

type IncomingMessage struct {
	Resource string          `json:"resource"`
	Event    string          `json:"method"`
	Content  json.RawMessage `json:"content"`
}

func (msg IncomingMessage) String() string {
	return fmt.Sprintf("%s %s: %s", msg.Event, msg.Resource, msg.Content)
}
