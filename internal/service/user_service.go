package service

import (
	"errors"
	"school-management-system/internal/models"
	"school-management-system/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

// OnboardUser:
// hashes the password and registers a new student or teacher
func (s *UserService) OnboardUser(name, email, plainPassword, role string) error {
	// Validate role permissions
	if role != "student" && role != "teacher" && role != "admin" {
		return errors.New("invalid user role specified")
	}

	// Hash the password
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &models.User{
		Name:  name,
		Email: email,
		Role:  role,
	}

	return s.userRepo.CreateUser(user, string(hashedBytes))
}