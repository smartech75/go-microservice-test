package handlers

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// DotEnvVariable -> get .env
func DotEnvVariable(key string) string {

	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}
