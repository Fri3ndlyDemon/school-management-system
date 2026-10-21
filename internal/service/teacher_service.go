package service

import (
	"errors"
	"school-management-system/internal/models"
	"school-management-system/internal/repository"
)

type TeacherService struct {
	teacherRepo *repository.TeacherRepository
}

func NewTeacherService(teacherRepo *repository.TeacherRepository) *TeacherService {
	return &TeacherService{teacherRepo: teacherRepo}
}

// Post a new assignment
func (s *TeacherService) CreateAssignment(courseID, teacherID int, title, description, dueDate string) error {
	if title == "" || courseID == 0 {
		return errors.New("title and course ID are required for assignments")
	}

	assignment := &models.Assignment{
		CourseID:    courseID,
		TeacherID:   teacherID,
		Title:       title,
		Description: description,
		DueDate:     dueDate,
	}

	return s.teacherRepo.CreateAssignment(assignment)
}

// Log attendance for a student
func (s *TeacherService) MarkAttendance(courseID, studentID int, date, status string) error {
	if status != "Present" && status != "Absent" {
		return errors.New("invalid attendance status")
	}
	return s.teacherRepo.RecordAttendance(courseID, studentID, date, status)
}