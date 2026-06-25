package repository

import (
	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	"gorm.io/gorm"
)

type TenantRepository struct {
	db *gorm.DB
}

func NewTenantRepository(db *gorm.DB) *TenantRepository {
	return &TenantRepository{db: db}
}

func (r *TenantRepository) WithTx(tx *gorm.DB) *TenantRepository {
	return &TenantRepository{db: tx}
}

func (r *TenantRepository) Create(t *entity.Tenant) error {
	return r.db.Create(t).Error
}

func (r *TenantRepository) GetByID(id string) (*entity.Tenant, error) {
	var t entity.Tenant
	err := r.db.First(&t, "id = ?", id).Error
	return &t, err
}

func (r *TenantRepository) GetBySlug(slug string) (*entity.Tenant, error) {
	var t entity.Tenant
	err := r.db.Where("slug = ?", slug).First(&t).Error
	return &t, err
}

func (r *TenantRepository) List() ([]*entity.Tenant, error) {
	var tenants []*entity.Tenant
	err := r.db.Order("name").Find(&tenants).Error
	return tenants, err
}

func (r *TenantRepository) Update(t *entity.Tenant) error {
	return r.db.Save(t).Error
}

func (r *TenantRepository) Delete(id string) error {
	return r.db.Delete(&entity.Tenant{}, "id = ?", id).Error
}
