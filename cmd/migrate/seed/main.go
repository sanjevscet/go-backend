package main

import (
	"log"

	"github.com/sanjevscet/go-backend.git/internal/db"
	"github.com/sanjevscet/go-backend.git/internal/store"
)

func main() {
	conn, err := db.New("postgres://sanjeev:sanjeev@localhost:11432/social?sslmode=disable", 3, 3, "15m")
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	store := store.NewStorage(conn)

	db.Seed(store, conn)
}
