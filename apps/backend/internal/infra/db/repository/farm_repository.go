package repository

import (
	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	"gorm.io/gorm"
)

type FarmRepository struct {
	db *gorm.DB
}

func NewFarmRepository(db *gorm.DB) *FarmRepository {
	return &FarmRepository{db: db}
}

func (r *FarmRepository) WithTx(tx *gorm.DB) *FarmRepository {
	return &FarmRepository{db: tx}
}

func (r *FarmRepository) Create(f *entity.Farm) error {
	return r.db.Create(f).Error
}

func (r *FarmRepository) GetByID(id string) (*entity.Farm, error) {
	var f entity.Farm
	err := r.db.First(&f, "id = ?", id).Error
	return &f, err
}

func (r *FarmRepository) ListByTenant(tenantID string) ([]*entity.Farm, error) {
	var farms []*entity.Farm
	err := r.db.Where("tenant_id = ?", tenantID).Order("name").Find(&farms).Error
	return farms, err
}

func (r *FarmRepository) Update(f *entity.Farm) error {
	return r.db.Save(f).Error
}

func (r *FarmRepository) Delete(id string) error {
	return r.db.Delete(&entity.Farm{}, "id = ?", id).Error
}
