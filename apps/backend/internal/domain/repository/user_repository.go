package repository

import "github.com/aralvesandrade/cafeos/internal/domain/entity"

type UserRepository interface {
	Create(user *entity.User) error
	GetByID(id string) (*entity.User, error)
	GetByEmail(email string) (*entity.User, error)
	ListByOrganization(organizationID string) ([]*entity.User, error)
	List() ([]*entity.User, error)
	Update(user *entity.User) error
	Delete(id string) error
}
