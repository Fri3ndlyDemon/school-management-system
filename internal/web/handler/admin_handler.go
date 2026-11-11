package handler

import (
	"database/sql"
	"net/http"
	"school-management-system/internal/middleware"
	"school-management-system/internal/service"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type AdminHandler struct {
	userService     *service.UserService
	scheduleService *service.ScheduleService
	proxyEngine     *service.ProxyEngine
	db              *sql.DB
}

func NewAdminHandler(userService *service.UserService, scheduleService *service.ScheduleService, proxyEngine *service.ProxyEngine, db *sql.DB) *AdminHandler {
	return &AdminHandler{
		userService:     userService,
		scheduleService: scheduleService,
		proxyEngine:     proxyEngine,
		db:              db,
	}
}

// Handle Onboarding Form Submission (calls OnboardUser)
func (h *AdminHandler) HandleOnboarding(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")
	role := r.FormValue("role")

	// Matches UserService.OnboardUser
	err := h.userService.OnboardUser(name, email, password, role)
	if err != nil {
		http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

// Handle Master Schedule Form Submission (calls SetMasterSchedule)
func (h *AdminHandler) HandleScheduleCreation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	courseID, _ := strconv.Atoi(r.FormValue("course_id"))
	teacherID, _ := strconv.Atoi(r.FormValue("teacher_id"))
	dayOfWeek := r.FormValue("day_of_week")
	startTime := r.FormValue("start_time")
	endTime := r.FormValue("end_time")

	// Parse comma-separated proxy pool IDs (e.g. "12, 14, 19")
	rawProxyPool := r.FormValue("proxy_pool")
	var proxyPool []int
	if rawProxyPool != "" {
		parts := strings.Split(rawProxyPool, ",")
		for _, p := range parts {
			id, err := strconv.Atoi(strings.TrimSpace(p))
			if err == nil {
				proxyPool = append(proxyPool, id)
			}
		}
	}

	// Matches ScheduleService.SetMasterSchedule
	err := h.scheduleService.SetMasterSchedule(courseID, teacherID, proxyPool, dayOfWeek, startTime, endTime)
	if err != nil {
		http.Error(w, "Failed to save schedule: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

// Handle Leave Submission & Trigger Proxy Engine
func (h *AdminHandler) HandleLeaveSubmission(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	teacherID, _ := strconv.Atoi(r.FormValue("teacher_id"))
	startDate := r.FormValue("start_date")
	endDate := r.FormValue("end_date")
	reason := r.FormValue("reason")

	// Calls ProxyEngine method
	err := h.proxyEngine.ProcessTeacherLeave(teacherID, startDate, endDate, reason)
	if err != nil {
		http.Error(w, "Failed to process leave and assign proxy: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (h *AdminHandler) PostAnnouncement(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	claims, ok := r.Context().Value(middleware.UserContextKey).(jwt.MapClaims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	adminID := int(claims["user_id"].(float64))

	content := r.FormValue("content")
	if content == "" {
		http.Error(w, "Announcement content cannot be empty", http.StatusBadRequest)
		return
	}

	// Save announcement straight into chat_messages under the official channel
	query := `INSERT INTO chat_messages (channel, author_id, content) VALUES ('official', $1, $2)`
	_, err := h.db.Exec(query, adminID, content)
	if err != nil {
		http.Error(w, "Failed to publish announcement", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
