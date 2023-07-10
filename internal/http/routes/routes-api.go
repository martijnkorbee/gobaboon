package routes

import (
	"github.com/go-chi/chi/v5"
)

func (ar *AppRoutes) RoutesAPI() *chi.Mux {
	// API routes
	r := chi.NewRouter()

	// add your middleware here

	// used to check health
	r.Get("/ping", ar.Handlers.Ping) // default route

	// users
	r.Route("/users/{email}", func(r chi.Router) {
		r.MethodFunc("POST", "/", ar.Handlers.UsersAdd)
		r.MethodFunc("DELETE", "/", ar.Handlers.UsersDelete)
	})

	r.Route("/v1", func(r chi.Router) {
		// these routes are protected
		r.Use(ar.Middleware.AuthToken)

		// used to check health
		r.Get("/ping", ar.Handlers.Ping) // default route
	})

	return r
}
