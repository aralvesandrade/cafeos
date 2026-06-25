package middleware

import (
	"context"
	"net/http"
)

func Tenant() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tenantID := r.PathValue("tenant_id")
			if tenantID == "" {
				tenantID = r.Header.Get("X-Tenant-ID")
			}

			if tenantID != "" {
				ctx := context.WithValue(r.Context(), TenantIDKey, tenantID)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
