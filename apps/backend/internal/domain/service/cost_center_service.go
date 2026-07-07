package service

import (
	"errors"
	"time"

	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	"github.com/aralvesandrade/cafeos/internal/domain/repository"
	"github.com/google/uuid"
)

type CostCenterService struct {
	repo repository.CostCenterRepository
}

func NewCostCenterService(repo repository.CostCenterRepository) *CostCenterService {
	return &CostCenterService{repo: repo}
}

func (s *CostCenterService) Create(tenantID, name, code string, ccType entity.CostCenterType, description string) (*entity.CostCenter, error) {
	if name == "" {
		return nil, errors.New("cost center name is required")
	}
	if code == "" {
		return nil, errors.New("cost center code is required")
	}

	cc := &entity.CostCenter{
		ID:          uuid.New().String(),
		TenantID:    tenantID,
		Name:        name,
		Code:        code,
		Type:        ccType,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.repo.Create(cc); err != nil {
		return nil, err
	}
	return cc, nil
}

func (s *CostCenterService) GetByID(id string) (*entity.CostCenter, error) {
	return s.repo.GetByID(id)
}

func (s *CostCenterService) ListByTenant(tenantID string) ([]*entity.CostCenter, error) {
	return s.repo.ListByTenant(tenantID)
}

func (s *CostCenterService) Update(cc *entity.CostCenter) error {
	if cc.Name == "" {
		return errors.New("cost center name is required")
	}
	cc.UpdatedAt = time.Now()
	return s.repo.Update(cc)
}

func (s *CostCenterService) Delete(id string) error {
	return s.repo.Delete(id)
}
