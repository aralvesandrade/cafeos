package handler

import (
	"net/http"

	"github.com/aralvesandrade/cafeos/internal/api/middleware"
	"github.com/aralvesandrade/cafeos/internal/domain/entity"
)

// restrictedOwnerID reports whether the authenticated request belongs to a
// proprietario (farm owner) — the only role scoped to its own farms — and
// returns the user ID to scope by. Every other role sees the whole organization.
func restrictedOwnerID(r *http.Request) (userID string, restricted bool) {
	role, _ := r.Context().Value(middleware.RoleKey).(string)
	if role != entity.RoleKeyProprietario {
		return "", false
	}
	userID, _ = r.Context().Value(middleware.UserIDKey).(string)
	return userID, true
}
