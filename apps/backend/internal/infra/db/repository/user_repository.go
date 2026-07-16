package repository

import (
	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) WithTx(tx *gorm.DB) *UserRepository {
	return &UserRepository{db: tx}
}

func (r *UserRepository) Create(u *entity.User) error {
	return r.db.Create(u).Error
}

func (r *UserRepository) GetByID(id string) (*entity.User, error) {
	var u entity.User
	err := r.db.Preload("Role").Preload("Plan").Preload("ManagedByUser").First(&u, "id = ?", id).Error
	return &u, err
}

func (r *UserRepository) GetByEmail(email string) (*entity.User, error) {
	var u entity.User
	err := r.db.Preload("Role").Preload("Plan").Preload("ManagedByUser").Where("email = ?", email).First(&u).Error
	return &u, err
}

func (r *UserRepository) ListByOrganization(organizationID string) ([]*entity.User, error) {
	var users []*entity.User
	err := r.db.Preload("Role").Preload("Plan").Preload("ManagedByUser").Where("organization_id = ?", organizationID).Order("name").Find(&users).Error
	return users, err
}

func (r *UserRepository) ListByManager(managerID string) ([]*entity.User, error) {
	var users []*entity.User
	err := r.db.Preload("Role").Preload("Plan").Preload("ManagedByUser").
		Where("id = ? OR managed_by_user_id = ?", managerID, managerID).
		Order("name").Find(&users).Error
	return users, err
}

func (r *UserRepository) Update(u *entity.User) error {
	return r.db.Save(u).Error
}

func (r *UserRepository) List() ([]*entity.User, error) {
	var users []*entity.User
	err := r.db.Preload("Role").Preload("Plan").Preload("ManagedByUser").Order("name").Find(&users).Error
	return users, err
}

func (r *UserRepository) Delete(id string) error {
	return r.db.Delete(&entity.User{}, "id = ?", id).Error
}

func (r *UserRepository) CountByRole(roleID string) (int64, error) {
	var count int64
	err := r.db.Model(&entity.User{}).Where("role_id = ?", roleID).Count(&count).Error
	return count, err
}

func (r *UserRepository) CountByOrganization(organizationID string) (int64, error) {
	var count int64
	err := r.db.Model(&entity.User{}).Where("organization_id = ?", organizationID).Count(&count).Error
	return count, err
}
