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

func (hub *Hub) HandleMessage(msg EventMessage) {
	log.Println("Handling message:", msg)
	// TODO Check the resources - whitelist?
	// TODO Handle user renames

	var err error
	if msg.Resource != "things" {
		err = fmt.Errorf("Unknown resource: %s", msg.Resource)
	}
	msg.Resource = "things"

	// TODO errors will overwrite
	// TODO persist the changes
	switch msg.Event {
	case "create":
		msg.Event = CREATE
	case "delete":
		msg.Event = DELETE
	case "update":
		msg.Event = UPDATE
	default:
		err = fmt.Errorf("Unknown method: %s", msg.Event)
	}

	// TODO return an error that will be sent to the sender only
	if err != nil {
		log.Printf("error: %s", err)
		return
	}

	log.Println("Broadcasting:", msg)
	hub.Broadcast(msg)
}

// Handler is the main websocket handler for users
func (hub *Hub) Handler(ws *websocket.Conn) {
	// Wrap the user, session key, and websocket together
	conn := Connection{ws: ws}

	// Examine the request for the session key and user
	r := ws.Request()
	if cookie, err := r.Cookie(hub.config.Cookie.Name); err == nil {
		conn.key = cookie.Value
	}

	// TODO this will request users even if cookie value was "" - shortcircuit?
	if conn.User = hub.sessions.GetUser(conn.key); !conn.User.Exists() {
		log.Printf("No user with session: %s", conn.key)
		return
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
	msg = EventMessage{
		Resource: "things",
		Event:    LIST,
		Content:  []db.Thing{{Name: "Hello"}},
	}
	websocket.JSON.Send(ws, msg)

	// Main event loop
Events:
	for {
		var event EventMessage
		if err := websocket.JSON.Receive(ws, &event); err != nil {
			log.Printf("error: parse error: %s", err)
			break Events
		}
		hub.HandleMessage(event)
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
