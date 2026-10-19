package repository

import (
	"database/sql"
	"school-management-system/internal/models"
)

type ScheduleRepository struct {
	db *sql.DB
}

func NewScheduleRepository(db *sql.DB) *ScheduleRepository {
	return &ScheduleRepository{db: db}
}

// CreateSchedule adds a master routine entry with primary and proxy teachers
func (r *ScheduleRepository) CreateSchedule(s *models.Schedule) error {
	query := `
		INSERT INTO schedules (course_id, teacher_id, proxy_pool, day_of_week, start_time, end_time)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	// Note: proxy_pool can be stored as an array of integer IDs or text in PostgreSQL
	err := r.db.QueryRow(
		query,
		s.CourseID,
		s.TeacherID,
		s.ProxyPool, // Assuming []int or pq.Int64Array / pq.StringArray in your model
		s.DayOfWeek,
		s.StartTime,
		s.EndTime,
	).Scan(&s.ID)

	return err
}