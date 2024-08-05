package main

import (
	"backend/internal/config"
	"backend/internal/database"
	"backend/internal/routes"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	config.LoadConfig()
	database.InitZincSearch()

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	routes.InitRoutes(r)

	log.Print("Starting server on :8080")
	http.ListenAndServe(":8080", r)
}
