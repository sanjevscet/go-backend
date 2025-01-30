package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/sanjevscet/go-backend.git/internal/env"
)

func main() {

	// Load .envrc file
	err := godotenv.Load(".envrc")
	if err != nil {
		log.Fatalf("Error loading .envrc file")
	}
	cfg := config{
		addr: env.GetString("ADDR", ":1414"),
	}
	app := &application{
		config: cfg,
	}

	log.Printf("Starting server on %s", app.config.addr)
	mux := app.mount()

	log.Fatal(app.run(mux))
	// Start the server
}
