package routes

import (
	"backend/internal/handlers"

	"github.com/go-chi/chi/v5"
)

func InitRoutes(r *chi.Mux) {
	r.Get("/health", handlers.HealthCheck)
}
