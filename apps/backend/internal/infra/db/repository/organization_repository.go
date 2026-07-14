package repository

import (
	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	"gorm.io/gorm"
)

type OrganizationRepository struct {
	db *gorm.DB
}

func NewOrganizationRepository(db *gorm.DB) *OrganizationRepository {
	return &OrganizationRepository{db: db}
}

func (r *OrganizationRepository) WithTx(tx *gorm.DB) *OrganizationRepository {
	return &OrganizationRepository{db: tx}
}

func (r *OrganizationRepository) Create(t *entity.Organization) error {
	return r.db.Create(t).Error
}

func (r *OrganizationRepository) GetByID(id string) (*entity.Organization, error) {
	var t entity.Organization
	err := r.db.First(&t, "id = ?", id).Error
	return &t, err
}

func (r *OrganizationRepository) GetBySlug(slug string) (*entity.Organization, error) {
	var t entity.Organization
	err := r.db.Where("slug = ?", slug).First(&t).Error
	return &t, err
}

func (r *OrganizationRepository) List() ([]*entity.Organization, error) {
	var organizations []*entity.Organization
	err := r.db.Order("name").Find(&organizations).Error
	return organizations, err
}

func (r *OrganizationRepository) Update(t *entity.Organization) error {
	return r.db.Save(t).Error
}

func (r *OrganizationRepository) Delete(id string) error {
	return r.db.Delete(&entity.Organization{}, "id = ?", id).Error
}
