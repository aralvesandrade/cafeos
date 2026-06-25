package repository

import (
	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	"gorm.io/gorm"
)

type StockItemRepository struct {
	db *gorm.DB
}

func NewStockItemRepository(db *gorm.DB) *StockItemRepository {
	return &StockItemRepository{db: db}
}

func (r *StockItemRepository) WithTx(tx *gorm.DB) *StockItemRepository {
	return &StockItemRepository{db: tx}
}

func (r *StockItemRepository) Create(item *entity.StockItem) error {
	return r.db.Create(item).Error
}

func (r *StockItemRepository) GetByID(id string) (*entity.StockItem, error) {
	var item entity.StockItem
	err := r.db.Preload("Product").First(&item, "id = ?", id).Error
	return &item, err
}

func (r *StockItemRepository) ListByTenant(tenantID string) ([]*entity.StockItem, error) {
	var items []*entity.StockItem
	err := r.db.Preload("Product").Where("tenant_id = ?", tenantID).Order("product_id").Find(&items).Error
	return items, err
}

func (r *StockItemRepository) Update(item *entity.StockItem) error {
	return r.db.Save(item).Error
}

func (r *StockItemRepository) Delete(id string) error {
	return r.db.Delete(&entity.StockItem{}, "id = ?", id).Error
}

type StockMovementRepository struct {
	db *gorm.DB
}

func NewStockMovementRepository(db *gorm.DB) *StockMovementRepository {
	return &StockMovementRepository{db: db}
}

func (r *StockMovementRepository) WithTx(tx *gorm.DB) *StockMovementRepository {
	return &StockMovementRepository{db: tx}
}

func (r *StockMovementRepository) Create(mov *entity.StockMovement) error {
	return r.db.Create(mov).Error
}

func (r *StockMovementRepository) ListByTenant(tenantID string) ([]*entity.StockMovement, error) {
	var movs []*entity.StockMovement
	err := r.db.Where("tenant_id = ?", tenantID).Order("date DESC").Find(&movs).Error
	return movs, err
}

func (r *StockMovementRepository) ListByItem(itemID string) ([]*entity.StockMovement, error) {
	var movs []*entity.StockMovement
	err := r.db.Where("item_id = ?", itemID).Order("date DESC").Find(&movs).Error
	return movs, err
}
