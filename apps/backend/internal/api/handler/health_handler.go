package handler

import "net/http"

// HealthCheck retorna o status de saúde da API
// @Summary Verificação de saúde
// @Description Retorna o status de saúde da API
// @Tags system (Sistema)
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, map[string]string{"status": "ok", "service": "cafeos-api"}, http.StatusOK)
}
