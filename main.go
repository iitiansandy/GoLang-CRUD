package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/coddicted/go-mongo-crud/db"
	"github.com/coddicted/go-mongo-crud/routes"
	"github.com/joho/godotenv"
)

func main() {
	// Load env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found — reading environment variables directly")
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI is required")
	}
	if err := db.Connect(mongoURI); err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := routes.Register()
	fmt.Printf("Server listening on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
