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

func (r *ProducerRepository) ListByFarmID(farmID string) ([]*entity.Producer, error) {
	var producers []*entity.Producer
	err := r.db.Where("farm_id = ?", farmID).Find(&producers).Error
	return producers, err
}

func (r *ProducerRepository) ListByUserID(userID string) ([]*entity.Producer, error) {
	var producers []*entity.Producer
	err := r.db.Where("user_id = ?", userID).Find(&producers).Error
	return producers, err
}

func (r *ProducerRepository) ExistsByFarmAndUser(farmID, userID string) (bool, error) {
	var count int64
	err := r.db.Model(&entity.Producer{}).Where("farm_id = ? AND user_id = ?", farmID, userID).Count(&count).Error
	return count > 0, err
}

func (r *ProducerRepository) Update(p *entity.Producer) error {
	return r.db.Save(p).Error
}

func (r *ProducerRepository) CreateBatch(producers []*entity.Producer) error {
	if len(producers) == 0 {
		return nil
	}
	return r.db.Create(producers).Error
}

func (r *ProducerRepository) DeleteByUserID(userID string) error {
	return r.db.Where("user_id = ?", userID).Delete(&entity.Producer{}).Error
}

func (r *ProducerRepository) DeleteByFarmID(farmID string) error {
	return r.db.Where("farm_id = ?", farmID).Delete(&entity.Producer{}).Error
}
