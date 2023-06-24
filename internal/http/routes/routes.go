package routes

import (
	"baboonapp/http/handlers"
	"baboonapp/http/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type AppRoutes struct {
	Middleware middleware.Middleware
	Handlers   handlers.Handlers
}

func (ar *AppRoutes) Routes() *chi.Mux {
	// application routes
	r := chi.NewRouter()

	// add middleware

	// add your routes here
	r.Get("/", ar.Handlers.Home) // default home route

	// static routes
	fileServer := http.FileServer(http.Dir("./public"))
	r.Handle("/public/*", http.StripPrefix("/public", fileServer))

	return r
}
