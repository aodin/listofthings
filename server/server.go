package server

import (
	"fmt"
	"log"
	"net/http"

	"code.google.com/p/go.net/websocket"
	"github.com/aodin/volta/config"
	"github.com/aodin/volta/templates"

	"github.com/aodin/listofthings/db"
	"github.com/aodin/listofthings/server/feeds/v1"
)

// Wrap HTTP methods
type Server struct {
	config    config.Config
	templates *templates.Templates
	users     map[*User]bool
	counter   int64
	store     Store
}

func (srv *Server) ListenAndServe() error {
	return http.ListenAndServe(srv.config.Address(), nil)
}

// Index is the handler for the index
func (srv *Server) IndexHandler(w http.ResponseWriter, r *http.Request) {
	// Assign a session if one has not been set
	srv.templates.Execute(w, "index")
}

// Broadcast method for users
func (s *Server) BroadcastMessage(msg *v1.Message) {
	// Send the message to all users
	for user, _ := range s.users {
		err := websocket.JSON.Send(user.conn, msg)
		if err != nil {
			continue
		}
	}
}

// User Events
// -----------
func (s *Server) AddUser(user *User) {
	// Broadcast a user joined message to all users
	msg := &v1.Message{
		Body:    "join",
		Content: user,
	}

	for user, _ := range s.users {
		err := websocket.JSON.Send(user.conn, msg)
		if err != nil {
			continue
		}
	}

	// Add the user afterwards
	s.users[user] = true
}

func (s *Server) DeleteUser(user *User) {
	// Delete the user
	delete(s.users, user)

	msg := &v1.Message{
		Body:    "left",
		Content: user,
	}

	// Send the remaining users a user left message
	for user, _ := range s.users {
		err := websocket.JSON.Send(user.conn, msg)
		if err != nil {
			continue
		}
	}
}

// Message handler
func (s *Server) HandleMessage(msg *v1.ThingMessage) {
	var err error
	var thing *db.Thing

	// TODO Handle user renames

	switch msg.Method {
	case "create":
		thing, err = s.store.Create(msg.Item)
	case "delete":
		thing, err = s.store.Delete(msg.Item)
	case "update":
		thing, err = s.store.Update(msg.Item)
	default:
		err = fmt.Errorf("Unknown method '%s'", msg.Method)
	}

	// Build the return message
	var returnMsg v1.Message
	if err != nil {
		returnMsg.Body = "error"
		returnMsg.Content = err.Error()
	} else {
		returnMsg.Body = msg.Method
		returnMsg.Content = thing
	}
	s.BroadcastMessage(&returnMsg)
}

// Main websocket handler for users
// Receives new events
func (s *Server) EventsHandler(ws *websocket.Conn) {
	// Register the new connection
	s.counter += 1

	// Create a user object and join the server
	user := &User{ID: s.counter, conn: ws}
	log.Println("User connected:", user)

	// Determine the number of current users
	others := make([]*User, len(s.users))
	var userIndex int
	for user, _ := range s.users {
		others[userIndex] = user
		userIndex += 1
	}

	s.AddUser(user)
	defer s.DeleteUser(user)

	// Send the user a list of other users
	usersMsg := &v1.Message{
		Body:    "users",
		Content: others,
	}
	websocket.JSON.Send(ws, usersMsg)

	// Send the initial state of the resource list
	listMsg := &v1.Message{
		Body:    "init",
		Content: s.store.List(),
	}
	websocket.JSON.Send(ws, listMsg)

	// Report the error that ended the main event loop
	var err error

	// Main event loop
Events:
	for {
		var msg v1.ThingMessage
		err = websocket.JSON.Receive(ws, &msg)
		if err != nil {
			break Events
		}

		log.Printf("Message from %s: %s\n", user, msg)
		s.HandleMessage(&msg)

	}
	log.Printf("User %s exited with error %s\n", user, err)
}

// New creates a new server. It will panic on error
func New(config config.Config, store Store) *Server {
	srv := &Server{
		config: config,
		store:  store,
		templates: templates.New(
			config.TemplateDir,
			templates.Attrs{"StaticURL": config.StaticURL},
		),
	}

	// Make the user dictionary
	srv.users = make(map[*User]bool)

	// Routes
	http.HandleFunc("/", srv.IndexHandler)
	http.Handle("/events", websocket.Handler(srv.EventsHandler))
	http.Handle(
		config.StaticURL,
		http.StripPrefix(
			config.StaticURL,
			http.FileServer(http.Dir(config.StaticDir)),
		),
	)
	return srv
}
