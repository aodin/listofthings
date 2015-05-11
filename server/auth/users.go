package auth

import (
	"fmt"
	"sync"
)

type UserManager struct {
	sync.RWMutex
	users map[int64]User
}

type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (user User) String() string {
	if user.Name == "" {
		return fmt.Sprintf("%d", user.ID)
	}
	return fmt.Sprintf("%s (%d)", user.Name, user.ID)
}

func Users() *UserManager {
	return &UserManager{
		users: make(map[int64]User),
	}
}
