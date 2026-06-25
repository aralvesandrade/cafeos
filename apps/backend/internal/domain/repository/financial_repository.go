package repository

import "github.com/aralvesandrade/cafeos/internal/domain/entity"

type FinancialRepository interface {
	Create(tx *entity.FinancialTransaction) error
	GetByID(id string) (*entity.FinancialTransaction, error)
	ListByTenant(tenantID string) ([]*entity.FinancialTransaction, error)
	Update(tx *entity.FinancialTransaction) error
	Delete(id string) error
}
