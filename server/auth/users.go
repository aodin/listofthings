package auth

import (
	"fmt"
	"sync"
)

type UserManager struct {
	sync.RWMutex
	users map[int64]User
}

func (m *UserManager) Get(id int64) User {
	m.Lock()
	defer m.Unlock()
	return m.users[id]
}

func Users() *UserManager {
	return &UserManager{
		users: make(map[int64]User),
	}
}

type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (user User) Exists() bool {
	return user.ID != 0
}

func (user User) String() string {
	if user.Name == "" {
		return fmt.Sprintf("%d", user.ID)
	}
	return fmt.Sprintf("%s (%d)", user.Name, user.ID)
}
