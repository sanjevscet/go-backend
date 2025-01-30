package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/sanjevscet/go-backend.git/internal/db"
	"github.com/sanjevscet/go-backend.git/internal/env"
	"github.com/sanjevscet/go-backend.git/internal/store"
)

func main() {

	// Load .envrc file
	err := godotenv.Load(".envrc")
	if err != nil {
		log.Fatalf("Error loading .envrc file")
	}
	cfg := config{
		addr: env.GetString("ADDR", ":1414"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://sanjeev:sanjeev@localhost:11432/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 10),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
	}

	database, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime)
	if err != nil {
		log.Panic()
	}
	defer database.Close()
	log.Println("<<<<<<<<<<<<< DB Connected >>>>>>>>>>>>>>>")

	store := store.NewStorage(database)

	app := &application{
		config: cfg,
		store:  store,
	}

	log.Printf("Starting server on %s", app.config.addr)
	mux := app.mount()

	log.Fatal(app.run(mux))
	// Start the server
}
