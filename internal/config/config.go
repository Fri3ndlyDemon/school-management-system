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
	JWTSecret     string
}

// load environment variables from .env otherwise system environment variables
func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Note: No .env file found. Falling back to system environment variables.")
	}

	serverAddr := os.Getenv("SERVER_ADDRESS")
	if serverAddr == "" {
		log.Fatalf("SERVER_ADDRESS is not set in environment or .env file, terminating the application.")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatalf("DATABASE_URL is not set in environment or .env file, terminating the application.")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatalf("JWT_SECRET is not set in environment or .env file, terminating the application.")
	}

	return &Config{
		ServerAddress: serverAddr,
		DatabaseURL:   dbURL,
		JWTSecret:     jwtSecret,
	}
}
