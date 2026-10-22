package handler

import (
	"net/http"
	"strconv"
	"school-management-system/internal/service"
)

type TeacherHandler struct {
	teacherService *service.TeacherService
}

func NewTeacherHandler(teacherService *service.TeacherService) *TeacherHandler {
	return &TeacherHandler{teacherService: teacherService}
}

// Handle posting a new assignment
func (h *TeacherHandler) HandleCreateAssignment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	courseID, _ := strconv.Atoi(r.FormValue("course_id"))
	teacherID, _ := strconv.Atoi(r.FormValue("teacher_id")) // In real flow, you can also extract this from the JWT session claims
	title := r.FormValue("title")
	description := r.FormValue("description")
	dueDate := r.FormValue("due_date")

	err := h.teacherService.CreateAssignment(courseID, teacherID, title, description, dueDate)
	if err != nil {
		http.Error(w, "Failed to create assignment: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/teacher/assignments?success=true", http.StatusSeeOther)
}

// Handle recording student attendance
func (h *TeacherHandler) HandleMarkAttendance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	courseID, _ := strconv.Atoi(r.FormValue("course_id"))
	studentID, _ := strconv.Atoi(r.FormValue("student_id"))
	date := r.FormValue("date")
	status := r.FormValue("status") // "Present" or "Absent"

	err := h.teacherService.MarkAttendance(courseID, studentID, date, status)
	if err != nil {
		http.Error(w, "Failed to mark attendance: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/teacher/attendance?success=true", http.StatusSeeOther)
}