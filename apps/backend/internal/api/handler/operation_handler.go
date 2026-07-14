package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/aralvesandrade/cafeos/internal/api/middleware"
	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	"github.com/aralvesandrade/cafeos/internal/domain/service"
)

type OperationHandler struct {
	svc *service.OperationService
}

func NewOperationHandler(svc *service.OperationService) *OperationHandler {
	return &OperationHandler{svc: svc}
}

type createOperationRequest struct {
	PlotID       string  `json:"plot_id"`
	HarvestID    *string `json:"harvest_id"`
	CostCenterID *string `json:"cost_center_id"`
	Type         string  `json:"type"`
	Date         string  `json:"date"`
	Responsible  string  `json:"responsible"`
	ProductUsed  string  `json:"product_used"`
	Quantity     float64 `json:"quantity"`
	Cost         float64 `json:"cost"`
	Notes        string  `json:"notes"`
}

// Create registra uma nova operação agrícola
// @Summary Registrar operação
// @Description Registra uma operação agrícola (adubação, pulverização, irrigação, poda, colheita)
// @Tags operations (Operações)
// @Accept json
// @Produce json
// @Param organization_id path string true "ID da Organização"
// @Param operation body createOperationRequest true "Dados da operação"
// @Success 201 {object} SwaggerOperation
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/operations [post]
func (h *OperationHandler) Create(w http.ResponseWriter, r *http.Request) {
	organizationID, _ := r.Context().Value(middleware.OrganizationIDKey).(string)

	var req createOperationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	date, err := time.Parse(time.RFC3339, req.Date)
	if err != nil {
		date = time.Now()
	}

	op, err := h.svc.Create(organizationID, req.PlotID, entity.OperationType(req.Type), date, req.Responsible, req.ProductUsed, req.Quantity, req.Cost, req.Notes, req.HarvestID, req.CostCenterID)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, op, http.StatusCreated)
}

// GetByID retorna uma operação pelo seu ID
// @Summary Obter operação por ID
// @Description Retorna uma única operação
// @Tags operations (Operações)
// @Produce json
// @Param organization_id path string true "ID da Organização"
// @Param id path string true "ID da Operação"
// @Success 200 {object} SwaggerOperation
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/operations/{id} [get]
func (h *OperationHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	op, err := h.svc.GetByID(id)
	if err != nil {
		writeError(w, "operation not found", http.StatusNotFound)
		return
	}
	writeJSON(w, op, http.StatusOK)
}

// ListByPlot retorna todas as operações de um talhão
// @Summary Listar operações por talhão
// @Description Lista todas as operações de um talhão específico
// @Tags operations (Operações)
// @Produce json
// @Param organization_id path string true "ID da Organização"
// @Param plot_id path string true "ID do Talhão"
// @Success 200 {array} SwaggerOperation
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/plots/{plot_id}/operations [get]
func (h *OperationHandler) ListByPlot(w http.ResponseWriter, r *http.Request) {
	plotID := r.PathValue("plot_id")
	ops, err := h.svc.ListByPlot(plotID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, ops, http.StatusOK)
}

// List retorna todas as operações da organização autenticada
// @Summary Listar todas as operações
// @Description Lista todas as operações de todos os talhões da organização
// @Tags operations (Operações)
// @Produce json
// @Param organization_id path string true "ID da Organização"
// @Success 200 {array} SwaggerOperation
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/operations [get]
func (h *OperationHandler) List(w http.ResponseWriter, r *http.Request) {
	organizationID, _ := r.Context().Value(middleware.OrganizationIDKey).(string)
	ops, err := h.svc.ListByOrganization(organizationID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, ops, http.StatusOK)
}

// ListRecent retorna as operações mais recentes
// @Summary Listar operações recentes
// @Description Lista as operações mais recentes, limitadas por parâmetro de consulta
// @Tags operations (Operações)
// @Produce json
// @Param organization_id path string true "ID da Organização"
// @Param limit query int false "Máximo de resultados (padrão 10)"
// @Success 200 {array} SwaggerOperation
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/operations/recent [get]
func (h *OperationHandler) ListRecent(w http.ResponseWriter, r *http.Request) {
	organizationID, _ := r.Context().Value(middleware.OrganizationIDKey).(string)
	limitStr := r.URL.Query().Get("limit")
	limit := 10
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}

	ops, err := h.svc.ListRecent(organizationID, limit)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, ops, http.StatusOK)
}

// Delete remove uma operação
// @Summary Excluir operação
// @Description Exclui uma operação por ID
// @Tags operations (Operações)
// @Param organization_id path string true "ID da Organização"
// @Param id path string true "ID da Operação"
// @Success 204 "No Content"
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/operations/{id} [delete]
func (h *OperationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.svc.Delete(id); err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
