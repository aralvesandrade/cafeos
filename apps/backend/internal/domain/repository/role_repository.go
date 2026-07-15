package repository

import "github.com/aralvesandrade/cafeos/internal/domain/entity"

type RoleRepository interface {
	List() ([]*entity.Role, error)
	GetByID(id string) (*entity.Role, error)
	FindByKey(key string) (*entity.Role, error)
	Create(role *entity.Role) error
	Update(role *entity.Role) error
	Delete(id string) error
	Count() (int64, error)
}
