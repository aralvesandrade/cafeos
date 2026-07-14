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

func validatePlot(plot *entity.Plot) error {
	if plot.Name == "" {
		return errors.New("plot name is required")
	}
	if plot.AreaHA <= 0 {
		return errors.New("area must be greater than zero")
	}
	switch plot.Stage {
	case "", entity.PlotStageFormacao, entity.PlotStageProducao:
	default:
		return errors.New("invalid stage: must be 'formacao' or 'producao'")
	}
	if plot.ActivationDate != nil && plot.DeactivationDate != nil && plot.DeactivationDate.Before(*plot.ActivationDate) {
		return errors.New("deactivation date cannot be before activation date")
	}
	return nil
}

func (s *PlotService) Create(plot *entity.Plot) (*entity.Plot, error) {
	if plot.Stage == "" {
		plot.Stage = entity.PlotStageFormacao
	}
	if err := validatePlot(plot); err != nil {
		return nil, err
	}

	plot.ID = uuid.New().String()
	plot.CreatedAt = time.Now()
	plot.UpdatedAt = time.Now()

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

func (s *PlotService) ListByOrganization(organizationID string) ([]*entity.Plot, error) {
	return s.repo.ListByOrganization(organizationID)
}

func (s *PlotService) Update(plot *entity.Plot) error {
	if err := validatePlot(plot); err != nil {
		return err
	}
	return s.repo.Update(plot)
}

func (s *PlotService) Delete(id string) error {
	return s.repo.Delete(id)
}
