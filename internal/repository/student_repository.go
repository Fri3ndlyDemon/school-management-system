package repository

import (
	"database/sql"
	"school-management-system/internal/models"
)

type StudentRepository struct {
	db *sql.DB
}

func NewStudentRepository(db *sql.DB) *StudentRepository {
	return &StudentRepository{db: db}
}

// query and return all courses that a specific student is enrolled in
func (r *StudentRepository) GetStudentCourses(studentID int) ([]models.Course, error) {
	query := `
		SELECT c.id, c.title, c.code, c.description, c.teacher_id, c.created_at
		FROM courses c
		JOIN enrollments e ON c.id = e.course_id
		WHERE e.student_id = $1
	`
	rows, err := r.db.Query(query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []models.Course
	for rows.Next() {
		var c models.Course
		var desc sql.NullString
		var teacherID sql.NullInt64
		if err := rows.Scan(&c.ID, &c.Title, &c.Code, &desc, &teacherID, &c.CreatedAt); err != nil {
			return nil, err
		}
		if desc.Valid {
			c.Description = &desc.String
		}
		if teacherID.Valid {
			tID := int(teacherID.Int64)
			c.TeacherID = &tID
		}
		courses = append(courses, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return courses, nil
}

// fetch class routine schedules for all courses associated with the given student id
func (r *StudentRepository) GetSchedules(studentID int) ([]models.Schedule, error) {
	query := `
		SELECT s.id, s.course_id, s.teacher_id, s.proxy_pool, s.day_of_week, s.start_time, s.end_time, s.room
		FROM schedules s
		JOIN enrollments e ON s.course_id = e.course_id
		WHERE e.student_id = $1
	`
	rows, err := r.db.Query(query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []models.Schedule
	for rows.Next() {
		var s models.Schedule
		var teacherID sql.NullInt64
		if err := rows.Scan(&s.ID, &s.CourseID, &teacherID, &s.ProxyPool, &s.DayOfWeek, &s.StartTime, &s.EndTime, &s.Room); err != nil {
			return nil, err
		}
		if teacherID.Valid {
			tID := int(teacherID.Int64)
			s.TeacherID = &tID
		}
		schedules = append(schedules, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return schedules, nil
}

// fetch the most recent attendance logs for a student
func (r *StudentRepository) GetAttendanceRecords(studentID int) ([]map[string]interface{}, error) {
	query := `SELECT date, status FROM attendance WHERE student_id = $1 ORDER BY date DESC LIMIT 20`
	rows, err := r.db.Query(query, studentID)
	if err != nil {
		return []map[string]interface{}{}, nil
	}
	defer rows.Close()

	var records []map[string]interface{}
	for rows.Next() {
		var date, status string
		if err := rows.Scan(&date, &status); err == nil {
			records = append(records, map[string]interface{}{"Date": date, "Status": status})
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return records, nil
}

// fetch all exam grades and scores associated with a student
func (r *StudentRepository) GetExamResults(studentID int) ([]models.ExamResult, error) {
	query := `
		SELECT er.id, er.student_id, c.title, er.exam_title, er.score, er.max_score, er.grade, er.created_at
		FROM exam_results er
		JOIN courses c ON er.course_id = c.id
		WHERE er.student_id = $1
	`
	rows, err := r.db.Query(query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.ExamResult
	for rows.Next() {
		var er models.ExamResult
		if err := rows.Scan(&er.ID, &er.StudentID, &er.CourseName, &er.ExamTitle, &er.Score, &er.MaxScore, &er.Grade, &er.CreatedAt); err != nil {
			return nil, err
		}
		results = append(results, er)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// fetch up to 10 recent notifications intended for a specific user
func (r *StudentRepository) GetNotifications(studentID int) ([]models.Notification, error) {
	query := `SELECT id, user_id, message, is_read, created_at FROM notifications WHERE user_id = $1 ORDER BY created_at DESC LIMIT 10`
	rows, err := r.db.Query(query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifs []models.Notification
	for rows.Next() {
		var n models.Notification
		if err := rows.Scan(&n.ID, &n.UserID, &n.Message, &n.IsRead, &n.CreatedAt); err != nil {
			return nil, err
		}
		notifs = append(notifs, n)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return notifs, nil
}

// fetch chat messages for a specific channel in chronological order
func (r *StudentRepository) GetChannelMessages(channel string) ([]models.ChatMessage, error) {
	query := `
		SELECT cm.id, cm.channel, cm.author_id, u.name, u.role, cm.content, cm.attachments_url, cm.sent_at
		FROM chat_messages cm
		JOIN users u ON cm.author_id = u.id
		WHERE cm.channel = $1
		ORDER BY cm.sent_at ASC
	`
	rows, err := r.db.Query(query, channel)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []models.ChatMessage
	for rows.Next() {
		var m models.ChatMessage
		var attachmentsURL sql.NullString
		if err := rows.Scan(&m.ID, &m.Channel, &m.AuthorID, &m.AuthorName, &m.AuthorRole, &m.Content, &attachmentsURL, &m.SentAt); err != nil {
			return nil, err
		}
		if attachmentsURL.Valid {
			m.AttachmentsURL = &attachmentsURL.String
		}
		msgs = append(msgs, m)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return msgs, nil
}

// insert a new text message into a designated chat channel
func (r *StudentRepository) SaveMessage(channel string, authorID int, content string) error {
	query := `INSERT INTO chat_messages (channel, author_id, content) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(query, channel, authorID, content)
	return err
}
