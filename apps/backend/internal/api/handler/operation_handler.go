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
	PlotID       string   `json:"plot_id"`
	HarvestID    *string  `json:"harvest_id"`
	CostCenterID *string  `json:"cost_center_id"`
	Type         string   `json:"type"`
	Date         string   `json:"date"`
	Responsible  string   `json:"responsible"`
	ProductUsed  string   `json:"product_used"`
	Quantity     float64  `json:"quantity"`
	Cost         float64  `json:"cost"`
	Notes        string   `json:"notes"`
}

// Create registers a new agricultural operation
// @Summary Register an operation
// @Description Register an agricultural operation (adubação, pulverização, irrigação, poda, colheita)
// @Tags operations
// @Accept json
// @Produce json
// @Param tenant_id path string true "Tenant ID"
// @Param operation body createOperationRequest true "Operation data"
// @Success 201 {object} SwaggerOperation
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{tenant_id}/operations [post]
func (h *OperationHandler) Create(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := r.Context().Value(middleware.TenantIDKey).(string)

	var req createOperationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	date, err := time.Parse(time.RFC3339, req.Date)
	if err != nil {
		date = time.Now()
	}

	op, err := h.svc.Create(tenantID, req.PlotID, entity.OperationType(req.Type), date, req.Responsible, req.ProductUsed, req.Quantity, req.Cost, req.Notes, req.HarvestID, req.CostCenterID)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, op, http.StatusCreated)
}

// GetByID returns an operation by its ID
// @Summary Get operation by ID
// @Description Returns a single operation
// @Tags operations
// @Produce json
// @Param tenant_id path string true "Tenant ID"
// @Param id path string true "Operation ID"
// @Success 200 {object} SwaggerOperation
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{tenant_id}/operations/{id} [get]
func (h *OperationHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	op, err := h.svc.GetByID(id)
	if err != nil {
		writeError(w, "operation not found", http.StatusNotFound)
		return
	}
	writeJSON(w, op, http.StatusOK)
}

// ListByPlot returns all operations for a given plot
// @Summary List operations by plot
// @Description List all operations for a specific plot (talhão)
// @Tags operations
// @Produce json
// @Param tenant_id path string true "Tenant ID"
// @Param plot_id path string true "Plot ID"
// @Success 200 {array} SwaggerOperation
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{tenant_id}/plots/{plot_id}/operations [get]
func (h *OperationHandler) ListByPlot(w http.ResponseWriter, r *http.Request) {
	plotID := r.PathValue("plot_id")
	ops, err := h.svc.ListByPlot(plotID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, ops, http.StatusOK)
}

// List returns all operations for the authenticated tenant
// @Summary List all operations
// @Description List all operations across all plots in the tenant
// @Tags operations
// @Produce json
// @Param tenant_id path string true "Tenant ID"
// @Success 200 {array} SwaggerOperation
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{tenant_id}/operations [get]
func (h *OperationHandler) List(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := r.Context().Value(middleware.TenantIDKey).(string)
	ops, err := h.svc.ListByTenant(tenantID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, ops, http.StatusOK)
}

// ListRecent returns the most recent operations
// @Summary List recent operations
// @Description List the most recent operations, limited by query param
// @Tags operations
// @Produce json
// @Param tenant_id path string true "Tenant ID"
// @Param limit query int false "Max results (default 10)"
// @Success 200 {array} SwaggerOperation
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{tenant_id}/operations/recent [get]
func (h *OperationHandler) ListRecent(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := r.Context().Value(middleware.TenantIDKey).(string)
	limitStr := r.URL.Query().Get("limit")
	limit := 10
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}

	ops, err := h.svc.ListRecent(tenantID, limit)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, ops, http.StatusOK)
}

// Delete removes an operation
// @Summary Delete an operation
// @Description Delete an operation by ID
// @Tags operations
// @Param tenant_id path string true "Tenant ID"
// @Param id path string true "Operation ID"
// @Success 204 "No Content"
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{tenant_id}/operations/{id} [delete]
func (h *OperationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.svc.Delete(id); err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
