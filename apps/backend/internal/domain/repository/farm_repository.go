package repository

import "github.com/aralvesandrade/cafeos/internal/domain/entity"

type FarmRepository interface {
	Create(farm *entity.Farm) error
	GetByID(id string) (*entity.Farm, error)
	ListByTenant(tenantID string) ([]*entity.Farm, error)
	Update(farm *entity.Farm) error
	Delete(id string) error
}
