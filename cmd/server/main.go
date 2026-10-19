package main

import (
	"html/template"
	"log"
	"net/http"

	"school-management-system/internal/config"
	"school-management-system/internal/database"
	"school-management-system/internal/service"
)

func main() {
	cfg := config.LoadConfig()
	log.Printf("Configuration loaded successfully. Target address: %s", cfg.ServerAddress)

	db, err := database.ConnectDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Critical error: Failed to connect to the database: %v", err)
	}
	defer db.Close()

	authService := service.NewAuthService(cfg.JWTSecret)
	authHandler := service.NewAuthHandler(authService)

	mux := http.NewServeMux()

	// show login page 
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("internal/web/templates/auth/login.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	})

	// server admin onboarding
	mux.HandleFunc("/admin/onboarding", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("internal/web/templates/admin/onboarding.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	})

	// server admin schedule page
	mux.HandleFunc("/admin/schedule", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("internal/web/templates/admin/schedule.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	})

	// server admin leaves page
	mux.HandleFunc("/admin/leaves", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("internal/web/templates/admin/leaves.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	})

	// API endpoints
	mux.HandleFunc("/api/auth/register", authHandler.Register)
	mux.HandleFunc("/api/auth/login", authHandler.Login)

	log.Printf("School Management System initialized successfully! Starting server on %s...", cfg.ServerAddress)
	if err := http.ListenAndServe(cfg.ServerAddress, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}