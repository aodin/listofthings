package server

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

type User struct {
	conn *websocket.Conn
	Id   int64
	Name string
}

// Wrap HTTP methods
type Server struct {
	config    Config
	templates map[string]*template.Template
	users     map[*User]bool
	counter   int64
	store     Storage
}

func (s *Server) ListenAndServe() error {
	address := fmt.Sprintf(":%d", s.config.Port)
	return http.ListenAndServe(address, nil)
}

// User Events
// -----------
func (s *Server) AddUser(user *User) {
	// Broadcast a user joined message to all users
	msg := &Message{
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

	msg := &Message{
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

func (s *Server) BroadcastMessage(msg *Message) {
	// Send the message to all users
	log.Println("Broadcasting to users:", len(s.users))
	for user, _ := range s.users {
		err := websocket.JSON.Send(user.conn, msg)
		log.Println("Broadcast message sent:", err)
		if err != nil {
			continue
		}
	}
}

func (s *Server) HandleThingMessage(msg *ThingMessage) {
	var err error
	var thing *Thing

	switch msg.Method {
	case "create":
		log.Println("Creating thing", msg.Item)
		thing, err = s.store.Create(msg.Item)
	case "delete":
		log.Println("Deleting thing", msg.Item)
		thing, err = s.store.Delete(msg.Item)
	case "update":
		log.Println("Updating thing", msg.Item)
		thing, err = s.store.Update(msg.Item)
	default:
		log.Println("Unknown method:", msg.Method)
		err = fmt.Errorf("Unknown method '%s'", msg.Method)
	}

	// Build the return message
	var returnMsg Message
	if err != nil {
		returnMsg.Body = "error"
		returnMsg.Content = err.Error()
	} else {
		returnMsg.Body = msg.Method
		returnMsg.Content = thing
	}
	s.BroadcastMessage(&returnMsg)
}

// Serve the list template
func (s *Server) RootHandler(w http.ResponseWriter, r *http.Request) {
	list, ok := s.templates["list"]
	if !ok {
		http.Error(w, "Template not found", 500)
		return
	}

	// TODO Output all things?
	attrs := map[string]interface{}{}

	list.Execute(w, attrs)
}

func (s *Server) EventsHandler(ws *websocket.Conn) {
	// Register the new connection
	s.counter += 1

	// Create a user object and join the server
	user := &User{Id: s.counter, conn: ws}
	s.users[user] = true
	log.Println("User connected:", user)

	things := s.store.List()

	listMsg := &Message{
		Body:    "init",
		Content: things,
	}

	// Send the current list
	websocket.JSON.Send(ws, listMsg)

	var err error

	// Main event loop
Events:
	for {
		var msg ThingMessage
		err = websocket.JSON.Receive(ws, &msg)
		if err != nil {
			break Events
		}

		log.Println("New message:", msg)
		s.HandleThingMessage(&msg)

		// TODO What to do with various message types?
		// s.BroadcastMessage(user, &msg)
	}
	log.Println("Exiting event loop:", err)
	// TODO This could be a defer method
	s.DeleteUser(user)
}

func New(config Config, store Storage) (*Server, error) {
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
	staticURL := "/static/"
	http.Handle(
		staticURL,
		http.StripPrefix(staticURL, http.FileServer(http.Dir(config.Static))),
	)
	return s, nil
}
