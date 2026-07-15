package repository

import "github.com/aralvesandrade/cafeos/internal/domain/entity"

type PlanRepository interface {
	Create(plan *entity.Plan) error
	GetByID(id string) (*entity.Plan, error)
	GetBySlug(slug string) (*entity.Plan, error)
	List() ([]*entity.Plan, error)
	Update(plan *entity.Plan) error
	Delete(id string) error
}
