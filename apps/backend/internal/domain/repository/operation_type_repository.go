package repository

import "github.com/aralvesandrade/cafeos/internal/domain/entity"

type OperationTypeRepository interface {
	Create(ot *entity.OperationType) error
	GetByID(id string) (*entity.OperationType, error)
	GetByOrganizationAndCode(organizationID, code string) (*entity.OperationType, error)
	ListByOrganization(organizationID string) ([]*entity.OperationType, error)
	Update(ot *entity.OperationType) error
	Delete(id string) error
}
