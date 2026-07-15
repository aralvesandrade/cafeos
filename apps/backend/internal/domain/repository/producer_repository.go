package repository

import "github.com/aralvesandrade/cafeos/internal/domain/entity"

type ProducerRepository interface {
	Create(producer *entity.Producer) error
	ListByFarmID(farmID string) ([]*entity.Producer, error)
	// ListByUserID returns every farm-link (across every organization) for
	// the given user — callers filter by organization when needed, since a
	// user only ever belongs to one organization.
	ListByUserID(userID string) ([]*entity.Producer, error)
	ExistsByFarmAndUser(farmID, userID string) (bool, error)
	Update(producer *entity.Producer) error
	DeleteByFarmID(farmID string) error
}
