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
	ruleEngine      *RuleEngine
	alertRepo       repository.AlertRepository
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
	ruleEngine *RuleEngine,
	alertRepo repository.AlertRepository,
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
		ruleEngine:      ruleEngine,
		alertRepo:       alertRepo,
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

	var computedIndicators []*entity.Indicator
	err = s.transactor.RunInTx(func(repos repository.TransactionProvider) error {
		if err := repos.Harvest().Update(harvest); err != nil {
			return err
		}
		indicators, err := s.calculateIndicatorsWithRepos(harvest, repos)
		computedIndicators = indicators
		return err
	})
	if err != nil {
		return err
	}

	s.eventBus.Publish(event.HarvestFinalized{
		HarvestID:      harvestID,
		OrganizationID: harvest.OrganizationID,
		Year:           harvest.Year,
	})

	s.evaluateRules(harvest, computedIndicators)

	return nil
}

// evaluateRules runs the Rule Engine over the harvest's freshly computed
// indicators and persists+publishes an Alert for every triggered rule.
// Best-effort: rule evaluation failures never fail Finalize.
func (s *HarvestService) evaluateRules(harvest *entity.Harvest, indicators []*entity.Indicator) {
	values := map[string]float64{}
	for _, ind := range indicators {
		switch ind.Type {
		case entity.IndSacasHA, entity.IndCustoSaca:
			values[string(ind.Type)] = ind.Value
		}
	}
	if len(values) == 0 {
		return
	}

	for _, result := range s.ruleEngine.EvaluateIndicators(values) {
		if !result.Triggered || result.Alert == nil {
			continue
		}
		alert := &entity.Alert{
			ID:             uuid.New().String(),
			OrganizationID: harvest.OrganizationID,
			HarvestID:      harvest.ID,
			RuleID:         result.Alert.RuleID,
			Message:        result.Alert.Message,
			Severity:       result.Alert.Severity,
			Status:         "aberto",
			CreatedAt:      result.Alert.CreatedAt,
		}
		if err := s.alertRepo.Create(alert); err != nil {
			continue
		}
		s.eventBus.Publish(event.AlertGenerated{
			OrganizationID: harvest.OrganizationID,
			RuleID:         alert.RuleID,
			Message:        alert.Message,
			Severity:       alert.Severity,
			CreatedAt:      alert.CreatedAt,
		})
	}
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

func (s *HarvestService) calculateIndicatorsWithRepos(harvest *entity.Harvest, repos repository.TransactionProvider) ([]*entity.Indicator, error) {
	productions, err := repos.HarvestProduction().ListByHarvest(harvest.ID)
	if err != nil {
		return nil, err
	}

	var totalProduction float64
	for _, p := range productions {
		totalProduction += p.Quantity
	}

	plots, err := repos.Plot().ListByOrganization(harvest.OrganizationID)
	if err != nil {
		return nil, err
	}

	var totalPlantedArea float64
	for _, p := range plots {
		totalPlantedArea += p.AreaHA
	}

	totalCost := s.calculateTotalCostWithRepos(harvest, repos)
	costByGroup, err := s.calculateCostByGroupWithRepos(harvest, repos)
	if err != nil {
		return nil, err
	}

	indicators := s.buildIndicators(harvest, totalProduction, totalPlantedArea, totalCost, costByGroup)

	for _, ind := range indicators {
		if err := repos.Indicator().Create(ind); err != nil {
			return nil, err
		}
	}

	return indicators, nil
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
