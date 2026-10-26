package main

import (
	"html/template"
	"log"
	"net/http"

	"school-management-system/internal/config"
	"school-management-system/internal/database"
	"school-management-system/internal/middleware"
	"school-management-system/internal/repository"
	"school-management-system/internal/service"
	"school-management-system/internal/web/handler"
)

// Helper function to render HTML templates cleanly
func page(templatePath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("internal/web/templates/" + templatePath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	}
}

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
	teacherRepo := repository.NewTeacherRepository(db)

	// Initialize Services
	authService := service.NewAuthService(cfg.JWTSecret)
	authHandler := handler.NewAuthHandler(authService)

	userService := service.NewUserService(userRepo)
	scheduleService := service.NewScheduleService(scheduleRepo)
	proxyEngine := service.NewProxyEngine(db)
	teacherService := service.NewTeacherService(teacherRepo)

	// Initialize Handlers
	adminHandler := handler.NewAdminHandler(userService, scheduleService, proxyEngine)
	teacherHandler := handler.NewTeacherHandler(teacherService)

	// Setup Authentication and Route Guard Helpers
	auth := middleware.AuthMiddleware(authService)
	guard := func(h http.HandlerFunc, roles ...string) http.Handler {
		return auth(middleware.RequireRole(roles...)(h))
	}

	mux := http.NewServeMux()

	// --- Static Assets (CSS, JS) ---
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("internal/web/static"))))

	// --- Public Routes ---
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("internal/web/templates/auth/login.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	})

	mux.HandleFunc("/api/auth/register", authHandler.Register)
	mux.HandleFunc("/api/auth/login", authHandler.Login)

	// --- Protected Admin Routes ---
	mux.Handle("/admin/onboarding", guard(page("admin/onboarding.html"), "admin"))
	mux.Handle("/admin/schedule", guard(page("admin/schedule.html"), "admin"))
	mux.Handle("/admin/leaves", guard(page("admin/leaves.html"), "admin"))

	mux.Handle("/api/admin/onboarding", guard(adminHandler.HandleOnboarding, "admin"))
	mux.Handle("/api/admin/schedules", guard(adminHandler.HandleScheduleCreation, "admin"))
	mux.Handle("/api/admin/leaves", guard(adminHandler.HandleLeaveSubmission, "admin"))

	// --- Protected Teacher Portal Routes ---
	mux.Handle("/teacher/assignments", guard(page("teacher/assignments.html"), "teacher", "admin"))
	mux.Handle("/teacher/attendance", guard(page("teacher/attendance.html"), "teacher", "admin"))

	mux.Handle("/api/teacher/assignments", guard(teacherHandler.HandleCreateAssignment, "teacher", "admin"))
	mux.Handle("/api/teacher/attendance", guard(teacherHandler.HandleMarkAttendance, "teacher", "admin"))

	log.Printf("School Management System initialized successfully! Starting server on %s...", cfg.ServerAddress)
	if err := http.ListenAndServe(cfg.ServerAddress, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}