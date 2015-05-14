package server

import (
	"net/http"

	"code.google.com/p/go.net/websocket"
	sql "github.com/aodin/aspect"
	"github.com/aodin/volta/config"
	"github.com/aodin/volta/templates"

	db "github.com/aodin/listofthings/db"
	"github.com/aodin/listofthings/server/auth"
	feeds "github.com/aodin/listofthings/server/feeds/v1"
)

// Wrap HTTP methods
type Server struct {
	config    config.Config
	sessions  *auth.SessionManager
	templates *templates.Templates
	users     *auth.UserManager
}

// TODO auth function?
func (srv *Server) RequireSession(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var session db.Session
		if cookie, err := r.Cookie(srv.config.Cookie.Name); err == nil {
			session = srv.sessions.Get(cookie.Value)
		}

		// If the cookie value is invalid, create a new user and session
		if !session.Exists() {
			user := srv.users.Create("", "")
			auth.SetCookie(w, srv.config.Cookie, srv.sessions.Create(user))
		}

		// Call the wrapped handler
		f(w, r)
	}
}

func (srv *Server) ListenAndServe() error {
	return http.ListenAndServe(srv.config.Address(), nil)
}

// Index is the handler for the index
func (srv *Server) IndexHandler(w http.ResponseWriter, r *http.Request) {
	// Assign a session if one has not been set
	srv.templates.Execute(w, "index")
}

// New creates a new server. It will panic on error
func New(config config.Config, conn sql.Connection) *Server {
	srv := &Server{
		config:   config,
		sessions: auth.Sessions(config, conn),
		templates: templates.New(
			config.TemplateDir,
			templates.Attrs{"StaticURL": config.StaticURL},
		),
		users: auth.Users(conn),
	}

	// Routes
	http.HandleFunc("/", srv.RequireSession(srv.IndexHandler))

	// Feeds
	hub := feeds.NewHub(srv.users, srv.sessions)
	http.Handle("/feeds/v1/things", websocket.Handler(hub.Handler))

	// Static Files
	http.Handle(
		config.StaticURL,
		http.StripPrefix(
			config.StaticURL,
			http.FileServer(http.Dir(config.StaticDir)),
		),
	)
	return srv
}
