package models

import "time"

type ExamResult struct {
	ID         int       `json:"id"`
	StudentID  int       `json:"student_id"`
	CourseName string    `json:"course_name"`
	ExamTitle  string    `json:"exam_title"`
	Score      float64   `json:"score"`
	MaxScore   float64   `json:"max_score"`
	Grade      string    `json:"grade"`
	CreatedAt  time.Time `json:"created_at"`
}
