package repository

import "github.com/aralvesandrade/cafeos/internal/domain/entity"

type CostCenterRepository interface {
	Create(cc *entity.CostCenter) error
	GetByID(id string) (*entity.CostCenter, error)
	ListByTenant(tenantID string) ([]*entity.CostCenter, error)
	Update(cc *entity.CostCenter) error
	Delete(id string) error
}
