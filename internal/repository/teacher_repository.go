package repository

import (
	"database/sql"
	"school-management-system/internal/models"
)

type TeacherRepository struct {
	db *sql.DB
}

func NewTeacherRepository(db *sql.DB) *TeacherRepository {
	return &TeacherRepository{db: db}
}

// --- Assignment Management ---

func (r *TeacherRepository) CreateAssignment(assignment *models.Assignment) error {
	query := `INSERT INTO assignments (course_id, teacher_id, title, description, due_date) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(query, assignment.CourseID, assignment.TeacherID, assignment.Title, assignment.Description, assignment.DueDate)
	return err
}

// --- Attendance Tracking ---

func (r *TeacherRepository) RecordAttendance(courseID, studentID int, date string, status string) error {
	query := `
		INSERT INTO attendance (course_id, student_id, date, status) 
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (course_id, student_id, date) 
		DO UPDATE SET status = EXCLUDED.status`
	_, err := r.db.Exec(query, courseID, studentID, date, status)
	return err
}