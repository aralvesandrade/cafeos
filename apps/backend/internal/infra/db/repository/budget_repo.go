package repository

import (
	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	"gorm.io/gorm"
)

type BudgetRepository struct {
	db *gorm.DB
}

func NewBudgetRepository(db *gorm.DB) *BudgetRepository {
	return &BudgetRepository{db: db}
}

func (r *BudgetRepository) WithTx(tx *gorm.DB) *BudgetRepository {
	return &BudgetRepository{db: tx}
}

func (r *BudgetRepository) Create(b *entity.Budget) error {
	return r.db.Create(b).Error
}

func (r *BudgetRepository) GetByID(id string) (*entity.Budget, error) {
	var b entity.Budget
	err := r.db.Preload("CostCenter").First(&b, "id = ?", id).Error
	return &b, err
}

func (r *BudgetRepository) ListByHarvest(harvestID string) ([]*entity.Budget, error) {
	var budgets []*entity.Budget
	err := r.db.Where("harvest_id = ?", harvestID).Preload("CostCenter").Order("created_at").Find(&budgets).Error
	return budgets, err
}

func (r *BudgetRepository) Update(b *entity.Budget) error {
	return r.db.Save(b).Error
}

func (r *BudgetRepository) Delete(id string) error {
	return r.db.Delete(&entity.Budget{}, "id = ?", id).Error
}
