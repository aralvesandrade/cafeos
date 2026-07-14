package repository

import (
	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	"gorm.io/gorm"
)

type AgriculturalProductRepository struct {
	db *gorm.DB
}

func NewAgriculturalProductRepository(db *gorm.DB) *AgriculturalProductRepository {
	return &AgriculturalProductRepository{db: db}
}

func (r *AgriculturalProductRepository) WithTx(tx *gorm.DB) *AgriculturalProductRepository {
	return &AgriculturalProductRepository{db: tx}
}

func (r *AgriculturalProductRepository) Create(p *entity.AgriculturalProduct) error {
	return r.db.Create(p).Error
}

func (r *AgriculturalProductRepository) GetByID(id string) (*entity.AgriculturalProduct, error) {
	var p entity.AgriculturalProduct
	err := r.db.First(&p, "id = ?", id).Error
	return &p, err
}

func (r *AgriculturalProductRepository) ListByOrganization(organizationID string) ([]*entity.AgriculturalProduct, error) {
	var products []*entity.AgriculturalProduct
	err := r.db.Where("organization_id = ?", organizationID).Order("name").Find(&products).Error
	return products, err
}

func (r *AgriculturalProductRepository) Update(p *entity.AgriculturalProduct) error {
	return r.db.Save(p).Error
}

func (r *AgriculturalProductRepository) Delete(id string) error {
	return r.db.Delete(&entity.AgriculturalProduct{}, "id = ?", id).Error
}
