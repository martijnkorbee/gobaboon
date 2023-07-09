// Server package is responsible for creating the server for the baboon application.
//
// Uses: go-chi as router: https://go-chi.io/,
// alex-edwards scs v2 as session manager: http://github.com/alexedwards/scs/v2
package server

import (
	"errors"
	"fmt"
	"github.com/martijnkorbee/gobaboon/pkg/render"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/martijnkorbee/gobaboon/pkg/logger"
)

var (
	// MaintenanceMode is used to set the server in maintenance mode (default false)
	MaintenanceMode bool = false
)

type Server struct {
	// Config holds all configuration variables for the server
	config ServerConfig

	// Log is the server's logger
	Log *logger.Logger

	// Session holds the server's sessionmanager
	Session *Session

	// Renderer is used to render html pages for app routes
	Renderer *render.Renderer

	// Middleware holds the default middleware added to a baboon server
	Middleware *Middleware

	// Router is the baboon server's router
	Router *chi.Mux
}

// ServerConfig stores all configuration variables for a baboon server
//
// TODO: find a way to have a generic cache interface for baboon application
type ServerConfig struct {
	// this is a required field
	// rootpath is the rootpath of the application using the baboon server,
	Rootpath string

	// Debug defaults to false
	Debug bool

	// host defaults to localhost
	Host string

	// port defaults to 4000
	Port string

	// Renderer defaults to the jet templating engine
	Renderer string

	// Cookie cookie configuration
	Cookie CookieConfig
}

// NewServer returns a new server loaded with the server config.properties.
func NewServer(sc ServerConfig) (*Server, error) {
	// check for rootpath
	if sc.Rootpath == "" {
		return nil, errors.New("rootpath is empty, required")
	}

	// set default variables
	if sc.Host == "" {
		sc.Host = "localhost"
	}
	if sc.Port == "" {
		sc.Port = "4000"
	}
	if sc.Renderer == "" {
		sc.Renderer = "jet"
	}

	// create server
	srv := Server{
		config: sc,
	}

	// create logger
	log := &logger.LoggerConfig{
		Rootpath: srv.config.Rootpath,
		Debug:    srv.config.Debug,
		Console:  true,
		Service:  "webserver",
	}
	// start logger
	srv.Log = log.Start()

	// add renderer (call before middleware)
	srv.Renderer = render.New(srv.config.Rootpath, srv.config.Renderer, srv.config.Debug)

	// assign session manager (call before middleware)
	srv.Session = NewSession(srv.config.Cookie, &srv)

	// add middleware (call before routes)
	srv.Middleware = &Middleware{
		roopath:      srv.config.Rootpath,
		nosurfCookie: NewCookie(srv.config.Cookie),
	}

	// add routes
	srv.Router = srv.routes(srv.Log)

	return &srv, nil
}

func (s *Server) Run() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", s.config.Host, s.config.Port),
		Handler:      s.Router,
		IdleTimeout:  30 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 600 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

// ActivateMaintenanceMode is replaced by SetMaintenanceMode. ActivateMaintenceMode calls SetMaintenanceMode true.
func (s *Server) ActivateMaintenanceMode() {
	s.SetMaintenanceMode(true)
}

// SetMaintenaceMode sets the server in or out of maintenance mode.
func (s *Server) SetMaintenanceMode(on bool) {
	MaintenanceMode = on
}
