package handler

import "net/http"

// HealthCheck returns the API health status
// @Summary Health check
// @Description Returns API health status
// @Tags system
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, map[string]string{"status": "ok", "service": "cafeos-api"}, http.StatusOK)
}
