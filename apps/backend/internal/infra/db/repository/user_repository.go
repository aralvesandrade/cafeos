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
	err := r.db.First(&u, "id = ?", id).Error
	return &u, err
}

func (r *UserRepository) GetByEmail(email string) (*entity.User, error) {
	var u entity.User
	err := r.db.Where("email = ?", email).First(&u).Error
	return &u, err
}

func (r *UserRepository) ListByTenant(tenantID string) ([]*entity.User, error) {
	var users []*entity.User
	err := r.db.Where("tenant_id = ?", tenantID).Order("name").Find(&users).Error
	return users, err
}

func (r *UserRepository) Update(u *entity.User) error {
	return r.db.Save(u).Error
}

func (r *UserRepository) List() ([]*entity.User, error) {
	var users []*entity.User
	err := r.db.Order("name").Find(&users).Error
	return users, err
}

func (r *UserRepository) Delete(id string) error {
	return r.db.Delete(&entity.User{}, "id = ?", id).Error
}
