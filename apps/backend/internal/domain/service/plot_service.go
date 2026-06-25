package service

import (
	"errors"
	"time"

	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	"github.com/aralvesandrade/cafeos/internal/domain/repository"
	"github.com/google/uuid"
)

type PlotService struct {
	repo repository.PlotRepository
}

func NewPlotService(repo repository.PlotRepository) *PlotService {
	return &PlotService{repo: repo}
}

func (s *PlotService) Create(tenantID, farmID, name, cultivar, soilType string, areaHA float64, plantingYear, altitude int) (*entity.Plot, error) {
	if name == "" {
		return nil, errors.New("plot name is required")
	}
	if areaHA <= 0 {
		return nil, errors.New("area must be greater than zero")
	}

	plot := &entity.Plot{
		ID:           uuid.New().String(),
		TenantID:     tenantID,
		FarmID:       farmID,
		Name:         name,
		AreaHA:       areaHA,
		Cultivar:     cultivar,
		PlantingYear: plantingYear,
		Altitude:     altitude,
		SoilType:     soilType,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.repo.Create(plot); err != nil {
		return nil, err
	}
	return plot, nil
}

func (s *PlotService) GetByID(id string) (*entity.Plot, error) {
	return s.repo.GetByID(id)
}

func (s *PlotService) ListByFarm(farmID string) ([]*entity.Plot, error) {
	return s.repo.ListByFarm(farmID)
}

func (s *PlotService) ListByTenant(tenantID string) ([]*entity.Plot, error) {
	return s.repo.ListByTenant(tenantID)
}

func (s *PlotService) Update(plot *entity.Plot) error {
	if plot.Name == "" {
		return errors.New("plot name is required")
	}
	return s.repo.Update(plot)
}

func (s *PlotService) Delete(id string) error {
	return s.repo.Delete(id)
}
