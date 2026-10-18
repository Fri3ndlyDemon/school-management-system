package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// all config values needed for the application
type Config struct {
	ServerAddress string
	DatabaseURL   string
}

// load environment variables from .env otherwise system environment variables
func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Note: No .env file found. Falling back to system environment variables.")
	}

	serverAddr := os.Getenv("SERVER_ADDRESS")
	if serverAddr == "" {
		serverAddr = ":8080"
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Println("Warning: DATABASE_URL is not set in environment or .env file.")
	}

	return &Config{
		ServerAddress: serverAddr,
		DatabaseURL:   dbURL,
	}
}
