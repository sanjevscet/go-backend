package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/sanjevscet/go-backend.git/internal/db"
	"github.com/sanjevscet/go-backend.git/internal/env"
	"github.com/sanjevscet/go-backend.git/internal/store"
)

const API_VERSION = "0.0.1"

func main() {

	// Load .envrc file
	err := godotenv.Load(".envrc")
	if err != nil {
		log.Fatalf("Error loading .envrc file")
	}
	dbConfig := dbConfig{
		addr:         env.GetString("DB_ADDR", "postgres://sanjeev:sanjeev@localhost:11432/social?sslmode=disable"),
		maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
		maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 10),
		maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
	}

	log.Printf("DB_ADDR: %s", dbConfig.addr)

	cfg := config{
		addr: env.GetString("ADDR", ":1414"),
		env:  env.GetString("ENV", "development"),
		db:   dbConfig,
	}

	database, err := db.New(
		dbConfig.addr,
		dbConfig.maxOpenConns,
		dbConfig.maxIdleConns,
		dbConfig.maxIdleTime)
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
