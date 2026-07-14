package handler

import (
	"encoding/json"
	"net/http"

	"github.com/aralvesandrade/cafeos/internal/api/middleware"
	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	"github.com/aralvesandrade/cafeos/internal/domain/repository"
	"github.com/aralvesandrade/cafeos/internal/domain/service"
)

type DashboardHandler struct {
	harvestRepo    repository.HarvestRepository
	indicatorRepo  repository.IndicatorRepository
	operationRepo  repository.OperationRepository
	plotRepo       repository.PlotRepository
	farmRepo       repository.FarmRepository
	productionRepo repository.HarvestProductionRepository
	farmSvc        *service.FarmService
}

func NewDashboardHandler(
	harvestRepo repository.HarvestRepository,
	indicatorRepo repository.IndicatorRepository,
	operationRepo repository.OperationRepository,
	plotRepo repository.PlotRepository,
	farmRepo repository.FarmRepository,
	productionRepo repository.HarvestProductionRepository,
	farmSvc *service.FarmService,
) *DashboardHandler {
	return &DashboardHandler{
		harvestRepo:    harvestRepo,
		indicatorRepo:  indicatorRepo,
		operationRepo:  operationRepo,
		plotRepo:       plotRepo,
		farmRepo:       farmRepo,
		productionRepo: productionRepo,
		farmSvc:        farmSvc,
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
// @Param farm_id query string false "Filtrar por ID da Fazenda"
// @Param harvest_id query string false "Filtrar por ID da Colheita"
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
	operations, _ := h.operationRepo.ListByOrganization(organizationID)
	plots, _ := h.plotRepo.ListByOrganization(organizationID)
	farms, _ := h.farmRepo.ListByOrganization(organizationID)

	scopeFarmIDs, err := h.resolveFarmScope(r, organizationID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	harvestID := r.URL.Query().Get("harvest_id")

	if scopeFarmIDs == nil && harvestID == "" {
		writeJSON(w, h.buildUnscopedDashboard(organizationID, harvests, operations, plots, farms), http.StatusOK)
		return
	}

	writeJSON(w, h.buildScopedDashboard(scopeFarmIDs, harvestID, harvests, operations, plots, farms), http.StatusOK)
}

// resolveFarmScope returns the set of farm IDs the response should be limited
// to, or nil for no scoping (full organization). proprietario is always
// scoped to its own farms; every other role is scoped only if ?farm_id= is given.
func (h *DashboardHandler) resolveFarmScope(r *http.Request, organizationID string) (map[string]bool, error) {
	farmIDParam := r.URL.Query().Get("farm_id")

	if userID, restricted := restrictedOwnerID(r); restricted {
		owned, err := h.farmSvc.OwnedFarmIDs(organizationID, userID)
		if err != nil {
			return nil, err
		}
		if farmIDParam != "" {
			if owned[farmIDParam] {
				return map[string]bool{farmIDParam: true}, nil
			}
			return map[string]bool{}, nil
		}
		return owned, nil
	}

	if farmIDParam != "" {
		return map[string]bool{farmIDParam: true}, nil
	}
	return nil, nil
}

// buildUnscopedDashboard is the original, unfiltered dashboard: totals come
// from the precomputed per-harvest Indicator rows. Kept byte-for-byte
// equivalent to the pre-filter behavior.
func (h *DashboardHandler) buildUnscopedDashboard(organizationID string, harvests []*entity.Harvest, operations []*entity.Operation, plots []*entity.Plot, farms []*entity.Farm) DashboardResponse {
	indicators, _ := h.indicatorRepo.ListByOrganization(organizationID)

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
	costPerBag := make([]CostPerBag, 0, len(harvests))
	for _, hv := range harvests {
		productionByHarvest = append(productionByHarvest, ProductionByHarvest{Year: hv.Description, Production: hv.EstimatedProduction})
		costPerBag = append(costPerBag, CostPerBag{Year: hv.Description, Cost: 0})
	}

	plotMap := make(map[string]string, len(plots))
	for _, p := range plots {
		plotMap[p.ID] = p.Name
	}

	if len(operations) > 10 {
		operations = operations[:10]
	}
	recentOps := make([]RecentOperationItem, 0, len(operations))
	for _, op := range operations {
		recentOps = append(recentOps, RecentOperationItem{
			ID: op.ID, Type: op.TypeName, Date: op.Date.Format("2006-01-02"),
			PlotName: plotMap[op.PlotID], Cost: op.Cost,
		})
	}

	return DashboardResponse{
		TotalFarms:          len(farms),
		TotalPlots:          len(plots),
		TotalProduction:     totalProduction,
		TotalCost:           totalCost,
		ProductionByHarvest: productionByHarvest,
		CostPerBag:          costPerBag,
		RecentOperations:    recentOps,
	}
}

// buildScopedDashboard recomputes every total from plot-level source data
// (HarvestProduction, Operation) restricted to scopeFarmIDs (nil = every farm)
// and, if set, a single harvest — since the precomputed Indicator rows have
// no per-farm breakdown.
func (h *DashboardHandler) buildScopedDashboard(scopeFarmIDs map[string]bool, harvestID string, harvests []*entity.Harvest, operations []*entity.Operation, plots []*entity.Plot, farms []*entity.Farm) DashboardResponse {
	if scopeFarmIDs != nil {
		filteredFarms := farms[:0]
		for _, f := range farms {
			if scopeFarmIDs[f.ID] {
				filteredFarms = append(filteredFarms, f)
			}
		}
		farms = filteredFarms

		filteredPlots := plots[:0]
		for _, p := range plots {
			if scopeFarmIDs[p.FarmID] {
				filteredPlots = append(filteredPlots, p)
			}
		}
		plots = filteredPlots
	}

	plotIDs := make(map[string]bool, len(plots))
	plotMap := make(map[string]string, len(plots))
	for _, p := range plots {
		plotIDs[p.ID] = true
		plotMap[p.ID] = p.Name
	}

	scopedHarvests := harvests
	if harvestID != "" {
		scopedHarvests = scopedHarvests[:0]
		for _, hv := range harvests {
			if hv.ID == harvestID {
				scopedHarvests = append(scopedHarvests, hv)
			}
		}
	}

	filteredOps := operations[:0]
	for _, op := range operations {
		if !plotIDs[op.PlotID] {
			continue
		}
		if harvestID != "" && (op.HarvestID == nil || *op.HarvestID != harvestID) {
			continue
		}
		filteredOps = append(filteredOps, op)
	}

	// Total cost sums every operation in scope (plot + optional harvest filter),
	// regardless of whether each operation is itself linked to a harvest —
	// operations commonly aren't tagged with a HarvestID in practice.
	var totalCost float64
	for _, op := range filteredOps {
		totalCost += op.Cost
	}

	var totalProduction float64
	productionByHarvest := make([]ProductionByHarvest, 0, len(scopedHarvests))
	costPerBag := make([]CostPerBag, 0, len(scopedHarvests))
	for _, hv := range scopedHarvests {
		productions, _ := h.productionRepo.ListByHarvest(hv.ID)
		var hvProduction float64
		for _, p := range productions {
			if plotIDs[p.PlotID] {
				hvProduction += p.Quantity
			}
		}

		var hvCost float64
		for _, op := range filteredOps {
			if op.HarvestID != nil && *op.HarvestID == hv.ID {
				hvCost += op.Cost
			}
		}

		totalProduction += hvProduction
		productionByHarvest = append(productionByHarvest, ProductionByHarvest{Year: hv.Description, Production: hvProduction})
		costPerBag = append(costPerBag, CostPerBag{Year: hv.Description, Cost: hvCost})
	}

	if len(filteredOps) > 10 {
		filteredOps = filteredOps[:10]
	}
	recentOps := make([]RecentOperationItem, 0, len(filteredOps))
	for _, op := range filteredOps {
		recentOps = append(recentOps, RecentOperationItem{
			ID: op.ID, Type: op.TypeName, Date: op.Date.Format("2006-01-02"),
			PlotName: plotMap[op.PlotID], Cost: op.Cost,
		})
	}

	return DashboardResponse{
		TotalFarms:          len(farms),
		TotalPlots:          len(plots),
		TotalProduction:     totalProduction,
		TotalCost:           totalCost,
		ProductionByHarvest: productionByHarvest,
		CostPerBag:          costPerBag,
		RecentOperations:    recentOps,
	}
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
