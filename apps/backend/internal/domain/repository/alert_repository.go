package repository

import "github.com/aralvesandrade/cafeos/internal/domain/entity"

type AlertRepository interface {
	Create(a *entity.Alert) error
	GetByID(id string) (*entity.Alert, error)
	ListByOrganization(organizationID string) ([]*entity.Alert, error)
	Update(a *entity.Alert) error
}
