package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"dstartest/internal/handlers"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Route("/", func(r chi.Router) {
		r.Use(app.sessions.LoadAndSave) // Session middleware
		r.Get("/", handlers.Root("landing-page"))
		r.Get("/landing-page", handlers.LandingPage())
		r.Route("/user", func(r chi.Router) {
			r.Get("/profile", handlers.UserProfile())
		})
	})
	r.Get("/healthz", handlers.Health(app.ready))
	return r
}
