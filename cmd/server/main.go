package main

import (
	"html/template"
	"log"
	"net/http"

	"school-management-system/internal/config"
	"school-management-system/internal/database"
	"school-management-system/internal/repository"
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

	// Initialize Repositories
	userRepo := repository.NewUserRepository(db)
	scheduleRepo := repository.NewScheduleRepository(db)

	// Initialize Services & Engines
	authService := service.NewAuthService(cfg.JWTSecret)
	authHandler := service.NewAuthHandler(authService)

	userService := service.NewUserService(userRepo)
	scheduleService := service.NewScheduleService(scheduleRepo)
	proxyEngine := service.NewProxyEngine(db)

	adminHandler := service.NewAdminHandler(userService, scheduleService, proxyEngine)

	mux := http.NewServeMux()

	// --- UI Page Routes (GET) ---
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("internal/web/templates/auth/login.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	})

	mux.HandleFunc("/admin/onboarding", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("internal/web/templates/admin/onboarding.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	})

	mux.HandleFunc("/admin/schedule", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("internal/web/templates/admin/schedule.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	})

	mux.HandleFunc("/admin/leaves", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("internal/web/templates/admin/leaves.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	})

	// --- API & Form Submission Routes (POST) ---
	mux.HandleFunc("/api/auth/register", authHandler.Register)
	mux.HandleFunc("/api/auth/login", authHandler.Login)
	mux.HandleFunc("/api/admin/onboarding", adminHandler.HandleOnboarding)
	mux.HandleFunc("/api/admin/schedules", adminHandler.HandleScheduleCreation)
	mux.HandleFunc("/api/admin/leaves", adminHandler.HandleLeaveSubmission)

	log.Printf("School Management System initialized successfully! Starting server on %s...", cfg.ServerAddress)
	if err := http.ListenAndServe(cfg.ServerAddress, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}