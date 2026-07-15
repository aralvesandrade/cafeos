package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/aralvesandrade/cafeos/internal/api/middleware"
	"github.com/aralvesandrade/cafeos/internal/domain/repository"
)

type AlertHandler struct {
	repo repository.AlertRepository
}

func NewAlertHandler(repo repository.AlertRepository) *AlertHandler {
	return &AlertHandler{repo: repo}
}

// List retorna todos os alertas da organização autenticada
// @Summary Listar alertas
// @Description Lista todos os alertas gerados pelo Rule Engine para a organização
// @Tags alerts (Alertas)
// @Produce json
// @Param organization_id path string true "ID da Organização"
// @Param status query string false "Filtrar por status (aberto, resolvido, descartado)"
// @Success 200 {array} entity.Alert
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/alerts [get]
func (h *AlertHandler) List(w http.ResponseWriter, r *http.Request) {
	organizationID, _ := r.Context().Value(middleware.OrganizationIDKey).(string)
	alerts, err := h.repo.ListByOrganization(organizationID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if status := r.URL.Query().Get("status"); status != "" {
		filtered := alerts[:0]
		for _, a := range alerts {
			if a.Status == status {
				filtered = append(filtered, a)
			}
		}
		alerts = filtered
	}

	writeJSON(w, alerts, http.StatusOK)
}

// Update atualiza o status de um alerta (resolver/descartar)
// @Summary Atualizar alerta
// @Description Atualiza o status de um alerta (resolvido/descartado)
// @Tags alerts (Alertas)
// @Accept json
// @Produce json
// @Param organization_id path string true "ID da Organização"
// @Param id path string true "ID do Alerta"
// @Param alert body object true "Novo status do alerta"
// @Success 200 {object} entity.Alert
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/alerts/{id} [put]
func (h *AlertHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	existing, err := h.repo.GetByID(id)
	if err != nil {
		writeError(w, "alert not found", http.StatusNotFound)
		return
	}

	var input struct {
		Status *string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if input.Status != nil {
		existing.Status = *input.Status
		if *input.Status != "aberto" {
			now := time.Now()
			existing.ResolvedAt = &now
		} else {
			existing.ResolvedAt = nil
		}
	}

	if err := h.repo.Update(existing); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, existing, http.StatusOK)
}
