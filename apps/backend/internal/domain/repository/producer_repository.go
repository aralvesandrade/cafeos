package repository

import "github.com/aralvesandrade/cafeos/internal/domain/entity"

type ProducerRepository interface {
	Create(producer *entity.Producer) error
	GetByFarmID(farmID string) (*entity.Producer, error)
	Update(producer *entity.Producer) error
	DeleteByFarmID(farmID string) error
}
