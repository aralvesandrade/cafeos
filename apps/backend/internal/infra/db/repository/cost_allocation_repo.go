package repository

import (
	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	"gorm.io/gorm"
)

type CostAllocationRepository struct {
	db *gorm.DB
}

func NewCostAllocationRepository(db *gorm.DB) *CostAllocationRepository {
	return &CostAllocationRepository{db: db}
}

func (r *CostAllocationRepository) WithTx(tx *gorm.DB) *CostAllocationRepository {
	return &CostAllocationRepository{db: tx}
}

func (r *CostAllocationRepository) Create(a *entity.CostAllocation) error {
	return r.db.Create(a).Error
}

func (r *CostAllocationRepository) GetByID(id string) (*entity.CostAllocation, error) {
	var a entity.CostAllocation
	err := r.db.Preload("Items.Plot").Preload("CostCenter").First(&a, "id = ?", id).Error
	return &a, err
}

func (r *CostAllocationRepository) ListByHarvest(harvestID string) ([]*entity.CostAllocation, error) {
	var allocs []*entity.CostAllocation
	err := r.db.Where("harvest_id = ?", harvestID).Preload("Items.Plot").Preload("CostCenter").Order("date").Find(&allocs).Error
	return allocs, err
}

func (r *CostAllocationRepository) Delete(id string) error {
	return r.db.Select("Items").Delete(&entity.CostAllocation{}, "id = ?", id).Error
}
