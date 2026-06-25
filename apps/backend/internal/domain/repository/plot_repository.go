package repository

import "github.com/aralvesandrade/cafeos/internal/domain/entity"

type PlotRepository interface {
	Create(plot *entity.Plot) error
	GetByID(id string) (*entity.Plot, error)
	ListByFarm(farmID string) ([]*entity.Plot, error)
	ListByTenant(tenantID string) ([]*entity.Plot, error)
	Update(plot *entity.Plot) error
	Delete(id string) error
}
