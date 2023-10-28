package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"playground/api/restAPI"
	"playground/utils/database"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	uriMongodb := os.Getenv("MONGO_URI")
	port := os.Getenv("PORT")
	// Connect to MongoDB
	database.ConnectMongoDB(uriMongodb)
	// Initialize Gin Rest API
	restAPI.InitializeGinRestAPI(port)

	defer database.DisconnectMongo(database.MongoClient)
}
