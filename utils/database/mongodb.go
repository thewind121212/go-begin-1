package database

// MongoDB

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

func ConnectMongoDB(uriMongodb string) {

	ctx := context.Background()
	clientOptions := options.Client().ApplyURI(uriMongodb)
	clientOptions.SetMaxConnecting(10)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		fmt.Println("Error to connect MongoDB", uriMongodb)
		defer DisconnectMongo(client)
	}

	//checking connection to mongodb

	err = client.Ping(ctx, nil)
	if err != nil {
		fmt.Println("Error to connect MongoDB")
		defer DisconnectMongo(client)
	}

	fmt.Println("Connected to MongoDB")

	MongoClient = client
}

func DisconnectMongo(client *mongo.Client) {
	fmt.Println("Disconnecting MongoDB")
	err := client.Disconnect(context.Background())
	if err != nil {
		panic(err)
	}
}
