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
	harvestRepo           repository.HarvestRepository
	productionRepo        repository.HarvestProductionRepository
	indicatorRepo         repository.IndicatorRepository
	plotRepo              repository.PlotRepository
	operationRepo         repository.OperationRepository
	eventBus              event.Bus
}

func NewHarvestService(
	harvestRepo repository.HarvestRepository,
	productionRepo repository.HarvestProductionRepository,
	indicatorRepo repository.IndicatorRepository,
	plotRepo repository.PlotRepository,
	operationRepo repository.OperationRepository,
	eventBus event.Bus,
) *HarvestService {
	return &HarvestService{
		harvestRepo:    harvestRepo,
		productionRepo: productionRepo,
		indicatorRepo:  indicatorRepo,
		plotRepo:       plotRepo,
		operationRepo:  operationRepo,
		eventBus:       eventBus,
	}
}

func (s *HarvestService) Create(tenantID string, year int, description string, estimatedProduction float64) (*entity.Harvest, error) {
	if year <= 0 {
		return nil, errors.New("year is required")
	}
	if estimatedProduction < 0 {
		return nil, errors.New("estimated production cannot be negative")
	}

	harvest := &entity.Harvest{
		ID:                 uuid.New().String(),
		TenantID:           tenantID,
		Year:               year,
		Description:        description,
		EstimatedProduction: estimatedProduction,
		Status:             entity.HarvestPlanejada,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	if err := s.harvestRepo.Create(harvest); err != nil {
		return nil, err
	}
	return harvest, nil
}

func (s *HarvestService) GetByID(id string) (*entity.Harvest, error) {
	return s.harvestRepo.GetByID(id)
}

func (s *HarvestService) ListByTenant(tenantID string) ([]*entity.Harvest, error) {
	return s.harvestRepo.ListByTenant(tenantID)
}

func (s *HarvestService) Finalize(harvestID string) error {
	harvest, err := s.harvestRepo.GetByID(harvestID)
	if err != nil {
		return err
	}

	harvest.Status = entity.HarvestFinalizada
	harvest.UpdatedAt = time.Now()

	if err := s.harvestRepo.Update(harvest); err != nil {
		return err
	}

	if err := s.calculateIndicators(harvest); err != nil {
		return err
	}

	s.eventBus.Publish(event.HarvestFinalized{
		HarvestID: harvestID,
		TenantID:  harvest.TenantID,
		Year:      harvest.Year,
	})

	return nil
}

func (s *HarvestService) RecordProduction(tenantID, harvestID, plotID string, quantity float64, notes string) (*entity.HarvestProduction, error) {
	if quantity < 0 {
		return nil, errors.New("quantity cannot be negative")
	}

	prod := &entity.HarvestProduction{
		ID:         uuid.New().String(),
		TenantID:   tenantID,
		HarvestID:  harvestID,
		PlotID:     plotID,
		Quantity:   quantity,
		RecordedAt: time.Now(),
		Notes:      notes,
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

	plots, err := s.plotRepo.ListByTenant(harvest.TenantID)
	if err != nil {
		return err
	}

	var totalPlantedArea float64
	for _, p := range plots {
		totalPlantedArea += p.AreaHA
	}

	operations, err := s.operationRepo.ListByTenant(harvest.TenantID)
	if err != nil {
		return err
	}

	var totalCost float64
	for _, op := range operations {
		totalCost += op.Cost
	}

	indicators := []*entity.Indicator{
		{
			ID:           uuid.New().String(),
			TenantID:     harvest.TenantID,
			HarvestID:    harvest.ID,
			Type:         entity.IndProducaoTotal,
			Value:        totalProduction,
			CalculatedAt: time.Now(),
		},
		{
			ID:           uuid.New().String(),
			TenantID:     harvest.TenantID,
			HarvestID:    harvest.ID,
			Type:         entity.IndCustoTotal,
			Value:        totalCost,
			CalculatedAt: time.Now(),
		},
	}

	if totalPlantedArea > 0 {
		indicators = append(indicators, &entity.Indicator{
			ID:           uuid.New().String(),
			TenantID:     harvest.TenantID,
			HarvestID:    harvest.ID,
			Type:         entity.IndSacasHA,
			Value:        totalProduction / totalPlantedArea,
			CalculatedAt: time.Now(),
		})
	}

	if totalProduction > 0 {
		indicators = append(indicators, &entity.Indicator{
			ID:           uuid.New().String(),
			TenantID:     harvest.TenantID,
			HarvestID:    harvest.ID,
			Type:         entity.IndCustoSaca,
			Value:        totalCost / totalProduction,
			CalculatedAt: time.Now(),
		})
	}

	for _, ind := range indicators {
		if err := s.indicatorRepo.Create(ind); err != nil {
			return err
		}
		s.eventBus.Publish(event.IndicatorUpdated{
			TenantID:    harvest.TenantID,
			HarvestID:   harvest.ID,
			IndicatorID: ind.ID,
			Type:        string(ind.Type),
			Value:       ind.Value,
		})
	}

	return nil
}
