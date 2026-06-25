package repository

import "github.com/aralvesandrade/cafeos/internal/domain/entity"

type TenantRepository interface {
	Create(tenant *entity.Tenant) error
	GetByID(id string) (*entity.Tenant, error)
	GetBySlug(slug string) (*entity.Tenant, error)
	List() ([]*entity.Tenant, error)
	Update(tenant *entity.Tenant) error
	Delete(id string) error
}
