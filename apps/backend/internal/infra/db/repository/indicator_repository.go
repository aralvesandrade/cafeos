package repository

import (
	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	"gorm.io/gorm"
)

type IndicatorRepository struct {
	db *gorm.DB
}

func NewIndicatorRepository(db *gorm.DB) *IndicatorRepository {
	return &IndicatorRepository{db: db}
}

func (r *IndicatorRepository) WithTx(tx *gorm.DB) *IndicatorRepository {
	return &IndicatorRepository{db: tx}
}

func (r *IndicatorRepository) Create(ind *entity.Indicator) error {
	return r.db.Create(ind).Error
}

func (r *IndicatorRepository) ListByHarvest(harvestID string) ([]*entity.Indicator, error) {
	var indicators []*entity.Indicator
	err := r.db.Where("harvest_id = ?", harvestID).Order("type").Find(&indicators).Error
	return indicators, err
}

func (r *IndicatorRepository) ListByOrganization(organizationID string) ([]*entity.Indicator, error) {
	var indicators []*entity.Indicator
	err := r.db.Where("organization_id = ?", organizationID).Order("calculated_at DESC").Find(&indicators).Error
	return indicators, err
}

func (r *IndicatorRepository) ListByOrganizationAndType(organizationID string, indicatorType entity.IndicatorType) ([]*entity.Indicator, error) {
	var indicators []*entity.Indicator
	err := r.db.Where("organization_id = ? AND type = ?", organizationID, indicatorType).
		Order("calculated_at DESC").Find(&indicators).Error
	return indicators, err
}
