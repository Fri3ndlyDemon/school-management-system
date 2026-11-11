package handler

import (
	"encoding/json"
	"net/http"
	"school-management-system/internal/service"
	"time"
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

	err = h.AuthService.CreateUser(req.Name, req.Email, hashedPassword, req.Role)
	if err != nil {
		http.Error(w, "Failed to create user (email might already exist)", http.StatusBadRequest)
		return
	}

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

	var email, password string

	if r.Header.Get("Content-Type") == "application/json" {
		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}
		email = req.Email
		password = req.Password
	} else {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}
		email = r.FormValue("email")
		password = r.FormValue("password")
	}

	user, err := h.AuthService.GetUserByEmail(email)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	match := h.AuthService.CheckPasswordHash(password, user.PasswordHash)
	if !match {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := h.AuthService.GenerateJWT(user.ID, user.Name, user.Role)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    token,
			Path:     "/",
			HttpOnly: true,
			Expires:  time.Now().Add(24 * time.Hour),
		})

		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{Token: token})
}
