package handler

import (
	"encoding/json"
	"net/http"

	"github.com/aralvesandrade/cafeos/internal/api/middleware"
	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	"github.com/aralvesandrade/cafeos/internal/domain/service"
)

// Types aliases for swagger documentation
type (
	SwaggerFarm              = entity.Farm
	SwaggerPlot              = entity.Plot
	SwaggerOperation         = entity.Operation
	SwaggerHarvest           = entity.Harvest
	SwaggerHarvestProduction = entity.HarvestProduction
	SwaggerIndicator         = entity.Indicator
)

type FarmHandler struct {
	svc *service.FarmService
}

func NewFarmHandler(svc *service.FarmService) *FarmHandler {
	return &FarmHandler{svc: svc}
}

type createFarmRequest struct {
	Name         string  `json:"name"`
	Owner        string  `json:"owner"`
	Location     string  `json:"location"`
	TotalAreaHA  float64 `json:"total_area_ha"`
	PlantedAreaHA float64 `json:"planted_area_ha"`
}

// Create registers a new farm for the authenticated tenant
// @Summary Create a farm
// @Description Register a new farm in the tenant
// @Tags farms
// @Accept json
// @Produce json
// @Param tenant_id path string true "Tenant ID"
// @Param farm body createFarmRequest true "Farm data"
// @Success 201 {object} SwaggerFarm
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{tenant_id}/farms [post]
func (h *FarmHandler) Create(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := r.Context().Value(middleware.TenantIDKey).(string)

	var req createFarmRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	farm, err := h.svc.Create(tenantID, req.Name, req.Owner, req.Location, req.TotalAreaHA, req.PlantedAreaHA)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, farm, http.StatusCreated)
}

// GetByID returns a farm by its ID
// @Summary Get farm by ID
// @Description Returns a single farm
// @Tags farms
// @Produce json
// @Param tenant_id path string true "Tenant ID"
// @Param id path string true "Farm ID"
// @Success 200 {object} SwaggerFarm
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{tenant_id}/farms/{id} [get]
func (h *FarmHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	farm, err := h.svc.GetByID(id)
	if err != nil {
		writeError(w, "farm not found", http.StatusNotFound)
		return
	}
	writeJSON(w, farm, http.StatusOK)
}

// List returns all farms for the authenticated tenant
// @Summary List farms
// @Description List all farms belonging to the tenant
// @Tags farms
// @Produce json
// @Param tenant_id path string true "Tenant ID"
// @Success 200 {array} SwaggerFarm
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{tenant_id}/farms [get]
func (h *FarmHandler) List(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := r.Context().Value(middleware.TenantIDKey).(string)
	farms, err := h.svc.ListByTenant(tenantID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, farms, http.StatusOK)
}

// Update updates an existing farm
// @Summary Update a farm
// @Description Update farm data by ID
// @Tags farms
// @Accept json
// @Produce json
// @Param tenant_id path string true "Tenant ID"
// @Param id path string true "Farm ID"
// @Param farm body SwaggerFarm true "Updated farm data"
// @Success 200 {object} SwaggerFarm
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{tenant_id}/farms/{id} [put]
func (h *FarmHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var farm entity.Farm
	if err := json.NewDecoder(r.Body).Decode(&farm); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	farm.ID = id

	if err := h.svc.Update(&farm); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, farm, http.StatusOK)
}

// Delete removes a farm
// @Summary Delete a farm
// @Description Delete a farm by ID
// @Tags farms
// @Param tenant_id path string true "Tenant ID"
// @Param id path string true "Farm ID"
// @Success 204 "No Content"
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{tenant_id}/farms/{id} [delete]
func (h *FarmHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.svc.Delete(id); err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
