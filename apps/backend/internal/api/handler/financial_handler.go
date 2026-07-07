package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/aralvesandrade/cafeos/internal/api/middleware"
	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	"github.com/aralvesandrade/cafeos/internal/domain/service"
)

type FinancialHandler struct {
	svc *service.FinancialService
}

func NewFinancialHandler(svc *service.FinancialService) *FinancialHandler {
	return &FinancialHandler{svc: svc}
}

type createFinancialRequest struct {
	Type         string  `json:"type"`
	CostCenterID *string `json:"cost_center_id"`
	Description  string  `json:"description"`
	Amount       float64 `json:"amount"`
	Date         string  `json:"date"`
	DueDate      string  `json:"due_date"`
	Notes        string  `json:"notes"`
}

func (h *FinancialHandler) Create(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := r.Context().Value(middleware.TenantIDKey).(string)
	var req createFinancialRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	date, _ := time.Parse("2006-01-02", req.Date)
	dueDate, _ := time.Parse("2006-01-02", req.DueDate)
	tx, err := h.svc.Create(tenantID, req.Type, req.CostCenterID, req.Description, req.Amount, date, dueDate, req.Notes)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, tx, http.StatusCreated)
}

func (h *FinancialHandler) List(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := r.Context().Value(middleware.TenantIDKey).(string)
	txs, err := h.svc.ListByTenant(tenantID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, txs, http.StatusOK)
}

func (h *FinancialHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	tx, err := h.svc.GetByID(id)
	if err != nil {
		writeError(w, "transaction not found", http.StatusNotFound)
		return
	}
	writeJSON(w, tx, http.StatusOK)
}

func (h *FinancialHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	existing, err := h.svc.GetByID(id)
	if err != nil {
		writeError(w, "transaction not found", http.StatusNotFound)
		return
	}
	var req struct {
		Type         *string  `json:"type"`
		CostCenterID *string  `json:"cost_center_id"`
		Description  *string  `json:"description"`
		Amount       *float64 `json:"amount"`
		Status       *string  `json:"status"`
		Notes        *string  `json:"notes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.Type != nil {
		existing.Type = entity.TransactionType(*req.Type)
	}
	if req.CostCenterID != nil {
		existing.CostCenterID = req.CostCenterID
	}
	if req.Description != nil {
		existing.Description = *req.Description
	}
	if req.Amount != nil {
		existing.Amount = *req.Amount
	}
	if req.Status != nil {
		existing.Status = *req.Status
	}
	if req.Notes != nil {
		existing.Notes = *req.Notes
	}
	if err := h.svc.Update(existing); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, existing, http.StatusOK)
}

func (h *FinancialHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.svc.Delete(id); err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
