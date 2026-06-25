package repository

import "github.com/aralvesandrade/cafeos/internal/domain/entity"

type HarvestRepository interface {
	Create(harvest *entity.Harvest) error
	GetByID(id string) (*entity.Harvest, error)
	ListByTenant(tenantID string) ([]*entity.Harvest, error)
	Update(harvest *entity.Harvest) error
}

type HarvestProductionRepository interface {
	Create(prod *entity.HarvestProduction) error
	GetByID(id string) (*entity.HarvestProduction, error)
	ListByHarvest(harvestID string) ([]*entity.HarvestProduction, error)
	ListByPlot(plotID string) ([]*entity.HarvestProduction, error)
	Update(prod *entity.HarvestProduction) error
}
