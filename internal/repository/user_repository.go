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

// CreateUser :
//  handles manual admin onboarding for students and teachers
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


//GetByEmail : 
// this looks up the email for entry
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	u := &models.User{}
	err := r.db.QueryRow(
		`SELECT id, name, email, password_hash, role, created_at FROM users WHERE email = $1`, email,
	).Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt)
	return u, err
}