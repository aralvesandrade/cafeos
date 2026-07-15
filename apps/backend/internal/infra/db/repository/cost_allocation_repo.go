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

func populateAllocationNames(a *entity.CostAllocation) {
	a.CostCenterName = a.CostCenter.Name
	for i := range a.Items {
		a.Items[i].PlotName = a.Items[i].Plot.Name
	}
}

func (r *CostAllocationRepository) GetByID(id string) (*entity.CostAllocation, error) {
	var a entity.CostAllocation
	err := r.db.Preload("Items.Plot").Preload("CostCenter").First(&a, "id = ?", id).Error
	if err == nil {
		populateAllocationNames(&a)
	}
	return &a, err
}

func (r *CostAllocationRepository) ListByHarvest(harvestID string) ([]*entity.CostAllocation, error) {
	var allocs []*entity.CostAllocation
	err := r.db.Where("harvest_id = ?", harvestID).Preload("Items.Plot").Preload("CostCenter").Order("date").Find(&allocs).Error
	for _, a := range allocs {
		populateAllocationNames(a)
	}
	return allocs, err
}

func (r *CostAllocationRepository) Delete(id string) error {
	if err := r.db.Where("allocation_id = ?", id).Delete(&entity.CostAllocationItem{}).Error; err != nil {
		return err
	}
	return r.db.Delete(&entity.CostAllocation{}, "id = ?", id).Error
}
