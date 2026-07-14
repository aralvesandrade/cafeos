package repository

import (
	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	"gorm.io/gorm"
)

type HarvestRepository struct {
	db *gorm.DB
}

func NewHarvestRepository(db *gorm.DB) *HarvestRepository {
	return &HarvestRepository{db: db}
}

func (r *HarvestRepository) WithTx(tx *gorm.DB) *HarvestRepository {
	return &HarvestRepository{db: tx}
}

func (r *HarvestRepository) Create(h *entity.Harvest) error {
	return r.db.Create(h).Error
}

func (r *HarvestRepository) GetByID(id string) (*entity.Harvest, error) {
	var h entity.Harvest
	err := r.db.First(&h, "id = ?", id).Error
	return &h, err
}

func (r *HarvestRepository) ListByOrganization(organizationID string) ([]*entity.Harvest, error) {
	var harvests []*entity.Harvest
	err := r.db.Where("organization_id = ?", organizationID).Order("year DESC").Find(&harvests).Error
	return harvests, err
}

func (r *HarvestRepository) Update(h *entity.Harvest) error {
	return r.db.Save(h).Error
}

type HarvestProductionRepository struct {
	db *gorm.DB
}

func NewHarvestProductionRepository(db *gorm.DB) *HarvestProductionRepository {
	return &HarvestProductionRepository{db: db}
}

func (r *HarvestProductionRepository) WithTx(tx *gorm.DB) *HarvestProductionRepository {
	return &HarvestProductionRepository{db: tx}
}

func (r *HarvestProductionRepository) Create(hp *entity.HarvestProduction) error {
	return r.db.Create(hp).Error
}

func (r *HarvestProductionRepository) GetByID(id string) (*entity.HarvestProduction, error) {
	var hp entity.HarvestProduction
	err := r.db.First(&hp, "id = ?", id).Error
	return &hp, err
}

func (r *HarvestProductionRepository) ListByHarvest(harvestID string) ([]*entity.HarvestProduction, error) {
	var prods []*entity.HarvestProduction
	err := r.db.Where("harvest_id = ?", harvestID).Order("plot_id").Find(&prods).Error
	return prods, err
}

func (r *HarvestProductionRepository) ListByPlot(plotID string) ([]*entity.HarvestProduction, error) {
	var prods []*entity.HarvestProduction
	err := r.db.Where("plot_id = ?", plotID).Order("recorded_at DESC").Find(&prods).Error
	return prods, err
}

func (r *HarvestProductionRepository) Update(hp *entity.HarvestProduction) error {
	return r.db.Save(hp).Error
}
