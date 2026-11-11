package handler

import (
	"net/http"
	"school-management-system/internal/middleware"
	"school-management-system/internal/service"
	"text/template"

	"github.com/golang-jwt/jwt/v5"
)

type StudentHandler struct {
	studentService *service.StudentService
}

func NewStudentHandler(studentService *service.StudentService) *StudentHandler {
	return &StudentHandler{studentService: studentService}
}

// render the chat channels view for students
func (h *StudentHandler) ChannelsView(w http.ResponseWriter, r *http.Request) {
	channel := r.URL.Query().Get("channel")
	if channel == "" {
		channel = "general"
	}

	messages, err := h.studentService.GetChannelMessages(channel)
	if err != nil {
		http.Error(w, "Failed to load channel messages", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"ActiveChannel": channel,
		"Messages":      messages,
	}

	tmpl, err := template.ParseFiles("internal/web/templates/student/channels.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}

// handle posting a new chat message to a channel
func (h *StudentHandler) PostMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	claims := r.Context().Value(middleware.UserContextKey).(jwt.MapClaims)
	authorID := int(claims["user_id"].(float64))
	role := claims["role"].(string)

	channel := r.FormValue("channel")
	content := r.FormValue("content")

	// restrict posting to the official channel
	if channel == "official" && role != "admin" && role != "teacher" {
		http.Error(w, "Forbidden: Only faculty/admin can post official announcements", http.StatusForbidden)
		return
	}

	err := h.studentService.SaveMessage(channel, authorID, content)
	if err != nil {
		http.Error(w, "Failed to send message", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/student/channels?channel="+channel, http.StatusSeeOther)
}
