package service

import (
	"errors"
	"time"

	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	"github.com/aralvesandrade/cafeos/internal/domain/repository"
	"github.com/google/uuid"
)

type OperationTypeService struct {
	repo repository.OperationTypeRepository
}

func NewOperationTypeService(repo repository.OperationTypeRepository) *OperationTypeService {
	return &OperationTypeService{repo: repo}
}

func (s *OperationTypeService) Create(organizationID, name, code, color string) (*entity.OperationType, error) {
	if name == "" {
		return nil, errors.New("operation type name is required")
	}
	if code == "" {
		return nil, errors.New("operation type code is required")
	}
	if color == "" {
		color = "default"
	}

	ot := &entity.OperationType{
		ID:             uuid.New().String(),
		OrganizationID: organizationID,
		Name:           name,
		Code:           code,
		Color:          color,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.repo.Create(ot); err != nil {
		return nil, err
	}
	return ot, nil
}

func (s *OperationTypeService) GetByID(id string) (*entity.OperationType, error) {
	return s.repo.GetByID(id)
}

func (s *OperationTypeService) ListByOrganization(organizationID string) ([]*entity.OperationType, error) {
	return s.repo.ListByOrganization(organizationID)
}

func (s *OperationTypeService) Update(ot *entity.OperationType) error {
	if ot.Name == "" {
		return errors.New("operation type name is required")
	}
	if ot.Code == "" {
		return errors.New("operation type code is required")
	}
	ot.UpdatedAt = time.Now()
	return s.repo.Update(ot)
}

func (s *OperationTypeService) Delete(id string) error {
	return s.repo.Delete(id)
}
