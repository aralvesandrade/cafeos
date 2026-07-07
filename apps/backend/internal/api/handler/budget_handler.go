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

func (h *BudgetHandler) Create(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := r.Context().Value(middleware.TenantIDKey).(string)

	var req createBudgetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	b, err := h.svc.Create(tenantID, req.HarvestID, req.CostCenterID, req.PlannedAmount, req.Description)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, b, http.StatusCreated)
}

func (h *BudgetHandler) ListByHarvest(w http.ResponseWriter, r *http.Request) {
	harvestID := r.PathValue("harvest_id")
	budgets, err := h.svc.ListByHarvest(harvestID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, budgets, http.StatusOK)
}

func (h *BudgetHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	b, err := h.svc.GetByID(id)
	if err != nil {
		writeError(w, "budget not found", http.StatusNotFound)
		return
	}
	writeJSON(w, b, http.StatusOK)
}

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

func (h *BudgetHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.svc.Delete(id); err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
