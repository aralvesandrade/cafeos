package repository

import "github.com/aralvesandrade/cafeos/internal/domain/entity"

type OperationRepository interface {
	Create(op *entity.Operation) error
	GetByID(id string) (*entity.Operation, error)
	ListByPlot(plotID string) ([]*entity.Operation, error)
	ListByTenant(tenantID string) ([]*entity.Operation, error)
	ListByTenantAndPeriod(tenantID string, start, end string) ([]*entity.Operation, error)
	Delete(id string) error
}
