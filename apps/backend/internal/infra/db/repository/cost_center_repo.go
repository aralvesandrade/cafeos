package repository

import (
	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	"gorm.io/gorm"
)

type CostCenterRepository struct {
	db *gorm.DB
}

func NewCostCenterRepository(db *gorm.DB) *CostCenterRepository {
	return &CostCenterRepository{db: db}
}

func (r *CostCenterRepository) WithTx(tx *gorm.DB) *CostCenterRepository {
	return &CostCenterRepository{db: tx}
}

func (r *CostCenterRepository) Create(cc *entity.CostCenter) error {
	return r.db.Create(cc).Error
}

func (r *CostCenterRepository) GetByID(id string) (*entity.CostCenter, error) {
	var cc entity.CostCenter
	err := r.db.First(&cc, "id = ?", id).Error
	return &cc, err
}

func (r *CostCenterRepository) ListByOrganization(organizationID string) ([]*entity.CostCenter, error) {
	var ccs []*entity.CostCenter
	err := r.db.Where("organization_id = ?", organizationID).Order("code").Find(&ccs).Error
	return ccs, err
}

func (r *CostCenterRepository) Update(cc *entity.CostCenter) error {
	return r.db.Save(cc).Error
}

func (r *CostCenterRepository) Delete(id string) error {
	return r.db.Delete(&entity.CostCenter{}, "id = ?", id).Error
}
