package repository

import "github.com/aralvesandrade/cafeos/internal/domain/entity"

type AgriculturalProductRepository interface {
	Create(product *entity.AgriculturalProduct) error
	GetByID(id string) (*entity.AgriculturalProduct, error)
	ListByOrganization(organizationID string) ([]*entity.AgriculturalProduct, error)
	Update(product *entity.AgriculturalProduct) error
	Delete(id string) error
}
