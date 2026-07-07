package service

import (
	"errors"
	"time"

	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	"github.com/aralvesandrade/cafeos/internal/domain/repository"
	"github.com/google/uuid"
)

type BudgetService struct {
	repo repository.BudgetRepository
}

func NewBudgetService(repo repository.BudgetRepository) *BudgetService {
	return &BudgetService{repo: repo}
}

func (s *BudgetService) Create(tenantID, harvestID, costCenterID string, plannedAmount float64, description string) (*entity.Budget, error) {
	if plannedAmount < 0 {
		return nil, errors.New("planned amount cannot be negative")
	}
	b := &entity.Budget{
		ID:            uuid.New().String(),
		TenantID:      tenantID,
		HarvestID:     harvestID,
		CostCenterID:  costCenterID,
		PlannedAmount: plannedAmount,
		Description:   description,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	if err := s.repo.Create(b); err != nil {
		return nil, err
	}
	return b, nil
}

func (s *BudgetService) GetByID(id string) (*entity.Budget, error) {
	return s.repo.GetByID(id)
}

func (s *BudgetService) ListByHarvest(harvestID string) ([]*entity.Budget, error) {
	return s.repo.ListByHarvest(harvestID)
}

func (s *BudgetService) Update(b *entity.Budget) error {
	b.UpdatedAt = time.Now()
	return s.repo.Update(b)
}

func (s *BudgetService) Delete(id string) error {
	return s.repo.Delete(id)
}
