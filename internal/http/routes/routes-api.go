package routes

import (
	"github.com/go-chi/chi/v5"
)

func (ar *AppRoutes) RoutesAPI() *chi.Mux {
	// API routes
	r := chi.NewRouter()

	// used to check health
	r.Get("/ping", ar.Handlers.Ping) // default route

	// add your middleware and routes here

	return r
}
