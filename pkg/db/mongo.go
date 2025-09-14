package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoInstance struct {
	Client   *mongo.Client
	Database *mongo.Database

	Users *mongo.Collection
}

func Connect(uri, dbName string) *MongoInstance {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("Mongo connection error:", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("Mongo ping failed:", err)
	}

	log.Println("Mongodb connected and reachable")

	database := client.Database(dbName)

	return &MongoInstance{
		Client:   client,
		Database: database,
		Users:    database.Collection("users"),
	}
}
