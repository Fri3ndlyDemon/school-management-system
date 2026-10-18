package main

import (
	"log"

	"school-management-system/internal/config"
	"school-management-system/internal/database"
)

func main() {
	// load config from env
	cfg := config.LoadConfig()
	log.Printf("Configuration loaded successfully. Target address: %s", cfg.ServerAddress)

	// initialize connection to database
	db, err := database.ConnectDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Critical error: Failed to connect to the database: %v", err)
	}

	// close database when main exits
	defer db.Close()

	log.Println("School Management System initialized successfully!")
}
