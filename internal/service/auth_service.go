package service

import (
	"database/sql"
	"errors"
	"school-management-system/internal/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	db        *sql.DB
	jwtSecret []byte
}

func NewAuthService(db *sql.DB, secret string) *AuthService {
	return &AuthService{
		db:        db,
		jwtSecret: []byte(secret),
	}
}

func (s *AuthService) CreateUser(name, email, passwordHash, role string) error {
	query := `INSERT INTO users (name, email, password_hash, role) VALUES ($1, $2, $3, $4)`
	_, err := s.db.Exec(query, name, email, passwordHash, role)
	return err
}

func (s *AuthService) GetUserByEmail(email string) (*models.User, error) {
	query := `SELECT id, name, email, password_hash, role, created_at FROM users WHERE email = $1`
	user := &models.User{}

	err := s.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

// convert password into bcrypt hash
func (s *AuthService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// compare password with bcrypt hash
func (s *AuthService) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// generate a single token with userid and role
func (s *AuthService) GenerateJWT(userID int, name string, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"name":    name,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // 24 hour expiration
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

// verify incoming token and validate
func (s *AuthService) ValidateJWT(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid or expired token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}
