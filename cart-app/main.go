package main

import (
	"cart-app/app/server"
	"cart-app/external"
	"log"
)

func main() {
	db, err := external.NewPostgresConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	srv := server.NewServer(db)
	log.Println("Starting cart API server on :8081")
	if err := srv.Run(":8081"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}