package service

import (
	"errors"
	"time"

	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	"github.com/aralvesandrade/cafeos/internal/domain/repository"
	"github.com/google/uuid"
)

type FarmService struct {
	repo repository.FarmRepository
}

func NewFarmService(repo repository.FarmRepository) *FarmService {
	return &FarmService{repo: repo}
}

func (s *FarmService) Create(tenantID, name, owner, location string, totalAreaHA, plantedAreaHA float64) (*entity.Farm, error) {
	if name == "" {
		return nil, errors.New("farm name is required")
	}
	if totalAreaHA <= 0 {
		return nil, errors.New("total area must be greater than zero")
	}
	if plantedAreaHA > totalAreaHA {
		return nil, errors.New("planted area cannot exceed total area")
	}

	farm := &entity.Farm{
		ID:           uuid.New().String(),
		TenantID:     tenantID,
		Name:         name,
		Owner:        owner,
		Location:     location,
		TotalAreaHA:  totalAreaHA,
		PlantedAreaHA: plantedAreaHA,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.repo.Create(farm); err != nil {
		return nil, err
	}
	return farm, nil
}

func (s *FarmService) GetByID(id string) (*entity.Farm, error) {
	return s.repo.GetByID(id)
}

func (s *FarmService) ListByTenant(tenantID string) ([]*entity.Farm, error) {
	return s.repo.ListByTenant(tenantID)
}

func (s *FarmService) Update(farm *entity.Farm) error {
	if farm.Name == "" {
		return errors.New("farm name is required")
	}
	if farm.PlantedAreaHA > farm.TotalAreaHA {
		return errors.New("planted area cannot exceed total area")
	}
	return s.repo.Update(farm)
}

func (s *FarmService) Delete(id string) error {
	return s.repo.Delete(id)
}
