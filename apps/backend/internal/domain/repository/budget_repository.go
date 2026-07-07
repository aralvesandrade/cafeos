package repository

import "github.com/aralvesandrade/cafeos/internal/domain/entity"

type BudgetRepository interface {
	Create(b *entity.Budget) error
	GetByID(id string) (*entity.Budget, error)
	ListByHarvest(harvestID string) ([]*entity.Budget, error)
	Update(b *entity.Budget) error
	Delete(id string) error
}
