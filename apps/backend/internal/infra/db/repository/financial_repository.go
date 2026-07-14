package repository

import (
	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	"gorm.io/gorm"
)

type FinancialRepository struct {
	db *gorm.DB
}

func NewFinancialRepository(db *gorm.DB) *FinancialRepository {
	return &FinancialRepository{db: db}
}

func (r *FinancialRepository) WithTx(tx *gorm.DB) *FinancialRepository {
	return &FinancialRepository{db: tx}
}

func (r *FinancialRepository) Create(t *entity.FinancialTransaction) error {
	return r.db.Create(t).Error
}

func (r *FinancialRepository) GetByID(id string) (*entity.FinancialTransaction, error) {
	var t entity.FinancialTransaction
	err := r.db.First(&t, "id = ?", id).Error
	return &t, err
}

func (r *FinancialRepository) ListByOrganization(organizationID string) ([]*entity.FinancialTransaction, error) {
	var items []*entity.FinancialTransaction
	err := r.db.Where("organization_id = ?", organizationID).Order("date DESC").Find(&items).Error
	return items, err
}

func (r *FinancialRepository) Update(t *entity.FinancialTransaction) error {
	return r.db.Save(t).Error
}

func (r *FinancialRepository) Delete(id string) error {
	return r.db.Delete(&entity.FinancialTransaction{}, "id = ?", id).Error
}
