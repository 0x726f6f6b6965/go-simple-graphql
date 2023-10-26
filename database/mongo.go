package database

import (
	"context"
	"time"

	"github.com/0x726f6f6b6965/go-simple-graphql/utils"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoInstance struct {
	// Client represents MongoDB client
	Client *mongo.Client
	// Database represents MongoDB database
	Database *mongo.Database
}

var Mongo MongoInstance

// Connect connects an application to the MongoDB database
func Connect(dbName string) error {
	// create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// create a new MongoDB client
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(utils.GetValue("MONGO_URI")))

	// if connection fails, return an error
	if err != nil {
		return err
	}

	// get the MongoDB database
	var db *mongo.Database = client.Database(dbName)

	// assign the MongoDB client and the MongoDB database
	Mongo = MongoInstance{
		Client:   client,
		Database: db,
	}

	return nil

}

// GetCollection returns collection based on the given name
func GetCollection(name string) *mongo.Collection {
	// return the collection from the MongoDB database
	return Mongo.Database.Collection(name)
}
