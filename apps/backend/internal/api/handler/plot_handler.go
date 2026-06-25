package handler

import (
	"encoding/json"
	"net/http"

	"github.com/aralvesandrade/cafeos/internal/api/middleware"
	"github.com/aralvesandrade/cafeos/internal/domain/service"
)

type PlotHandler struct {
	svc *service.PlotService
}

func NewPlotHandler(svc *service.PlotService) *PlotHandler {
	return &PlotHandler{svc: svc}
}

type createPlotRequest struct {
	FarmID       string  `json:"farm_id"`
	Name         string  `json:"name"`
	AreaHA       float64 `json:"area_ha"`
	Cultivar     string  `json:"cultivar"`
	PlantingYear int     `json:"planting_year"`
	Altitude     int     `json:"altitude"`
	SoilType     string  `json:"soil_type"`
}

// Create registers a new plot for the authenticated tenant
// @Summary Create a plot
// @Description Register a new plot (talhão) in a farm
// @Tags plots
// @Accept json
// @Produce json
// @Param tenant_id path string true "Tenant ID"
// @Param plot body createPlotRequest true "Plot data"
// @Success 201 {object} SwaggerPlot
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{tenant_id}/plots [post]
func (h *PlotHandler) Create(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := r.Context().Value(middleware.TenantIDKey).(string)

	var req createPlotRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	plot, err := h.svc.Create(tenantID, req.FarmID, req.Name, req.Cultivar, req.SoilType, req.AreaHA, req.PlantingYear, req.Altitude)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, plot, http.StatusCreated)
}

// GetByID returns a plot by its ID
// @Summary Get plot by ID
// @Description Returns a single plot (talhão)
// @Tags plots
// @Produce json
// @Param tenant_id path string true "Tenant ID"
// @Param id path string true "Plot ID"
// @Success 200 {object} SwaggerPlot
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{tenant_id}/plots/{id} [get]
func (h *PlotHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	plot, err := h.svc.GetByID(id)
	if err != nil {
		writeError(w, "plot not found", http.StatusNotFound)
		return
	}
	writeJSON(w, plot, http.StatusOK)
}

// ListByFarm returns all plots for a given farm
// @Summary List plots by farm
// @Description List all plots (talhões) belonging to a farm
// @Tags plots
// @Produce json
// @Param tenant_id path string true "Tenant ID"
// @Param farm_id path string true "Farm ID"
// @Success 200 {array} SwaggerPlot
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{tenant_id}/farms/{farm_id}/plots [get]
func (h *PlotHandler) ListByFarm(w http.ResponseWriter, r *http.Request) {
	farmID := r.PathValue("farm_id")
	plots, err := h.svc.ListByFarm(farmID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, plots, http.StatusOK)
}

// List returns all plots for the authenticated tenant
// @Summary List all plots
// @Description List all plots across all farms in the tenant
// @Tags plots
// @Produce json
// @Param tenant_id path string true "Tenant ID"
// @Success 200 {array} SwaggerPlot
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{tenant_id}/plots [get]
func (h *PlotHandler) List(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := r.Context().Value(middleware.TenantIDKey).(string)
	plots, err := h.svc.ListByTenant(tenantID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, plots, http.StatusOK)
}

// Update updates an existing plot
// @Summary Update a plot
// @Description Update plot data by ID (partial update - only provided fields are changed)
// @Tags plots
// @Accept json
// @Produce json
// @Param tenant_id path string true "Tenant ID"
// @Param id path string true "Plot ID"
// @Param plot body SwaggerPlot true "Updated plot data"
// @Success 200 {object} SwaggerPlot
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{tenant_id}/plots/{id} [put]
func (h *PlotHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	existing, err := h.svc.GetByID(id)
	if err != nil {
		writeError(w, "plot not found", http.StatusNotFound)
		return
	}

	var input struct {
		Name         *string  `json:"name"`
		FarmID       *string  `json:"farm_id"`
		AreaHA       *float64 `json:"area_ha"`
		Cultivar     *string  `json:"cultivar"`
		SoilType     *string  `json:"soil_type"`
		Altitude     *int     `json:"altitude"`
		PlantingYear *int     `json:"planting_year"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if input.Name != nil {
		existing.Name = *input.Name
	}
	if input.FarmID != nil {
		existing.FarmID = *input.FarmID
	}
	if input.AreaHA != nil {
		existing.AreaHA = *input.AreaHA
	}
	if input.Cultivar != nil {
		existing.Cultivar = *input.Cultivar
	}
	if input.SoilType != nil {
		existing.SoilType = *input.SoilType
	}
	if input.Altitude != nil {
		existing.Altitude = *input.Altitude
	}
	if input.PlantingYear != nil {
		existing.PlantingYear = *input.PlantingYear
	}

	if err := h.svc.Update(existing); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, existing, http.StatusOK)
}

// Delete removes a plot
// @Summary Delete a plot
// @Description Delete a plot by ID
// @Tags plots
// @Param tenant_id path string true "Tenant ID"
// @Param id path string true "Plot ID"
// @Success 204 "No Content"
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{tenant_id}/plots/{id} [delete]
func (h *PlotHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.svc.Delete(id); err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
