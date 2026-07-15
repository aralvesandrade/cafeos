package repository

import (
	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	"gorm.io/gorm"
)

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) ListForOrganization(organizationID string) ([]*entity.Role, error) {
	var roles []*entity.Role
	query := r.db.Order("is_system DESC, name ASC")
	if organizationID == "" {
		query = query.Where("organization_id IS NULL")
	} else {
		query = query.Where("organization_id IS NULL OR organization_id = ?", organizationID)
	}
	err := query.Find(&roles).Error
	return roles, err
}

func (r *RoleRepository) GetByID(id string) (*entity.Role, error) {
	var role entity.Role
	err := r.db.Where("id = ?", id).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) FindByKey(organizationID, key string) (*entity.Role, error) {
	var role entity.Role
	err := r.db.Where("organization_id = ? AND key = ?", organizationID, key).First(&role).Error
	if err == nil {
		return &role, nil
	}
	if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	err = r.db.Where("organization_id IS NULL AND key = ?", key).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) Create(role *entity.Role) error {
	return r.db.Create(role).Error
}

func (r *RoleRepository) Update(role *entity.Role) error {
	return r.db.Save(role).Error
}

func (r *RoleRepository) Delete(id string) error {
	return r.db.Delete(&entity.Role{}, "id = ?", id).Error
}

func (r *RoleRepository) CountByOrganization(organizationID string) (int64, error) {
	var count int64
	err := r.db.Model(&entity.Role{}).Where("organization_id = ?", organizationID).Count(&count).Error
	return count, err
}
