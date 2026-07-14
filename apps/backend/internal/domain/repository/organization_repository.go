package repository

import "github.com/aralvesandrade/cafeos/internal/domain/entity"

type OrganizationRepository interface {
	Create(organization *entity.Organization) error
	GetByID(id string) (*entity.Organization, error)
	GetBySlug(slug string) (*entity.Organization, error)
	List() ([]*entity.Organization, error)
	Update(organization *entity.Organization) error
	Delete(id string) error
}
