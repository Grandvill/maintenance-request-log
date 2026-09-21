package middleware

import (
	"net/http"
)

// RequireRoles restricts route access to specified roles
func RequireRoles(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := GetUserFromContext(r.Context())
			if !ok {
				respondJSON(w, http.StatusUnauthorized, map[string]string{
					"error": "Unauthorized: user identity not found in context",
				})
				return
			}

			allowed := false
			for _, role := range allowedRoles {
				if user.Role == role {
					allowed = true
					break
				}
			}

			if !allowed {
				respondJSON(w, http.StatusForbidden, map[string]string{
					"error": "Forbidden: your role does not have permission to perform this action",
				})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

