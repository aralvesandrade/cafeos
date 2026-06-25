package handler

import (
	"encoding/json"
	"net/http"

	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	"github.com/aralvesandrade/cafeos/internal/domain/repository"
)

type DashboardHandler struct {
	harvestRepo    repository.HarvestRepository
	indicatorRepo  repository.IndicatorRepository
	operationRepo  repository.OperationRepository
	plotRepo       repository.PlotRepository
}

func NewDashboardHandler(
	harvestRepo repository.HarvestRepository,
	indicatorRepo repository.IndicatorRepository,
	operationRepo repository.OperationRepository,
	plotRepo repository.PlotRepository,
) *DashboardHandler {
	return &DashboardHandler{
		harvestRepo:   harvestRepo,
		indicatorRepo: indicatorRepo,
		operationRepo: operationRepo,
		plotRepo:      plotRepo,
	}
}

type DashboardResponse struct {
	TotalProduction  float64             `json:"total_production"`
	TotalCost        float64             `json:"total_cost"`
	SacasPerHA       float64             `json:"sacas_per_ha"`
	CostPerSaca      float64             `json:"cost_per_saca"`
	Harvests         []*entity.Harvest   `json:"harvests"`
	RecentOperations []*entity.Operation `json:"recent_operations"`
	Indicators       []*entity.Indicator `json:"indicators"`
}

// GetDashboard returns consolidated dashboard data
// @Summary Get dashboard
// @Description Returns consolidated production, costs, indicators and recent operations
// @Tags dashboard
// @Produce json
// @Param tenant_id path string true "Tenant ID"
// @Success 200 {object} DashboardResponse
// @Security BearerAuth
// @Router /api/v1/{tenant_id}/dashboard [get]
func (h *DashboardHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	tenantID := r.PathValue("tenant_id")
	if tenantID == "" {
		tenantID = r.URL.Query().Get("tenant_id")
	}

	harvests, _ := h.harvestRepo.ListByTenant(tenantID)
	indicators, _ := h.indicatorRepo.ListByTenant(tenantID)
	operations, _ := h.operationRepo.ListByTenant(tenantID)
	plots, _ := h.plotRepo.ListByTenant(tenantID)

	var totalProduction, totalCost float64
	var latestIndicators []*entity.Indicator

	for _, ind := range indicators {
		switch ind.Type {
		case entity.IndProducaoTotal:
			totalProduction = ind.Value
		case entity.IndCustoTotal:
			totalCost = ind.Value
		}
	}

	var totalPlantedArea float64
	for _, p := range plots {
		totalPlantedArea += p.AreaHA
	}

	sacasPerHA := 0.0
	if totalPlantedArea > 0 {
		sacasPerHA = totalProduction / totalPlantedArea
	}

	costPerSaca := 0.0
	if totalProduction > 0 {
		costPerSaca = totalCost / totalProduction
	}

	if len(indicators) > 0 {
		seen := make(map[string]bool)
		for _, ind := range indicators {
			if !seen[string(ind.Type)] {
				latestIndicators = append(latestIndicators, ind)
				seen[string(ind.Type)] = true
			}
		}
	}

	recentLimit := 10
	if len(operations) > recentLimit {
		operations = operations[:recentLimit]
	}

	dashboard := DashboardResponse{
		TotalProduction:  totalProduction,
		TotalCost:        totalCost,
		SacasPerHA:       sacasPerHA,
		CostPerSaca:      costPerSaca,
		Harvests:         harvests,
		RecentOperations: operations,
		Indicators:       latestIndicators,
	}

	writeJSON(w, dashboard, http.StatusOK)
}

func writeJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
