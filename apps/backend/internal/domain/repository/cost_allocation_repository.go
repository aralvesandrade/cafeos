package repository

import "github.com/aralvesandrade/cafeos/internal/domain/entity"

type CostAllocationRepository interface {
	Create(a *entity.CostAllocation) error
	GetByID(id string) (*entity.CostAllocation, error)
	ListByHarvest(harvestID string) ([]*entity.CostAllocation, error)
	Delete(id string) error
}
