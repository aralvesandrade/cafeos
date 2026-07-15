package middleware

import (
	"net/http"

	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	"github.com/aralvesandrade/cafeos/internal/domain/service"
)

func RequireModule(svc *service.PermissionService, module entity.Module, need entity.AccessLevel) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, ok := r.Context().Value(RoleKey).(string)
			if !ok || role == "" {
				http.Error(w, "unauthorized", http.StatusForbidden)
				return
			}
			organizationID, ok := r.Context().Value(OrganizationIDKey).(string)
			if !ok || organizationID == "" {
				http.Error(w, "unauthorized", http.StatusForbidden)
				return
			}

			access, err := svc.GetAccess(organizationID, entity.UserRole(role), module)
			if err != nil {
				http.Error(w, "failed to resolve permissions", http.StatusInternalServerError)
				return
			}
			if !access.Satisfies(need) {
				http.Error(w, "forbidden: insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
