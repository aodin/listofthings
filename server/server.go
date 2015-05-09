package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"code.google.com/p/go.net/websocket"

	"github.com/aodin/listofthings/db"
	"github.com/aodin/listofthings/server/config"
	"github.com/aodin/listofthings/server/feeds/v1"
)

type User struct {
	conn *websocket.Conn
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (u *User) String() string {
	if u.Name == "" {
		return fmt.Sprintf("%d", u.ID)
	}
	return fmt.Sprintf("%s (%d)", u.Name, u.ID)
}

// Wrap HTTP methods
type Server struct {
	config    config.Config
	templates map[string]*template.Template
	users     map[*User]bool
	counter   int64
	store     Store
}

func (s *Server) ListenAndServe() error {
	address := fmt.Sprintf(":%d", s.config.Port)
	return http.ListenAndServe(address, nil)
}

// Serve the list template
func (s *Server) RootHandler(w http.ResponseWriter, r *http.Request) {
	list, ok := s.templates["list"]
	if !ok {
		http.Error(w, "Template not found", 500)
		return
	}
	// TODO Bootstrap the things using JSON embeded in the template?
	attrs := map[string]interface{}{}
	list.Execute(w, attrs)
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

func New(config config.Config, store Store) (*Server, error) {
	// Parse the templates
	tmplPath := filepath.Join(config.Templates, "list.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return nil, err
	}
	templates := map[string]*template.Template{"list": tmpl}

	// Build the server
	s := &Server{config: config, templates: templates, store: store}

	// Make the user dictionary
	s.users = make(map[*User]bool)

	// Bind the root HTTP hander
	http.HandleFunc("/", s.RootHandler)

	// And the websocket handler
	http.Handle("/events", websocket.Handler(s.EventsHandler))

	// Serve the static files
	staticURL := "/dist/"
	http.Handle(
		staticURL,
		http.StripPrefix(staticURL, http.FileServer(http.Dir(config.Static))),
	)
	return s, nil
}
