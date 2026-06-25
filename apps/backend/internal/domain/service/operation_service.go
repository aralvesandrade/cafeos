package service

import (
	"errors"
	"time"

	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	"github.com/aralvesandrade/cafeos/internal/domain/repository"
	"github.com/aralvesandrade/cafeos/internal/event"
	"github.com/google/uuid"
)

type OperationService struct {
	repo     repository.OperationRepository
	eventBus event.Bus
}

func NewOperationService(repo repository.OperationRepository, eventBus event.Bus) *OperationService {
	return &OperationService{repo: repo, eventBus: eventBus}
}

func (s *OperationService) Create(tenantID, plotID string, opType entity.OperationType, date time.Time, responsible, productUsed string, quantity, cost float64, notes string) (*entity.Operation, error) {
	if plotID == "" {
		return nil, errors.New("plot is required")
	}
	if opType == "" {
		return nil, errors.New("operation type is required")
	}
	if cost < 0 {
		return nil, errors.New("cost cannot be negative")
	}

	op := &entity.Operation{
		ID:          uuid.New().String(),
		TenantID:    tenantID,
		PlotID:      plotID,
		Type:        opType,
		Date:        date,
		Responsible: responsible,
		ProductUsed: productUsed,
		Quantity:    quantity,
		Cost:        cost,
		Notes:       notes,
		CreatedAt:   time.Now(),
	}

	if err := s.repo.Create(op); err != nil {
		return nil, err
	}

	s.eventBus.Publish(event.OperationRegistered{
		OperationID: op.ID,
		TenantID:    tenantID,
		PlotID:      plotID,
		Type:        string(op.Type),
		Cost:        cost,
		Date:        date,
	})

	return op, nil
}

func (s *OperationService) GetByID(id string) (*entity.Operation, error) {
	return s.repo.GetByID(id)
}

func (s *OperationService) ListByPlot(plotID string) ([]*entity.Operation, error) {
	return s.repo.ListByPlot(plotID)
}

func (s *OperationService) ListByTenant(tenantID string) ([]*entity.Operation, error) {
	return s.repo.ListByTenant(tenantID)
}

func (s *OperationService) ListRecent(tenantID string, limit int) ([]*entity.Operation, error) {
	all, err := s.repo.ListByTenant(tenantID)
	if err != nil {
		return nil, err
	}
	if len(all) > limit {
		all = all[:limit]
	}
	return all, nil
}

func (s *OperationService) Delete(id string) error {
	return s.repo.Delete(id)
}
