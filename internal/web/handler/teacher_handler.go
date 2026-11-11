package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"school-management-system/internal/service"
	"strconv"
)

type TeacherHandler struct {
	teacherService *service.TeacherService
	db             *sql.DB // Added database connection for direct message posting
}

func NewTeacherHandler(teacherService *service.TeacherService, db *sql.DB) *TeacherHandler {
	return &TeacherHandler{
		teacherService: teacherService,
		db:             db,
	}
}

// Handle posting a new assignment
func (h *TeacherHandler) HandleCreateAssignment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	courseID, _ := strconv.Atoi(r.FormValue("course_id"))
	teacherID, _ := strconv.Atoi(r.FormValue("teacher_id"))
	title := r.FormValue("title")
	description := r.FormValue("description")
	dueDate := r.FormValue("due_date")

	// Create the assignment in the database
	err := h.teacherService.CreateAssignment(courseID, teacherID, title, description, dueDate)
	if err != nil {
		http.Error(w, "Failed to create assignment: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Format the message into a nice clean update format
	formattedMessage := fmt.Sprintf(
		"<strong>NEW ASSIGNMENT POSTED</strong><br><br>"+
			"<strong>Title:</strong> %s<br>"+
			"<strong>Course ID:</strong> %d<br>"+
			"<strong>Due Date:</strong> %s<br><br>"+
			"<em>%s</em>",
		title, courseID, dueDate, description,
	)

	// Automatically insert the message into the updates channel table
	if h.db != nil {
		query := `INSERT INTO chat_messages (channel, author_id, content, sent_at) VALUES ($1, $2, $3, CURRENT_TIMESTAMP)`
		// We use teacherID as the author_id since it's an integer referencing users(id)
		_, err = h.db.Exec(query, "updates", teacherID, formattedMessage)
		if err != nil {
			fmt.Printf("Failed to broadcast update message: %v\n", err)
		}
	}

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

// Handle recording student attendance
func (h *TeacherHandler) HandleMarkAttendance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
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

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
