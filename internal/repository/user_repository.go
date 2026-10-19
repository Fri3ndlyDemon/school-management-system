package repository

import (
	"database/sql"
	"school-management-system/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser handles manual admin onboarding for students and teachers
func (r *UserRepository) CreateUser(user *models.User, hashedPassword string) error {
	query := `
		INSERT INTO users (name, email, password_hash, role, created_at)
		VALUES ($1, $2, $3, $4, NOW())
		RETURNING id, created_at
	`
	err := r.db.QueryRow(
		query,
		user.Name,
		user.Email,
		hashedPassword,
		user.Role,
	).Scan(&user.ID, &user.CreatedAt)

	return err
}