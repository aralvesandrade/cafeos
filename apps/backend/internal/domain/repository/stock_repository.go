package repository

import "github.com/aralvesandrade/cafeos/internal/domain/entity"

type StockItemRepository interface {
	Create(item *entity.StockItem) error
	GetByID(id string) (*entity.StockItem, error)
	ListByOrganization(organizationID string) ([]*entity.StockItem, error)
	Update(item *entity.StockItem) error
	Delete(id string) error
}

type StockMovementRepository interface {
	Create(mov *entity.StockMovement) error
	ListByOrganization(organizationID string) ([]*entity.StockMovement, error)
	ListByItem(itemID string) ([]*entity.StockMovement, error)
}
