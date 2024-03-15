package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	middlewares "github.com/smartech75/go-microservice-test/user-service/handlers"
)

var client *mongo.Client

func DotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	fmt.Println("Getting value: " + key + ": " + os.Getenv(key))

	return os.Getenv(key)
}

// Dbconnect -> connects mongo
func Dbconnect() *mongo.Client {

	// // Mongodb connect for local server
	clientOptions := options.Client().ApplyURI(middlewares.DotEnvVariable("MONGO_URL"))

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal("⛒ Connection Failed to Database")
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal("⛒ Connection Failed to Database")
		log.Fatal(err)
	}
	color.Green("⛁ Connected to Database")
	return client
}
