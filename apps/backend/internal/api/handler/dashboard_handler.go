package handler

import (
	"encoding/json"
	"net/http"

	"github.com/aralvesandrade/cafeos/internal/api/middleware"
	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	"github.com/aralvesandrade/cafeos/internal/domain/repository"
)

type DashboardHandler struct {
	harvestRepo   repository.HarvestRepository
	indicatorRepo repository.IndicatorRepository
	operationRepo repository.OperationRepository
	plotRepo      repository.PlotRepository
	farmRepo      repository.FarmRepository
}

func NewDashboardHandler(
	harvestRepo repository.HarvestRepository,
	indicatorRepo repository.IndicatorRepository,
	operationRepo repository.OperationRepository,
	plotRepo repository.PlotRepository,
	farmRepo repository.FarmRepository,
) *DashboardHandler {
	return &DashboardHandler{
		harvestRepo:   harvestRepo,
		indicatorRepo: indicatorRepo,
		operationRepo: operationRepo,
		plotRepo:      plotRepo,
		farmRepo:      farmRepo,
	}
}

type ProductionByHarvest struct {
	Year       string  `json:"year"`
	Production float64 `json:"production"`
}

type CostPerBag struct {
	Year string  `json:"year"`
	Cost float64 `json:"cost"`
}

type RecentOperationItem struct {
	ID       string  `json:"id"`
	Type     string  `json:"type"`
	Date     string  `json:"date"`
	PlotName string  `json:"plot_name"`
	Cost     float64 `json:"cost"`
}

type DashboardResponse struct {
	TotalFarms          int                   `json:"total_farms"`
	TotalPlots          int                   `json:"total_plots"`
	TotalProduction     float64               `json:"total_production"`
	TotalCost           float64               `json:"total_cost"`
	ProductionByHarvest []ProductionByHarvest `json:"production_by_harvest"`
	CostPerBag          []CostPerBag          `json:"cost_per_bag"`
	RecentOperations    []RecentOperationItem `json:"recent_operations"`
}

// GetDashboard retorna os dados consolidados do painel
// @Summary Obter painel
// @Description Retorna produção, custos, indicadores e operações recentes consolidados
// @Tags dashboard (Painel)
// @Produce json
// @Param organization_id path string true "ID da Organização"
// @Success 200 {object} DashboardResponse
// @Security BearerAuth
// @Router /api/v1/{organization_id}/dashboard [get]
func (h *DashboardHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	organizationID, _ := r.Context().Value(middleware.OrganizationIDKey).(string)
	if organizationID == "" {
		writeError(w, "organization not found", http.StatusUnauthorized)
		return
	}

	harvests, _ := h.harvestRepo.ListByOrganization(organizationID)
	indicators, _ := h.indicatorRepo.ListByOrganization(organizationID)
	operations, _ := h.operationRepo.ListByOrganization(organizationID)
	plots, _ := h.plotRepo.ListByOrganization(organizationID)
	farms, _ := h.farmRepo.ListByOrganization(organizationID)

	var totalProduction, totalCost float64

	for _, ind := range indicators {
		switch ind.Type {
		case entity.IndProducaoTotal:
			totalProduction = ind.Value
		case entity.IndCustoTotal:
			totalCost = ind.Value
		}
	}

	productionByHarvest := make([]ProductionByHarvest, 0, len(harvests))
	for _, h := range harvests {
		productionByHarvest = append(productionByHarvest, ProductionByHarvest{
			Year:       h.Description,
			Production: h.EstimatedProduction,
		})
	}

	plotMap := make(map[string]string)
	for _, p := range plots {
		plotMap[p.ID] = p.Name
	}

	costPerBag := make([]CostPerBag, 0, len(harvests))
	for _, h := range harvests {
		costPerBag = append(costPerBag, CostPerBag{
			Year: h.Description,
			Cost: 0,
		})
	}

	recentLimit := 10
	if len(operations) > recentLimit {
		operations = operations[:recentLimit]
	}

	recentOps := make([]RecentOperationItem, 0, len(operations))
	for _, op := range operations {
		recentOps = append(recentOps, RecentOperationItem{
			ID:       op.ID,
			Type:     string(op.Type),
			Date:     op.Date.Format("2006-01-02"),
			PlotName: plotMap[op.PlotID],
			Cost:     op.Cost,
		})
	}

	dashboard := DashboardResponse{
		TotalFarms:          len(farms),
		TotalPlots:          len(plots),
		TotalProduction:     totalProduction,
		TotalCost:           totalCost,
		ProductionByHarvest: productionByHarvest,
		CostPerBag:          costPerBag,
		RecentOperations:    recentOps,
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
