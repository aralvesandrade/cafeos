package service

import (
	"errors"
	"time"

	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	"github.com/aralvesandrade/cafeos/internal/domain/repository"
	"github.com/google/uuid"
)

type BudgetService struct {
	repo        repository.BudgetRepository
	harvestRepo repository.HarvestRepository
	opRepo      repository.OperationRepository
	maintRepo   repository.MaintenanceRepository
	shiftRepo   repository.WorkShiftRepository
	finRepo     repository.FinancialRepository
	allocRepo   repository.CostAllocationRepository
}

func NewBudgetService(
	repo repository.BudgetRepository,
	harvestRepo repository.HarvestRepository,
	opRepo repository.OperationRepository,
	maintRepo repository.MaintenanceRepository,
	shiftRepo repository.WorkShiftRepository,
	finRepo repository.FinancialRepository,
	allocRepo repository.CostAllocationRepository,
) *BudgetService {
	return &BudgetService{
		repo:        repo,
		harvestRepo: harvestRepo,
		opRepo:      opRepo,
		maintRepo:   maintRepo,
		shiftRepo:   shiftRepo,
		finRepo:     finRepo,
		allocRepo:   allocRepo,
	}
}

func (s *BudgetService) Create(organizationID, harvestID, costCenterID string, plannedAmount float64, description string) (*entity.Budget, error) {
	if plannedAmount < 0 {
		return nil, errors.New("planned amount cannot be negative")
	}
	b := &entity.Budget{
		ID:             uuid.New().String(),
		OrganizationID: organizationID,
		HarvestID:      harvestID,
		CostCenterID:   costCenterID,
		PlannedAmount:  plannedAmount,
		Description:    description,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	if err := s.repo.Create(b); err != nil {
		return nil, err
	}
	return b, nil
}

func (s *BudgetService) GetByID(id string) (*entity.Budget, error) {
	return s.repo.GetByID(id)
}

// calculateRealizedByCostCenter sums Operation, Maintenance, WorkShift,
// FinancialTransaction (despesa) and CostAllocation costs for the harvest,
// bucketed by CostCenterID — every cost center counts, not just ones with a
// SENAR CostGroup classification (unlike calculateCostByGroupWithRepos in
// harvest_service.go, which this mirrors).
func (s *BudgetService) calculateRealizedByCostCenter(harvest *entity.Harvest) (map[string]float64, error) {
	byCostCenter := map[string]float64{}
	add := func(costCenterID *string, amount float64) {
		if costCenterID == nil {
			return
		}
		byCostCenter[*costCenterID] += amount
	}

	ops, err := s.opRepo.ListByOrganization(harvest.OrganizationID)
	if err != nil {
		return nil, err
	}
	for _, op := range ops {
		if op.HarvestID != nil && *op.HarvestID == harvest.ID {
			add(op.CostCenterID, op.Cost)
		}
	}

	start := harvest.StartDate
	end := harvest.EndDate

	maints, err := s.maintRepo.ListByOrganization(harvest.OrganizationID)
	if err != nil {
		return nil, err
	}
	for _, m := range maints {
		if inDateRange(m.Date, start, end) {
			add(m.CostCenterID, m.Cost)
		}
	}

	shifts, err := s.shiftRepo.ListByOrganization(harvest.OrganizationID)
	if err != nil {
		return nil, err
	}
	for _, ws := range shifts {
		if inDateRange(ws.Date, start, end) {
			add(ws.CostCenterID, ws.Cost)
		}
	}

	transactions, err := s.finRepo.ListByOrganization(harvest.OrganizationID)
	if err != nil {
		return nil, err
	}
	for _, ft := range transactions {
		if ft.Type == entity.TranDespesa && inDateRange(ft.Date, start, end) {
			add(ft.CostCenterID, ft.Amount)
		}
	}

	allocs, err := s.allocRepo.ListByHarvest(harvest.ID)
	if err != nil {
		return nil, err
	}
	for _, a := range allocs {
		id := a.CostCenterID
		add(&id, a.TotalAmount)
	}

	return byCostCenter, nil
}

func (s *BudgetService) ListByHarvest(harvestID string) ([]*entity.Budget, error) {
	budgets, err := s.repo.ListByHarvest(harvestID)
	if err != nil {
		return nil, err
	}
	if len(budgets) == 0 {
		return budgets, nil
	}

	harvest, err := s.harvestRepo.GetByID(harvestID)
	if err != nil {
		return nil, err
	}

	realized, err := s.calculateRealizedByCostCenter(harvest)
	if err != nil {
		return nil, err
	}

	for _, b := range budgets {
		b.CostCenterName = b.CostCenter.Name
		b.RealizedAmount = realized[b.CostCenterID]
		b.Variance = b.PlannedAmount - b.RealizedAmount
	}

	return budgets, nil
}

func (s *BudgetService) Update(b *entity.Budget) error {
	b.UpdatedAt = time.Now()
	return s.repo.Update(b)
}

func (s *BudgetService) Delete(id string) error {
	return s.repo.Delete(id)
}
