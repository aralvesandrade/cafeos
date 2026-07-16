package handler

import (
	"net/http"

	"github.com/aralvesandrade/cafeos/internal/api/middleware"
	"github.com/aralvesandrade/cafeos/internal/domain/entity"
)

// restrictedOwnerID reports whether the authenticated request should be scoped
// to farms the user is linked to via the Producer table. System roles
// (platform_owner, organization_admin) bypass scoping. Every other role sees
// only farms they are a producer on.
func restrictedOwnerID(r *http.Request) (userID string, restricted bool) {
	role, _ := r.Context().Value(middleware.RoleKey).(string)
	if role == entity.SystemRolePlatformOwner || role == entity.SystemRoleOrganizationAdmin {
		return "", false
	}
	userID, _ = r.Context().Value(middleware.UserIDKey).(string)
	return userID, true
}
