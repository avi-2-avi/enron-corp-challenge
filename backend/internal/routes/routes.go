package routes

import (
	"backend/internal/handlers"

	"github.com/go-chi/chi/v5"
)

func InitRoutes(r *chi.Mux) {
	r.Route("/emails", func(r chi.Router) {
		r.Get("/", handlers.GetEmails)
		r.Get("/{id}", handlers.GetEmail)
	})

	r.Get("/health", handlers.HealthCheck)
}
