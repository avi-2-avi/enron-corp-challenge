package main

import (
	"backend/internal/config"
	"backend/internal/database"
	"backend/internal/routes"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	config.LoadConfig()
	database.InitZincSearch()

	frontendBaseUrl := os.Getenv("FRONTEND_BASE_URL")
	if frontendBaseUrl == "" {
		frontendBaseUrl = "http://localhost:5173"
		log.Println("FRONTEND_BASE_URL not set, defaulting to", frontendBaseUrl)
	}

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{frontendBaseUrl},
		AllowedMethods:   []string{"GET"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	routes.InitRoutes(r)

	log.Print("Starting server on :8080")
	http.ListenAndServe(":8080", r)
}
