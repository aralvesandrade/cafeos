package repository

import (
	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	"gorm.io/gorm"
)

type AlertRepository struct {
	db *gorm.DB
}

func NewAlertRepository(db *gorm.DB) *AlertRepository {
	return &AlertRepository{db: db}
}

func (r *AlertRepository) WithTx(tx *gorm.DB) *AlertRepository {
	return &AlertRepository{db: tx}
}

func (r *AlertRepository) Create(a *entity.Alert) error {
	return r.db.Create(a).Error
}

func (r *AlertRepository) GetByID(id string) (*entity.Alert, error) {
	var a entity.Alert
	err := r.db.First(&a, "id = ?", id).Error
	return &a, err
}

func (r *AlertRepository) ListByOrganization(organizationID string) ([]*entity.Alert, error) {
	var alerts []*entity.Alert
	err := r.db.Where("organization_id = ?", organizationID).Order("created_at DESC").Find(&alerts).Error
	return alerts, err
}

func (r *AlertRepository) Update(a *entity.Alert) error {
	return r.db.Save(a).Error
}
