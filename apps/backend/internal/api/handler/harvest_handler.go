package handler

import (
	"encoding/json"
	"net/http"

	"github.com/aralvesandrade/cafeos/internal/api/middleware"
	"github.com/aralvesandrade/cafeos/internal/domain/service"
)

type HarvestHandler struct {
	svc *service.HarvestService
}

func NewHarvestHandler(svc *service.HarvestService) *HarvestHandler {
	return &HarvestHandler{svc: svc}
}

type createHarvestRequest struct {
	Year                int     `json:"year"`
	Description         string  `json:"description"`
	EstimatedProduction float64 `json:"estimated_production"`
}

// Create registers a new harvest season
// @Summary Create a harvest
// @Description Register a new harvest (safra) for a year
// @Tags harvests
// @Accept json
// @Produce json
// @Param tenant_id path string true "Tenant ID"
// @Param harvest body createHarvestRequest true "Harvest data"
// @Success 201 {object} SwaggerHarvest
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{tenant_id}/harvests [post]
func (h *HarvestHandler) Create(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := r.Context().Value(middleware.TenantIDKey).(string)

	var req createHarvestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	harvest, err := h.svc.Create(tenantID, req.Year, req.Description, req.EstimatedProduction)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, harvest, http.StatusCreated)
}

// GetByID returns a harvest by its ID
// @Summary Get harvest by ID
// @Description Returns a single harvest (safra)
// @Tags harvests
// @Produce json
// @Param tenant_id path string true "Tenant ID"
// @Param id path string true "Harvest ID"
// @Success 200 {object} SwaggerHarvest
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{tenant_id}/harvests/{id} [get]
func (h *HarvestHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	harvest, err := h.svc.GetByID(id)
	if err != nil {
		writeError(w, "harvest not found", http.StatusNotFound)
		return
	}
	writeJSON(w, harvest, http.StatusOK)
}

// List returns all harvests for the authenticated tenant
// @Summary List harvests
// @Description List all harvests (safras) for the tenant
// @Tags harvests
// @Produce json
// @Param tenant_id path string true "Tenant ID"
// @Success 200 {array} SwaggerHarvest
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{tenant_id}/harvests [get]
func (h *HarvestHandler) List(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := r.Context().Value(middleware.TenantIDKey).(string)
	harvests, err := h.svc.ListByTenant(tenantID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, harvests, http.StatusOK)
}

// Finalize marks a harvest as completed and calculates indicators
// @Summary Finalize a harvest
// @Description Finalize a harvest (safra), calculating all indicators (sacas/ha, custo/saca, etc.)
// @Tags harvests
// @Param tenant_id path string true "Tenant ID"
// @Param id path string true "Harvest ID"
// @Success 200 "OK"
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{tenant_id}/harvests/{id}/finalize [put]
func (h *HarvestHandler) Finalize(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.svc.Finalize(id); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

type recordProductionRequest struct {
	PlotID   string  `json:"plot_id"`
	Quantity float64 `json:"quantity"`
	Notes    string  `json:"notes"`
}

// RecordProduction records production for a plot in a harvest
// @Summary Record harvest production
// @Description Record production quantity for a specific plot in a harvest
// @Tags harvests
// @Accept json
// @Produce json
// @Param tenant_id path string true "Tenant ID"
// @Param id path string true "Harvest ID"
// @Param production body recordProductionRequest true "Production data"
// @Success 201 {object} SwaggerHarvestProduction
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{tenant_id}/harvests/{id}/production [post]
func (h *HarvestHandler) RecordProduction(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := r.Context().Value(middleware.TenantIDKey).(string)
	harvestID := r.PathValue("id")

	var req recordProductionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	prod, err := h.svc.RecordProduction(tenantID, harvestID, req.PlotID, req.Quantity, req.Notes)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, prod, http.StatusCreated)
}

// GetProduction returns all production records for a harvest
// @Summary Get harvest production
// @Description Get all production records for a given harvest
// @Tags harvests
// @Produce json
// @Param tenant_id path string true "Tenant ID"
// @Param id path string true "Harvest ID"
// @Success 200 {array} SwaggerHarvestProduction
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{tenant_id}/harvests/{id}/production [get]
func (h *HarvestHandler) GetProduction(w http.ResponseWriter, r *http.Request) {
	harvestID := r.PathValue("id")
	productions, err := h.svc.GetProductionByHarvest(harvestID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, productions, http.StatusOK)
}
