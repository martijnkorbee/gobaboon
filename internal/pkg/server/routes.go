package server

import (
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog"
	"github.com/martijnkorbee/gobaboon/pkg/logger"
	"github.com/rs/cors"
)

func (s *Server) routes(log *logger.Logger) *chi.Mux {
	r := chi.NewRouter()

	// default go-chi middleware
	r.Use(middleware.RealIP)
	r.Use(httplog.RequestLogger(*log.Logger))

	// Basic CORS
	// for more info, see: https://github.com/rs/cors
	r.Use(cors.New(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}).Handler)

	// added default middleware
	r.Use(s.Session.Manager.LoadAndSave)
	r.Use(s.Middleware.NoSurf) // CSRF protection
	r.Use(s.Middleware.CheckMaintenanceMode)

	return r
}
