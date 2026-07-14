package handler

import (
	"encoding/json"
	"net/http"

	"github.com/aralvesandrade/cafeos/internal/api/middleware"
	"github.com/aralvesandrade/cafeos/internal/domain/service"
)

type BudgetHandler struct {
	svc *service.BudgetService
}

func NewBudgetHandler(svc *service.BudgetService) *BudgetHandler {
	return &BudgetHandler{svc: svc}
}

type createBudgetRequest struct {
	HarvestID     string  `json:"harvest_id"`
	CostCenterID  string  `json:"cost_center_id"`
	PlannedAmount float64 `json:"planned_amount"`
	Description   string  `json:"description"`
}

// Create registra um novo orçamento para um centro de custo da colheita
// @Summary Criar orçamento
// @Description Registra um novo orçamento para uma colheita e centro de custo
// @Tags budgets (Orçamentos)
// @Accept json
// @Produce json
// @Param organization_id path string true "ID da Organização"
// @Param budget body createBudgetRequest true "Dados do orçamento"
// @Success 201 {object} entity.Budget
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/budgets [post]
func (h *BudgetHandler) Create(w http.ResponseWriter, r *http.Request) {
	organizationID, _ := r.Context().Value(middleware.OrganizationIDKey).(string)

	var req createBudgetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	b, err := h.svc.Create(organizationID, req.HarvestID, req.CostCenterID, req.PlannedAmount, req.Description)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, b, http.StatusCreated)
}

// ListByHarvest retorna todos os orçamentos de uma colheita
// @Summary Listar orçamentos por colheita
// @Description Lista todos os orçamentos pertencentes a uma colheita
// @Tags budgets (Orçamentos)
// @Produce json
// @Param organization_id path string true "ID da Organização"
// @Param harvest_id path string true "ID da Colheita"
// @Success 200 {array} entity.Budget
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/harvests/{harvest_id}/budgets [get]
func (h *BudgetHandler) ListByHarvest(w http.ResponseWriter, r *http.Request) {
	harvestID := r.PathValue("harvest_id")
	budgets, err := h.svc.ListByHarvest(harvestID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, budgets, http.StatusOK)
}

// GetByID retorna um orçamento pelo seu ID
// @Summary Obter orçamento por ID
// @Description Retorna um único orçamento
// @Tags budgets (Orçamentos)
// @Produce json
// @Param organization_id path string true "ID da Organização"
// @Param id path string true "ID do Orçamento"
// @Success 200 {object} entity.Budget
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/budgets/{id} [get]
func (h *BudgetHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	b, err := h.svc.GetByID(id)
	if err != nil {
		writeError(w, "budget not found", http.StatusNotFound)
		return
	}
	writeJSON(w, b, http.StatusOK)
}

// Update atualiza um orçamento existente
// @Summary Atualizar orçamento
// @Description Atualiza dados do orçamento por ID (atualização parcial - somente os campos informados são alterados)
// @Tags budgets (Orçamentos)
// @Accept json
// @Produce json
// @Param organization_id path string true "ID da Organização"
// @Param id path string true "ID do Orçamento"
// @Param budget body entity.Budget true "Dados atualizados do orçamento"
// @Success 200 {object} entity.Budget
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/budgets/{id} [put]
func (h *BudgetHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	existing, err := h.svc.GetByID(id)
	if err != nil {
		writeError(w, "budget not found", http.StatusNotFound)
		return
	}

	var input struct {
		PlannedAmount *float64 `json:"planned_amount"`
		Description   *string  `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if input.PlannedAmount != nil {
		existing.PlannedAmount = *input.PlannedAmount
	}
	if input.Description != nil {
		existing.Description = *input.Description
	}
	if err := h.svc.Update(existing); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, existing, http.StatusOK)
}

// Delete remove um orçamento pelo seu ID
// @Summary Excluir orçamento
// @Description Exclui um orçamento por ID
// @Tags budgets (Orçamentos)
// @Param organization_id path string true "ID da Organização"
// @Param id path string true "ID do Orçamento"
// @Success 204 "No Content"
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/budgets/{id} [delete]
func (h *BudgetHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.svc.Delete(id); err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
