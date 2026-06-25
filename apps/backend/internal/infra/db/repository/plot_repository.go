package repository

import (
	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	"gorm.io/gorm"
)

type PlotRepository struct {
	db *gorm.DB
}

func NewPlotRepository(db *gorm.DB) *PlotRepository {
	return &PlotRepository{db: db}
}

func (r *PlotRepository) WithTx(tx *gorm.DB) *PlotRepository {
	return &PlotRepository{db: tx}
}

func (r *PlotRepository) Create(p *entity.Plot) error {
	return r.db.Create(p).Error
}

func (r *PlotRepository) GetByID(id string) (*entity.Plot, error) {
	var p entity.Plot
	err := r.db.First(&p, "id = ?", id).Error
	return &p, err
}

func (r *PlotRepository) ListByFarm(farmID string) ([]*entity.Plot, error) {
	var plots []*entity.Plot
	err := r.db.Where("farm_id = ?", farmID).Order("name").Find(&plots).Error
	return plots, err
}

func (r *PlotRepository) ListByTenant(tenantID string) ([]*entity.Plot, error) {
	var plots []*entity.Plot
	err := r.db.Where("tenant_id = ?", tenantID).Order("name").Find(&plots).Error
	return plots, err
}

func (r *PlotRepository) Update(p *entity.Plot) error {
	return r.db.Save(p).Error
}

func (r *PlotRepository) Delete(id string) error {
	return r.db.Delete(&entity.Plot{}, "id = ?", id).Error
}
