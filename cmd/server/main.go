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

	//Initialize Handlers
	adminHandler := handler.NewAdminHandler(userService, scheduleService, proxyEngine)
	teacherHandler := handler.NewTeacherHandler(teacherService)

	mux := http.NewServeMux()

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

	// --- Protected Admin Routes (Requires "ADMIN" Role) ---

	// 1. Admin Onboarding UI Page
	mux.Handle("/admin/onboarding", middleware.RequireRole("ADMIN")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("internal/web/templates/admin/onboarding.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	})))

	// 2. Admin Schedule UI Page
	mux.Handle("/admin/schedule", middleware.RequireRole("ADMIN")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("internal/web/templates/admin/schedule.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	})))

	// 3. Admin Leaves UI Page
	mux.Handle("/admin/leaves", middleware.RequireRole("ADMIN")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("internal/web/templates/admin/leaves.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	})))

	// --- Teacher Portal UI Pages (GET) ---
	mux.Handle("/teacher/assignments", middleware.RequireRole("TEACHER", "ADMIN")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("internal/web/templates/teacher/assignments.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	})))

	mux.Handle("/teacher/attendance", middleware.RequireRole("TEACHER", "ADMIN")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("internal/web/templates/teacher/attendance.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	})))

	// --- Teacher Form Submission Endpoints (POST) ---
	// Note: Make sure you initialize teacherHandler with your TeacherService 
	// and register these API endpoints to point to teacherHandler methods.
	mux.Handle("/api/teacher/assignments", middleware.RequireRole("TEACHER", "ADMIN")(http.HandlerFunc(teacherHandler.HandleCreateAssignment)))
	mux.Handle("/api/teacher/attendance", middleware.RequireRole("TEACHER", "ADMIN")(http.HandlerFunc(teacherHandler.HandleMarkAttendance)))


	
	// 4. Admin Form Submission Endpoints (POST)
	mux.Handle("/api/admin/onboarding", middleware.RequireRole("ADMIN")(http.HandlerFunc(adminHandler.HandleOnboarding)))
	mux.Handle("/api/admin/schedules", middleware.RequireRole("ADMIN")(http.HandlerFunc(adminHandler.HandleScheduleCreation)))
	mux.Handle("/api/admin/leaves", middleware.RequireRole("ADMIN")(http.HandlerFunc(adminHandler.HandleLeaveSubmission)))

	log.Printf("School Management System initialized successfully! Starting server on %s...", cfg.ServerAddress)
	if err := http.ListenAndServe(cfg.ServerAddress, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

