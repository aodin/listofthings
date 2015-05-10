package server

import (
	"fmt"

	"code.google.com/p/go.net/websocket"
)

type User struct {
	conn *websocket.Conn
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (u User) String() string {
	if u.Name == "" {
		return fmt.Sprintf("%d", u.ID)
	}
	return fmt.Sprintf("%s (%d)", u.Name, u.ID)
}
