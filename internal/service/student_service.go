package service

import (
	"school-management-system/internal/models"
	"school-management-system/internal/repository"
)

type StudentService struct {
	repo *repository.StudentRepository
}

func NewStudentService(repo *repository.StudentRepository) *StudentService {
	return &StudentService{repo: repo}
}

// gather all relevant data for a student's dashboard
func (s *StudentService) GetDashboardBundle(studentID int) (map[string]interface{}, error) {
	courses, _ := s.repo.GetStudentCourses(studentID)
	schedules, _ := s.repo.GetSchedules(studentID)
	attendance, _ := s.repo.GetAttendanceRecords(studentID)
	results, _ := s.repo.GetExamResults(studentID)
	notifications, _ := s.repo.GetNotifications(studentID)

	presentCount := 0
	for _, att := range attendance {
		if att["Status"] == "Present" {
			presentCount++
		}
	}
	attendancePct := 0.0
	if len(attendance) > 0 {
		attendancePct = (float64(presentCount) / float64(len(attendance))) * 100.0
	}

	return map[string]interface{}{
		"Courses":              courses,
		"Schedules":            schedules,
		"AttendanceRecords":    attendance,
		"AttendancePercentage": attendancePct,
		"ExamResults":          results,
		"Notifications":        notifications,
	}, nil
}

// fetch chat messages for a specific channel in chronological order
func (s *StudentService) GetChannelMessages(channel string) ([]models.ChatMessage, error) {
	return s.repo.GetChannelMessages(channel)
}

// insert a new text message into a designated chat channel
func (s *StudentService) SaveMessage(channel string, authorID int, content string) error {
	return s.repo.SaveMessage(channel, authorID, content)
}
