package handler

import (
	"encoding/json"
	"net/http"
	"school-management-system/internal/service"
)

type AuthHandler struct {
	AuthService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		AuthService: authService,
	}
}

// holds input data for new registration
type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// holds input data for login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// holds output data for login response
type LoginResponse struct {
	Token string `json:"token"`
}

// user registration handler
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// password hashing
	hashedPassword, err := h.AuthService.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Failed to process password", http.StatusInternalServerError)
		return
	}

	// save user to database - incomplete
	_ = hashedPassword

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

// login handler
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// fetch user from database - incomplete
	// dummy values
	dummyUserID := 1
	dummyRole := "student"
	dummyHash := "$2a$10$..." // stored hash

	// compare password with hash
	if !h.AuthService.CheckPasswordHash(req.Password, dummyHash) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// generate jwt
	token, err := h.AuthService.GenerateJWT(dummyUserID, dummyRole)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{Token: token})
}
