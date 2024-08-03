package database

import (
	"log"
	"os"
)

var zincBaseURL string
var zincUsername string
var zincPassword string

func InitZincSearch() {
	zincBaseURL = os.Getenv("ZINC_BASE_URL")
	zincUsername = os.Getenv("ZINC_USERNAME")
	zincPassword = os.Getenv("ZINC_PASSWORD")

	if zincBaseURL == "" || zincUsername == "" || zincPassword == "" {
		log.Fatal("ZincSearch configuration is missing")
	}
}
