package main

import (
	"log"
	"net/http"

	"school-management-system/internal/config"
	"school-management-system/internal/database"
	"school-management-system/internal/service"
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

	// initialize auth service and handler
	authService := service.NewAuthService(cfg.JWTSecret)
	authHandler := service.NewAuthHandler(authService)

	mux := http.NewServeMux()

	// public endpoints
	mux.HandleFunc("/api/auth/register", authHandler.Register)
	mux.HandleFunc("/api/auth/login", authHandler.Login)

	// start the http server
	log.Printf("School Management System initialized successfully! Starting server on %s...", cfg.ServerAddress)
	if err := http.ListenAndServe(cfg.ServerAddress, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
