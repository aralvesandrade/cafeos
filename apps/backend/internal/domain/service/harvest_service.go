package service

import (
	"errors"
	"time"

	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	"github.com/aralvesandrade/cafeos/internal/domain/repository"
	"github.com/aralvesandrade/cafeos/internal/event"
	"github.com/google/uuid"
)

type HarvestService struct {
	harvestRepo     repository.HarvestRepository
	productionRepo  repository.HarvestProductionRepository
	indicatorRepo   repository.IndicatorRepository
	plotRepo        repository.PlotRepository
	operationRepo   repository.OperationRepository
	maintenanceRepo repository.MaintenanceRepository
	shiftRepo       repository.WorkShiftRepository
	financialRepo   repository.FinancialRepository
	allocationRepo  repository.CostAllocationRepository
	transactor      repository.Transactor
	eventBus        event.Bus
}

func NewHarvestService(
	harvestRepo repository.HarvestRepository,
	productionRepo repository.HarvestProductionRepository,
	indicatorRepo repository.IndicatorRepository,
	plotRepo repository.PlotRepository,
	operationRepo repository.OperationRepository,
	maintenanceRepo repository.MaintenanceRepository,
	shiftRepo repository.WorkShiftRepository,
	financialRepo repository.FinancialRepository,
	allocationRepo repository.CostAllocationRepository,
	transactor repository.Transactor,
	eventBus event.Bus,
) *HarvestService {
	return &HarvestService{
		harvestRepo:     harvestRepo,
		productionRepo:  productionRepo,
		indicatorRepo:   indicatorRepo,
		plotRepo:        plotRepo,
		operationRepo:   operationRepo,
		maintenanceRepo: maintenanceRepo,
		shiftRepo:       shiftRepo,
		financialRepo:   financialRepo,
		allocationRepo:  allocationRepo,
		transactor:      transactor,
		eventBus:        eventBus,
	}
}

