package middleware

import (
	"context"
	"net/http"
)

func Organization() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			organizationID := r.PathValue("organization_id")
			if organizationID == "" {
				organizationID = r.Header.Get("X-Organization-ID")
			}

			if organizationID != "" {
				ctx := context.WithValue(r.Context(), OrganizationIDKey, organizationID)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
