package repository

import (
	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	"gorm.io/gorm"
)

type ProducerRepository struct {
	db *gorm.DB
}

func NewProducerRepository(db *gorm.DB) *ProducerRepository {
	return &ProducerRepository{db: db}
}

func (r *ProducerRepository) WithTx(tx *gorm.DB) *ProducerRepository {
	return &ProducerRepository{db: tx}
}

func (r *ProducerRepository) Create(p *entity.Producer) error {
	return r.db.Create(p).Error
}

func (r *ProducerRepository) GetByFarmID(farmID string) (*entity.Producer, error) {
	var p entity.Producer
	err := r.db.First(&p, "farm_id = ?", farmID).Error
	return &p, err
}

func (r *ProducerRepository) Update(p *entity.Producer) error {
	return r.db.Save(p).Error
}

func (r *ProducerRepository) DeleteByFarmID(farmID string) error {
	return r.db.Where("farm_id = ?", farmID).Delete(&entity.Producer{}).Error
}
