package handler

import (
	"html/template"
	"net/http"
	"school-management-system/internal/middleware"
	"school-management-system/internal/service"
)

type DashboardHandler struct {
	UserService *service.UserService
}

func NewDashboardHandler(userService *service.UserService) *DashboardHandler {
	return &DashboardHandler{
		UserService: userService,
	}
}

func (h *DashboardHandler) RenderDashboard(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// extract name and role from claims safely
	name := "User"
	if n, ok := claims["name"].(string); ok {
		name = n
	}

	role, _ := claims["role"].(string)

	data := struct {
		Name              string
		Role              string
		Success           bool
		ExamResults       interface{}
		AttendanceRecords interface{}
		Schedules         interface{}
	}{
		Name:    name,
		Role:    role,
		Success: r.URL.Query().Get("success") == "true",
	}

	// Parse all sub-template files concurrently so Go resolves defined blocks correctly
	tmpl, err := template.ParseFiles(
		"internal/web/templates/dashboard.html",
		"internal/web/templates/admin/onboarding.html",
		"internal/web/templates/admin/schedule.html",
		"internal/web/templates/admin/leaves.html",
		"internal/web/templates/admin/announcement.html",
		"internal/web/templates/teacher/assignments.html",
		"internal/web/templates/teacher/attendance.html",
		"internal/web/templates/student/results.html",
		"internal/web/templates/student/attendance.html",
		"internal/web/templates/student/calendar.html",
	)
	if err != nil {
		http.Error(w, "Template parsing error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Template execution error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
