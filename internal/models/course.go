package models

import "time"

type Course struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Code        string    `json:"code"`
	Description *string   `json:"description,omitempty"`
	TeacherID   *int      `json:"teacher_id,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}
