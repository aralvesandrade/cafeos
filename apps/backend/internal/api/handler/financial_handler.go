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

// Create registra uma nova transação financeira para o tenant autenticado
// @Summary Criar transação financeira
// @Description Registra uma nova transação financeira (receita/despesa) no tenant
// @Tags financial (Financeiro)
// @Accept json
// @Produce json
// @Param tenant_id path string true "ID do Tenant"
// @Param financial body createFinancialRequest true "Dados da transação financeira"
// @Success 201 {object} entity.FinancialTransaction
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{tenant_id}/financial [post]
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

// List retorna todas as transações financeiras do tenant autenticado
// @Summary Listar transações financeiras
// @Description Lista todas as transações financeiras pertencentes ao tenant
// @Tags financial (Financeiro)
// @Produce json
// @Param tenant_id path string true "ID do Tenant"
// @Success 200 {array} entity.FinancialTransaction
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{tenant_id}/financial [get]
func (h *FinancialHandler) List(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := r.Context().Value(middleware.TenantIDKey).(string)
	txs, err := h.svc.ListByTenant(tenantID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, txs, http.StatusOK)
}

// GetByID retorna uma transação financeira pelo seu ID
// @Summary Obter transação financeira por ID
// @Description Retorna uma única transação financeira
// @Tags financial (Financeiro)
// @Produce json
// @Param tenant_id path string true "ID do Tenant"
// @Param id path string true "ID da Transação Financeira"
// @Success 200 {object} entity.FinancialTransaction
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{tenant_id}/financial/{id} [get]
func (h *FinancialHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	tx, err := h.svc.GetByID(id)
	if err != nil {
		writeError(w, "transaction not found", http.StatusNotFound)
		return
	}
	writeJSON(w, tx, http.StatusOK)
}

// Update atualiza uma transação financeira existente
// @Summary Atualizar transação financeira
// @Description Atualiza dados da transação financeira por ID (atualização parcial - somente os campos informados são alterados)
// @Tags financial (Financeiro)
// @Accept json
// @Produce json
// @Param tenant_id path string true "ID do Tenant"
// @Param id path string true "ID da Transação Financeira"
// @Param financial body entity.FinancialTransaction true "Dados atualizados da transação financeira"
// @Success 200 {object} entity.FinancialTransaction
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{tenant_id}/financial/{id} [put]
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

// Delete remove uma transação financeira pelo seu ID
// @Summary Excluir transação financeira
// @Description Exclui uma transação financeira por ID
// @Tags financial (Financeiro)
// @Param tenant_id path string true "ID do Tenant"
// @Param id path string true "ID da Transação Financeira"
// @Success 204 "No Content"
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{tenant_id}/financial/{id} [delete]
func (h *FinancialHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.svc.Delete(id); err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
