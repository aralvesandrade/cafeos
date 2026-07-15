package repository

import (
	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	"gorm.io/gorm"
)

type PlanRepository struct {
	db *gorm.DB
}

func NewPlanRepository(db *gorm.DB) *PlanRepository {
	return &PlanRepository{db: db}
}

func (r *PlanRepository) Create(p *entity.Plan) error {
	return r.db.Create(p).Error
}

func (r *PlanRepository) GetByID(id string) (*entity.Plan, error) {
	var p entity.Plan
	err := r.db.First(&p, "id = ?", id).Error
	return &p, err
}

func (r *PlanRepository) GetBySlug(slug string) (*entity.Plan, error) {
	var p entity.Plan
	err := r.db.Where("slug = ?", slug).First(&p).Error
	return &p, err
}

func (r *PlanRepository) List() ([]*entity.Plan, error) {
	var plans []*entity.Plan
	err := r.db.Order("display_order, name").Find(&plans).Error
	return plans, err
}

func (r *PlanRepository) Update(p *entity.Plan) error {
	return r.db.Save(p).Error
}

func (r *PlanRepository) Delete(id string) error {
	return r.db.Delete(&entity.Plan{}, "id = ?", id).Error
}
