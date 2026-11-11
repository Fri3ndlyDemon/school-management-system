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

// fetch chat messages for a specific channel in chronological order
func (s *StudentService) GetChannelMessages(channel string) ([]models.ChatMessage, error) {
	return s.repo.GetChannelMessages(channel)
}

// insert a new text message into a designated chat channel
func (s *StudentService) SaveMessage(channel string, authorID int, content string) error {
	return s.repo.SaveMessage(channel, authorID, content)
}
