package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// GetValue returns value from the .env file
func GetValue(key string) string {
	// load the .env file
	var err error = godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error while loading .env file\n")
	}
	// return the value from the .env file
	// based on the provided key
	return os.Getenv(key)
}
