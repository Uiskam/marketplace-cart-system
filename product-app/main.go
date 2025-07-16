package main

import (
	"log"
	"product-app/app/server"
	"product-app/external"
)

func main() {
	db, err := external.NewPostgresConnection()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	srv := server.NewServer(db)

	log.Println("Starting server on :8080")
	if err := srv.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
