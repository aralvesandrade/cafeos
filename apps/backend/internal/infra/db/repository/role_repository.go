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

func (r *RoleRepository) List() ([]*entity.Role, error) {
	var roles []*entity.Role
	err := r.db.Order("is_system DESC, name ASC").Find(&roles).Error
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

func (r *RoleRepository) FindByKey(key string) (*entity.Role, error) {
	var role entity.Role
	err := r.db.Where("key = ?", key).First(&role).Error
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

func (r *RoleRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&entity.Role{}).Count(&count).Error
	return count, err
}
