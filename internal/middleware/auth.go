package middleware

import (
	"context"
	"net/http"
	"strings"

	"school-management-system/internal/service"

	"github.com/golang-jwt/jwt/v5"
)

// custom type to avoid context key collisions
type contextKey string

const UserContextKey contextKey = "userClaims"

// obtain jwt claims, validate and inject into request context
func AuthMiddleware(authService *service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenStr := ""

			// check authorization header first
			if h := r.Header.Get("Authorization"); strings.HasPrefix(h, "Bearer ") {
				tokenStr = h[7:]
			} else if c, err := r.Cookie("token"); err == nil {
				// fallback to cookie if header is missing
				tokenStr = c.Value
			}

			// if no token is found anywhere, redirect to login
			if tokenStr == "" {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			// validate token and extract claims using your AuthService
			claims, err := authService.ValidateJWT(tokenStr)
			if err != nil {
				http.Error(w, "Unauthorized: Invalid or expired token", http.StatusUnauthorized)
				return
			}

			// inject claims into request context for next handlers
			ctx := context.WithValue(r.Context(), UserContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetClaims(ctx context.Context) (jwt.MapClaims, bool) {
	claims, ok := ctx.Value(UserContextKey).(jwt.MapClaims)
	return claims, ok
}
