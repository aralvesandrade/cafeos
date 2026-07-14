package service

import (
	"errors"
	"time"

	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	"github.com/aralvesandrade/cafeos/internal/domain/repository"
	"github.com/google/uuid"
)

type CostAllocationService struct {
	repo     repository.CostAllocationRepository
	plotRepo repository.PlotRepository
}

func NewCostAllocationService(repo repository.CostAllocationRepository, plotRepo repository.PlotRepository) *CostAllocationService {
	return &CostAllocationService{repo: repo, plotRepo: plotRepo}
}

type CreateAllocationInput struct {
	OrganizationID string
	HarvestID      string
	CostCenterID   string
	Description    string
	TotalAmount    float64
	Method         entity.AllocationMethod
	Date           time.Time
	Percentages    map[string]float64 // plotID -> percentage (for custom_percentage)
}

func (s *CostAllocationService) Create(input CreateAllocationInput) (*entity.CostAllocation, error) {
	if input.Description == "" {
		return nil, errors.New("description is required")
	}
	if input.TotalAmount <= 0 {
		return nil, errors.New("total amount must be greater than zero")
	}

	allocation := &entity.CostAllocation{
		ID:             uuid.New().String(),
		OrganizationID: input.OrganizationID,
		HarvestID:      input.HarvestID,
		CostCenterID:   input.CostCenterID,
		Description:    input.Description,
		TotalAmount:    input.TotalAmount,
		Method:         input.Method,
		Date:           input.Date,
		CreatedAt:      time.Now(),
	}

	switch input.Method {
	case entity.AllocAreaProportional:
		plots, err := s.plotRepo.ListByOrganization(input.OrganizationID)
		if err != nil {
			return nil, err
		}
		var totalArea float64
		for _, p := range plots {
			totalArea += p.AreaHA
		}
		if totalArea <= 0 {
			return nil, errors.New("no plots with area found for area-proportional allocation")
		}
		for _, p := range plots {
			percentage := (p.AreaHA / totalArea) * 100
			allocation.Items = append(allocation.Items, entity.CostAllocationItem{
				ID:         uuid.New().String(),
				PlotID:     p.ID,
				Amount:     input.TotalAmount * (percentage / 100),
				Percentage: percentage,
			})
		}

	case entity.AllocCustomPercentage:
		if len(input.Percentages) == 0 {
			return nil, errors.New("percentages required for custom_percentage method")
		}
		var sum float64
		for _, v := range input.Percentages {
			sum += v
		}
		if sum < 99.9 || sum > 100.1 {
			return nil, errors.New("percentages must sum to 100")
		}
		for plotID, pct := range input.Percentages {
			allocation.Items = append(allocation.Items, entity.CostAllocationItem{
				ID:         uuid.New().String(),
				PlotID:     plotID,
				Amount:     input.TotalAmount * (pct / 100),
				Percentage: pct,
			})
		}

	default:
		return nil, errors.New("invalid allocation method")
	}

	if err := s.repo.Create(allocation); err != nil {
		return nil, err
	}
	return allocation, nil
}

func (s *CostAllocationService) GetByID(id string) (*entity.CostAllocation, error) {
	return s.repo.GetByID(id)
}

func (s *CostAllocationService) ListByHarvest(harvestID string) ([]*entity.CostAllocation, error) {
	return s.repo.ListByHarvest(harvestID)
}

func (s *CostAllocationService) Delete(id string) error {
	return s.repo.Delete(id)
}
