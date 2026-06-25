package middleware

import (
	"net/http"

	"github.com/aralvesandrade/cafeos/internal/domain/entity"
)

func RequireRole(allowedRoles ...entity.UserRole) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, ok := r.Context().Value(RoleKey).(string)
			if !ok {
				http.Error(w, "unauthorized", http.StatusForbidden)
				return
			}

			for _, allowed := range allowedRoles {
				if string(allowed) == role {
					next.ServeHTTP(w, r)
					return
				}
			}

			http.Error(w, "forbidden: insufficient permissions", http.StatusForbidden)
		})
	}
}
