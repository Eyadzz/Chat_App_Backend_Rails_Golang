package infrastructure

import (
	"github.com/joho/godotenv"
	"log"
)

func LoadEnvironmentVariables() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
