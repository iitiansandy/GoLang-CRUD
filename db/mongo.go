package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

// Connect connects to MongoDB and assigns the global Client
func Connect(uri string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return err
	}

	// Ping to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		return err
	}

	Client = client
	log.Println("Connected to MongoDB")
	return nil
}
