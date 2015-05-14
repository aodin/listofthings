package v1

import (
	"log"
	"sync"

	"code.google.com/p/go.net/websocket"

	"github.com/aodin/listofthings/server/auth"
)

// Hub matches session keys to connections
type Hub struct {
	sync.RWMutex
	users       *auth.UserManager
	sessions    *auth.SessionManager
	connections map[string]*websocket.Conn
}

// Broadcast sends a message to all users
func (hub *Hub) Broadcast(msg Message) {
	for _, connection := range hub.connections {
		// TODO error ignored
		_ = websocket.JSON.Send(connection, msg)
	}
}

func (hub *Hub) Join(session auth.Session, ws *websocket.Conn) {
	// Get the user if it exists
	var user auth.User
	if session.UserID != 0 {
		user = hub.users.Get(session.UserID)
	}

	if !user.Exists() {
		// Create a anonymous user
		user = auth.User{Name: "Anonymous User"}
	}

	msg := EventMessage{
		Resource: "users",
		Event:    CREATE,
	}
	hub.Broadcast(msg)
	hub.Lock()
	defer hub.Unlock()
	hub.connections[session.Key] = ws
}

func (hub *Hub) Leave(session auth.Session) {
	hub.Lock()
	defer hub.Unlock()
	delete(hub.connections, session.Key)

	msg := EventMessage{
		Resource: "users",
		Event:    DELETE,
		Content:  hub.users.Get(session.UserID),
	}
	hub.Broadcast(msg)
}

func (hub *Hub) Users() []auth.User {
	// This list will include the requesting user
	users := make([]auth.User, len(hub.connections))
	var i int
	for _, _ = range hub.connections {
		// TODO Match users
		users[i] = auth.User{Name: "Anonymous User"}
		i += 1
	}
	return users
}

// func (hub *Hub) Handler() {
// 	var err error
// 	var thing *db.Thing

// 	// TODO Handle user renames

// 	switch msg.Method {
// 	case "create":
// 		thing, err = s.store.Create(msg.Item)
// 	case "delete":
// 		thing, err = s.store.Delete(msg.Item)
// 	case "update":
// 		thing, err = s.store.Update(msg.Item)
// 	default:
// 		err = fmt.Errorf("Unknown method '%s'", msg.Method)
// 	}

// 	// Build the return message
// 	var returnMsg v1.Message
// 	if err != nil {
// 		returnMsg.Body = "error"
// 		returnMsg.Content = err.Error()
// 	} else {
// 		returnMsg.Body = msg.Method
// 		returnMsg.Content = thing
// 	}
// 	s.BroadcastMessage(&returnMsg)
// }

// Handler is the main websocket handler for users
func (hub *Hub) Handler(ws *websocket.Conn) {
	// Wait for the connect message
	var connect ConnectMessage
	if err := websocket.JSON.Receive(ws, &connect); err != nil {
		// Exit early
		log.Println("Failed to connect:", err)
		return
	}

	// Create a connection with the session
	session := hub.sessions.Get(connect.SessionKey)
	if !session.Exists() {
		// TODO Force an expiration in session? Create a new session?
		log.Printf("Session does not exist: %s", connect.SessionKey)
		return
	}

	hub.Join(session, ws)

	// Send the initial state of the users list
	msg := EventMessage{
		Resource: "users",
		Event:    LIST,
		Content:  hub.Users(),
	}
	websocket.JSON.Send(ws, msg)

	// Send the initial state of the data

	// Main event loop
Events:
	for {
		var event EventMessage
		if err := websocket.JSON.Receive(ws, &event); err != nil {
			break Events
		}
		// s.HandleMessage(&msg)
	}

	hub.Leave(session)
	log.Printf("Session %s exited\n", session.Key)
}

func NewHub(users *auth.UserManager, sessions *auth.SessionManager) *Hub {
	// Create a memory store with a limited number of items
	// store := NewMemoryStore(25)

	// // Add a few things
	// names := []string{"Bass-o-matic", "Swill", "Jam Hawkers"}

	// for _, name := range names {
	// 	thing := db.NewThing(name)
	// 	if _, err := store.Create(&thing); err != nil {
	// 		log.Panic(err)
	// 	}
	// }

	// http.Handle("/events", websocket.Handler(srv.EventsHandler))

	return &Hub{
		connections: make(map[string]*websocket.Conn),
		users:       users,
		sessions:    sessions,
	}
}
