package v1

import (
	"fmt"
	"log"
	"sync"

	"code.google.com/p/go.net/websocket"
	"github.com/aodin/volta/config"

	db "github.com/aodin/listofthings/db"
	"github.com/aodin/listofthings/server/auth"
)

type Connection struct {
	db.User
	key string // Session key
	ws  *websocket.Conn
}

func (c Connection) String() string {
	return fmt.Sprintf("%s (id: %d)", c.User, c.User.ID)
	// return fmt.Sprintf("%s (id: %d, key: %s)", c.User, c.User.ID, c.key)
}

// Hub matches session keys to connections
type Hub struct {
	sync.RWMutex
	config      config.Config
	sessions    *auth.SessionManager
	connections map[string]Connection
}

// Broadcast sends a message to all users
func (hub *Hub) Broadcast(msg Message) {
	for _, user := range hub.connections {
		// TODO error ignored
		_ = websocket.JSON.Send(user.ws, msg)
	}
}

func (hub *Hub) Join(connection Connection) {
	// Log and broadcast the event
	log.Printf("%s joined\n", connection)
	msg := EventMessage{
		Resource: "users",
		Event:    CREATE,
		Content:  connection.User,
	}
	hub.Broadcast(msg)

	// Add the connection to the hub
	hub.Lock()
	defer hub.Unlock()
	hub.connections[connection.key] = connection
}

func (hub *Hub) Leave(connection Connection) {
	// Remove the connection from the hub
	hub.Lock()
	delete(hub.connections, connection.key)
	hub.Unlock()

	// Log and broadcast the event
	log.Printf("%s left\n", connection)
	msg := EventMessage{
		Resource: "users",
		Event:    DELETE,
		Content:  connection.User,
	}
	hub.Broadcast(msg)
}

func (hub *Hub) Users() []db.User {
	// This list will include the requesting user
	// TODO Does order matter?
	users := make([]db.User, len(hub.connections))
	var i int
	for _, connection := range hub.connections {
		users[i] = connection.User
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
	// Wrap the user, session key, and websocket together
	conn := Connection{ws: ws}

	// Examine the request for the session key and user
	r := ws.Request()
	if cookie, err := r.Cookie(hub.config.Cookie.Name); err == nil {
		conn.key = cookie.Value
		if conn.User = hub.sessions.GetUser(conn.key); !conn.User.Exists() {
			log.Printf("No user with session: %s", conn.key)
			return
		}
	}

	hub.Join(conn)

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

	hub.Leave(conn)
}

func NewHub(config config.Config, sessions *auth.SessionManager) *Hub {
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
		config:      config,
		sessions:    sessions,
		connections: make(map[string]Connection),
	}
}
