package service

import (
	"errors"
	"time"

	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	"github.com/aralvesandrade/cafeos/internal/domain/repository"
	"github.com/google/uuid"
)

type StockService struct {
	itemRepo repository.StockItemRepository
	movRepo  repository.StockMovementRepository
}

func NewStockService(itemRepo repository.StockItemRepository, movRepo repository.StockMovementRepository) *StockService {
	return &StockService{itemRepo: itemRepo, movRepo: movRepo}
}

func (s *StockService) CreateItem(organizationID, productID, unit, batch, location, notes string, farmID *string, quantity, minStock float64, expiryDate *time.Time) (*entity.StockItem, error) {
	if productID == "" {
		return nil, errors.New("product is required")
	}
	item := &entity.StockItem{
		ID:             uuid.New().String(),
		OrganizationID: organizationID,
		ProductID:      productID,
		FarmID:         farmID,
		Quantity:       quantity,
		Unit:           unit,
		Batch:          batch,
		ExpiryDate:     expiryDate,
		MinStock:       minStock,
		Location:       location,
		Notes:          notes,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	if err := s.itemRepo.Create(item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *StockService) GetItemByID(id string) (*entity.StockItem, error) {
	return s.itemRepo.GetByID(id)
}

func (s *StockService) ListItems(organizationID string) ([]*entity.StockItem, error) {
	return s.itemRepo.ListByOrganization(organizationID)
}

func (s *StockService) UpdateItem(item *entity.StockItem) error {
	item.UpdatedAt = time.Now()
	return s.itemRepo.Update(item)
}

func (s *StockService) DeleteItem(id string) error {
	return s.itemRepo.Delete(id)
}

func (s *StockService) GetMovementByID(id string) (*entity.StockMovement, error) {
	return s.movRepo.GetByID(id)
}

func (s *StockService) RecordMovement(organizationID, itemID, movType, reference, notes string, quantity float64, date time.Time) (*entity.StockMovement, error) {
	if quantity <= 0 {
		return nil, errors.New("quantity must be greater than zero")
	}
	mov := &entity.StockMovement{
		ID:             uuid.New().String(),
		OrganizationID: organizationID,
		ItemID:         itemID,
		Type:           movType,
		Quantity:       quantity,
		Date:           date,
		Reference:      reference,
		Notes:          notes,
		CreatedAt:      time.Now(),
	}
	if err := s.movRepo.Create(mov); err != nil {
		return nil, err
	}
	// Update item quantity
	item, err := s.itemRepo.GetByID(itemID)
	if err == nil {
		if movType == "in" {
			item.Quantity += quantity
		} else {
			item.Quantity -= quantity
		}
		s.itemRepo.Update(item)
	}
	return mov, nil
}

func (s *StockService) ListMovements(organizationID string) ([]*entity.StockMovement, error) {
	return s.movRepo.ListByOrganization(organizationID)
}

func (s *StockService) ListMovementsByItem(itemID string) ([]*entity.StockMovement, error) {
	return s.movRepo.ListByItem(itemID)
}
