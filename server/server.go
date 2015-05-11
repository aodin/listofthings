package server

import (
	"net/http"

	"github.com/aodin/volta/config"
	"github.com/aodin/volta/templates"

	"github.com/aodin/listofthings/server/auth"
	// "github.com/aodin/listofthings/server/feeds/v1"
)

// Wrap HTTP methods
type Server struct {
	config    config.Config
	sessions  *auth.SessionManager
	templates *templates.Templates
	users     *auth.UserManager
}

func (srv *Server) RequireSession(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie(srv.config.Cookie.Name)
		if err != nil {
			auth.SetCookie(w, srv.config.Cookie, srv.sessions.Create())
		} else {
			// TODO If the cookie value is invalid, reset it
		}
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
func New(config config.Config) *Server {
	srv := &Server{
		config:   config,
		sessions: auth.Sessions(config),
		templates: templates.New(
			config.TemplateDir,
			templates.Attrs{"StaticURL": config.StaticURL},
		),
		users: auth.Users(),
	}

	// Routes
	http.HandleFunc("/", srv.RequireSession(srv.IndexHandler))
	http.Handle(
		config.StaticURL,
		http.StripPrefix(
			config.StaticURL,
			http.FileServer(http.Dir(config.StaticDir)),
		),
	)
	return srv
}