func (s *HarvestService) Create(organizationID string, year int, description string, estimatedProduction float64) (*entity.Harvest, error) {
	if year <= 0 {
		return nil, errors.New("year is required")
	}
	if estimatedProduction < 0 {
		return nil, errors.New("estimated production cannot be negative")
	}

	harvest := &entity.Harvest{
		ID:                  uuid.New().String(),
		OrganizationID:      organizationID,
		Year:                year,
		Description:         description,
		EstimatedProduction: estimatedProduction,
		Status:              entity.HarvestPlanejada,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	if err := s.harvestRepo.Create(harvest); err != nil {
		return nil, err
	}
	return harvest, nil
}

func (s *HarvestService) GetByID(id string) (*entity.Harvest, error) {
	return s.harvestRepo.GetByID(id)
}

func (s *HarvestService) ListByOrganization(organizationID string) ([]*entity.Harvest, error) {
	return s.harvestRepo.ListByOrganization(organizationID)
}

func (s *HarvestService) Finalize(harvestID string) error {
	harvest, err := s.harvestRepo.GetByID(harvestID)
	if err != nil {
		return err
	}

	harvest.Status = entity.HarvestFinalizada
	harvest.UpdatedAt = time.Now()

	err = s.transactor.RunInTx(func(repos repository.TransactionProvider) error {
		if err := repos.Harvest().Update(harvest); err != nil {
			return err
		}
		return s.calculateIndicatorsWithRepos(harvest, repos)
	})
	if err != nil {
		return err
	}

	s.eventBus.Publish(event.HarvestFinalized{
		HarvestID:      harvestID,
		OrganizationID: harvest.OrganizationID,
		Year:           harvest.Year,
	})

	return nil
}

func (s *HarvestService) RecordProduction(organizationID, harvestID, plotID string, quantity float64, notes string) (*entity.HarvestProduction, error) {
	if quantity < 0 {
		return nil, errors.New("quantity cannot be negative")
	}

	prod := &entity.HarvestProduction{
		ID:             uuid.New().String(),
		OrganizationID: organizationID,
		HarvestID:      harvestID,
		PlotID:         plotID,
		Quantity:       quantity,
		RecordedAt:     time.Now(),
		Notes:          notes,
	}

	if err := s.productionRepo.Create(prod); err != nil {
		return nil, err
	}

	return prod, nil
}

func (s *HarvestService) GetProductionByHarvest(harvestID string) ([]*entity.HarvestProduction, error) {
	return s.productionRepo.ListByHarvest(harvestID)
}

func (s *HarvestService) calculateIndicators(harvest *entity.Harvest) error {
	productions, err := s.productionRepo.ListByHarvest(harvest.ID)
	if err != nil {
		return err
	}

	var totalProduction float64
	for _, p := range productions {
		totalProduction += p.Quantity
	}

	plots, err := s.plotRepo.ListByOrganization(harvest.OrganizationID)
	if err != nil {
		return err
	}

	var totalPlantedArea float64
	for _, p := range plots {
		totalPlantedArea += p.AreaHA
	}

	totalCost := s.calculateTotalCost(harvest,
		s.operationRepo.ListByOrganization,
		s.maintenanceRepo.ListByOrganization,
		s.shiftRepo.ListByOrganization,
		s.financialRepo.ListByOrganization,
		s.allocationRepo.ListByHarvest,
	)

	indicators := s.buildIndicators(harvest, totalProduction, totalPlantedArea, totalCost, map[entity.CostGroup]float64{})

	for _, ind := range indicators {
		if err := s.indicatorRepo.Create(ind); err != nil {
			return err
		}
	}

	return nil
}

func (s *HarvestService) calculateIndicatorsWithRepos(harvest *entity.Harvest, repos repository.TransactionProvider) error {
	productions, err := repos.HarvestProduction().ListByHarvest(harvest.ID)
	if err != nil {
		return err
	}

	var totalProduction float64
	for _, p := range productions {
		totalProduction += p.Quantity
	}

	plots, err := repos.Plot().ListByOrganization(harvest.OrganizationID)
	if err != nil {
		return err
	}

	var totalPlantedArea float64
	for _, p := range plots {
		totalPlantedArea += p.AreaHA
	}

	totalCost := s.calculateTotalCostWithRepos(harvest, repos)
	costByGroup, err := s.calculateCostByGroupWithRepos(harvest, repos)
	if err != nil {
		return err
	}

	indicators := s.buildIndicators(harvest, totalProduction, totalPlantedArea, totalCost, costByGroup)

	for _, ind := range indicators {
		if err := repos.Indicator().Create(ind); err != nil {
			return err
		}
	}

	return nil
}

func (s *HarvestService) calculateTotalCost(harvest *entity.Harvest,
	listOps func(string) ([]*entity.Operation, error),
	listMaint func(string) ([]*entity.Maintenance, error),
	listShifts func(string) ([]*entity.WorkShift, error),
	listFin func(string) ([]*entity.FinancialTransaction, error),
	listAllocs func(string) ([]*entity.CostAllocation, error),
) float64 {
	var totalCost float64

	// Operations by HarvestID
	ops, _ := listOps(harvest.OrganizationID)
	for _, op := range ops {
		if op.HarvestID != nil && *op.HarvestID == harvest.ID {
			totalCost += op.Cost
		}
	}

	start := harvest.StartDate
	end := harvest.EndDate

	// Maintenance by date range
	maints, _ := listMaint(harvest.OrganizationID)
	for _, m := range maints {
		if inDateRange(m.Date, start, end) {
			totalCost += m.Cost
		}
	}

	// WorkShift by date range
	shifts, _ := listShifts(harvest.OrganizationID)
	for _, ws := range shifts {
		if inDateRange(ws.Date, start, end) {
			totalCost += ws.Cost
		}
	}

	// FinancialTransaction (despesas) by date range
	transactions, _ := listFin(harvest.OrganizationID)
	for _, ft := range transactions {
		if ft.Type == entity.TranDespesa && inDateRange(ft.Date, start, end) {
			totalCost += ft.Amount
		}
	}

	// CostAllocation by HarvestID
	allocs, _ := listAllocs(harvest.ID)
	for _, a := range allocs {
		totalCost += a.TotalAmount
	}

	return totalCost
}

func (s *HarvestService) calculateTotalCostWithRepos(harvest *entity.Harvest, repos repository.TransactionProvider) float64 {
	var totalCost float64

	ops, _ := repos.Operation().ListByOrganization(harvest.OrganizationID)
	for _, op := range ops {
		if op.HarvestID != nil && *op.HarvestID == harvest.ID {
			totalCost += op.Cost
		}
	}

	start := harvest.StartDate
	end := harvest.EndDate

	maints, _ := repos.Maintenance().ListByOrganization(harvest.OrganizationID)
	for _, m := range maints {
		if inDateRange(m.Date, start, end) {
			totalCost += m.Cost
		}
	}

	shifts, _ := repos.WorkShift().ListByOrganization(harvest.OrganizationID)
	for _, ws := range shifts {
		if inDateRange(ws.Date, start, end) {
			totalCost += ws.Cost
		}
	}

	transactions, _ := repos.Financial().ListByOrganization(harvest.OrganizationID)
	for _, ft := range transactions {
		if ft.Type == entity.TranDespesa && inDateRange(ft.Date, start, end) {
			totalCost += ft.Amount
		}
	}

	allocs, _ := repos.CostAllocation().ListByHarvest(harvest.ID)
	for _, a := range allocs {
		totalCost += a.TotalAmount
	}

	return totalCost
}

// calculateCostByGroupWithRepos mirrors calculateTotalCostWithRepos but
// buckets each cost item by its CostCenter.CostGroup (SENAR/CEPEA
// classification), so COE/COT/CT can be computed. Cost items without a
// CostCenterID, or whose CostCenter has no CostGroup set, are not counted
// in any bucket (but still count towards the legacy custo_total).
func (s *HarvestService) calculateCostByGroupWithRepos(harvest *entity.Harvest, repos repository.TransactionProvider) (map[entity.CostGroup]float64, error) {
	costCenters, err := repos.CostCenter().ListByOrganization(harvest.OrganizationID)
	if err != nil {
		return nil, err
	}
	groupByCostCenter := make(map[string]entity.CostGroup, len(costCenters))
	for _, cc := range costCenters {
		groupByCostCenter[cc.ID] = cc.CostGroup
	}

	byGroup := map[entity.CostGroup]float64{}
	add := func(costCenterID *string, amount float64) {
		if costCenterID == nil {
			return
		}
		group, ok := groupByCostCenter[*costCenterID]
		if !ok || group == "" {
			return
		}
		byGroup[group] += amount
	}

	ops, _ := repos.Operation().ListByOrganization(harvest.OrganizationID)
	for _, op := range ops {
		if op.HarvestID != nil && *op.HarvestID == harvest.ID {
			add(op.CostCenterID, op.Cost)
		}
	}

	start := harvest.StartDate
	end := harvest.EndDate

	maints, _ := repos.Maintenance().ListByOrganization(harvest.OrganizationID)
	for _, m := range maints {
		if inDateRange(m.Date, start, end) {
			add(m.CostCenterID, m.Cost)
		}
	}

	shifts, _ := repos.WorkShift().ListByOrganization(harvest.OrganizationID)
	for _, ws := range shifts {
		if inDateRange(ws.Date, start, end) {
			add(ws.CostCenterID, ws.Cost)
		}
	}

	transactions, _ := repos.Financial().ListByOrganization(harvest.OrganizationID)
	for _, ft := range transactions {
		if ft.Type == entity.TranDespesa && inDateRange(ft.Date, start, end) {
			add(ft.CostCenterID, ft.Amount)
		}
	}

	allocs, _ := repos.CostAllocation().ListByHarvest(harvest.ID)
	for _, a := range allocs {
		id := a.CostCenterID
		add(&id, a.TotalAmount)
	}

	return byGroup, nil
}

func inDateRange(d time.Time, start, end *time.Time) bool {
	if start != nil && d.Before(*start) {
		return false
	}
	if end != nil && d.After(*end) {
		return false
	}
	return true
}

func (s *HarvestService) buildIndicators(harvest *entity.Harvest, totalProduction, totalPlantedArea, totalCost float64, costByGroup map[entity.CostGroup]float64) []*entity.Indicator {
	now := time.Now()
	add := func(indicators []*entity.Indicator, t entity.IndicatorType, value float64) []*entity.Indicator {
		return append(indicators, &entity.Indicator{
			ID:             uuid.New().String(),
			OrganizationID: harvest.OrganizationID,
			HarvestID:      harvest.ID,
			Type:           t,
			Value:          value,
			CalculatedAt:   now,
		})
	}

	var indicators []*entity.Indicator
	indicators = add(indicators, entity.IndProducaoTotal, totalProduction)
	indicators = add(indicators, entity.IndCustoTotal, totalCost)

	if totalPlantedArea > 0 {
		indicators = add(indicators, entity.IndAreaProducao, totalPlantedArea)
		indicators = add(indicators, entity.IndSacasHA, totalProduction/totalPlantedArea)
	}

	if totalProduction > 0 {
		indicators = add(indicators, entity.IndCustoSaca, totalCost/totalProduction)
	}

	coe := costByGroup[entity.CostGroupOperacionalEfetivo]
	cot := coe + costByGroup[entity.CostGroupMaoDeObraFamiliar] + costByGroup[entity.CostGroupCapitalDepreciacao]
	ct := cot + costByGroup[entity.CostGroupRemuneracaoCapital]

	indicators = add(indicators, entity.IndCOE, coe)
	indicators = add(indicators, entity.IndCOT, cot)
	indicators = add(indicators, entity.IndCTProducao, ct)

	if totalPlantedArea > 0 {
		indicators = add(indicators, entity.IndCOEPorArea, coe/totalPlantedArea)
		indicators = add(indicators, entity.IndCOTPorArea, cot/totalPlantedArea)
		indicators = add(indicators, entity.IndCTProducaoPorArea, ct/totalPlantedArea)
	}

	if totalProduction > 0 {
		indicators = add(indicators, entity.IndCOEPorSaca, coe/totalProduction)
		indicators = add(indicators, entity.IndCOTPorSaca, cot/totalProduction)
		indicators = add(indicators, entity.IndCTProducaoPorSaca, ct/totalProduction)
	}

	return indicators
}
