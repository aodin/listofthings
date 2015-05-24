package v1

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"code.google.com/p/go.net/websocket"
	sql "github.com/aodin/aspect"
	pg "github.com/aodin/aspect/postgres"
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
}

// Hub matches session keys to connections
type Hub struct {
	sync.RWMutex
	config      config.Config
	conn        sql.Connection
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
	msg := OutgoingMessage{
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
	msg := OutgoingMessage{
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

func unmarshalThing(msg IncomingMessage) (thing db.Thing, err error) {
	if err = json.Unmarshal(msg.Content, &thing); err != nil {
		return
	}
	// TODO remove the timestamps from the content?
	// TODO whitelisted fields?
	thing.Content = string(msg.Content)
	err = thing.Error()
	return
}

type Things []db.Thing

func (t Things) Mutate() {
	for i, thing := range t {
		// TODO ignored error
		json.Unmarshal([]byte(thing.Content), &thing)

		// Only copy content fields - preserve system info
		t[i].Name = thing.Name
	}
}

// TODO error?
func getThings(conn sql.Connection) (things Things) {
	things = Things{}
	conn.MustQueryAll(db.Things.Select().Where(
		db.Things.C["deleted_at"].IsNull(),
	).OrderBy(db.Things.C["id"]), &things)
	things.Mutate()
	return
}

func (hub *Hub) HandleMessage(in IncomingMessage) {
	log.Println("Handling message:", in)
	// TODO Check the resources - whitelist?
	// TODO Handle user renames

	var err error
	var out OutgoingMessage
	if in.Resource != "things" {
		err = fmt.Errorf("Unknown resource: %s", in.Resource)
	}
	out.Resource = "things"

	// TODO errors will overwrite
	// TODO persist the changes
	var thing db.Thing
	switch in.Event {
	case "create":
		out.Event = CREATE
		if thing, err = unmarshalThing(in); err == nil {
			stmt := pg.Insert(db.Things).Values(thing).Returning(db.Things)
			err = hub.conn.QueryOne(stmt, &thing)
			out.Content = thing
		}
	case "delete":
		out.Event = DELETE
		// TODO confirm ID?
		if thing, err = unmarshalThing(in); err == nil {
			stmt := db.Things.Delete().Where(
				db.Things.C["id"].Equals(thing.ID),
			)
			_, err = hub.conn.Execute(stmt)
			out.Content = thing
		}
	case "update":
		out.Event = UPDATE
		// TODO confirm ID?
		if thing, err = unmarshalThing(in); err == nil {
			stmt := db.Things.Update().Values(thing.Values()).Where(
				db.Things.C["id"].Equals(thing.ID),
			)
			_, err = hub.conn.Execute(stmt)
			out.Content = thing
		}
	default:
		err = fmt.Errorf("Unknown method: %s", in.Event)
	}

	// TODO return an error that will be sent to the sender only
	if err != nil {
		log.Printf("error: %s", err)
		return
	}

	log.Println("Broadcasting:", out)
	hub.Broadcast(out)
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
	msg := OutgoingMessage{
		Resource: "users",
		Event:    LIST,
		Content:  hub.Users(),
	}
	websocket.JSON.Send(ws, msg)

	// Send the initial state of the data
	msg = OutgoingMessage{
		Resource: "things",
		Event:    LIST,
		Content:  getThings(hub.conn),
	}
	websocket.JSON.Send(ws, msg)

	// Main event loop
Events:
	for {
		var event IncomingMessage
		if err := websocket.JSON.Receive(ws, &event); err != nil {
			log.Printf("error: parse error: %s", err)
			break Events
		}
		hub.HandleMessage(event)
	}

	hub.Leave(conn)
}

func NewHub(config config.Config, conn sql.Connection, sessions *auth.SessionManager) *Hub {
	return &Hub{
		config:      config,
		conn:        conn,
		sessions:    sessions,
		connections: make(map[string]Connection),
	}
}
