package middleware

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

// check if user has required roles to access the endpoint
func RequireRole(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// extract claims
			claims, ok := r.Context().Value(UserContextKey).(jwt.MapClaims)
			if !ok {
				http.Error(w, "Unauthorized: No user session found", http.StatusUnauthorized)
				return
			}

			// extract the role
			role, ok := claims["role"].(string)
			if !ok {
				http.Error(w, "Forbidden: Invalid role claims", http.StatusForbidden)
				return
			}

			// check if the role is in the allowed roles
			isAllowed := false
			for _, allowedRole := range allowedRoles {
				if role == allowedRole {
					isAllowed = true
					break
				}
			}

			if !isAllowed {
				http.Error(w, "Forbidden: Insufficient permissions for this action", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
