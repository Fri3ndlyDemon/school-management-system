package service

import (
	"net/http"
	"strconv"
	"strings"
)

type AdminHandler struct {
	userService     *UserService
	scheduleService *ScheduleService
	proxyEngine     *ProxyEngine
}

func NewAdminHandler(userService *UserService, scheduleService *ScheduleService, proxyEngine *ProxyEngine) *AdminHandler {
	return &AdminHandler{
		userService:     userService,
		scheduleService: scheduleService,
		proxyEngine:     proxyEngine,
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

	http.Redirect(w, r, "/admin/onboarding?success=true", http.StatusSeeOther)
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

	http.Redirect(w, r, "/admin/schedule?success=true", http.StatusSeeOther)
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

	http.Redirect(w, r, "/admin/leaves?success=true", http.StatusSeeOther)
}