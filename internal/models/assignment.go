package models

type Assignment struct {
	ID          int    `json:"id"`
	CourseID    int    `json:"course_id"`
	TeacherID   int    `json:"teacher_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
}