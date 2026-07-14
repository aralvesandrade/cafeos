package service

import (
	"errors"
	"time"

	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	"github.com/aralvesandrade/cafeos/internal/domain/repository"
	"github.com/google/uuid"
)

type FinancialService struct {
	repo repository.FinancialRepository
}

func NewFinancialService(repo repository.FinancialRepository) *FinancialService {
	return &FinancialService{repo: repo}
}

func (s *FinancialService) Create(organizationID, txType string, costCenterID, farmID *string, description string, amount float64, date time.Time, dueDate time.Time, notes string) (*entity.FinancialTransaction, error) {
	if description == "" {
		return nil, errors.New("description is required")
	}
	if amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}
	tx := &entity.FinancialTransaction{
		ID:             uuid.New().String(),
		OrganizationID: organizationID,
		Type:           entity.TransactionType(txType),
		CostCenterID:   costCenterID,
		FarmID:         farmID,
		Description:    description,
		Amount:         amount,
		Date:           date,
		DueDate:        dueDate,
		Status:         "pending",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	if err := s.repo.Create(tx); err != nil {
		return nil, err
	}
	return tx, nil
}

func (s *FinancialService) GetByID(id string) (*entity.FinancialTransaction, error) {
	return s.repo.GetByID(id)
}

func (s *FinancialService) ListByOrganization(organizationID string) ([]*entity.FinancialTransaction, error) {
	return s.repo.ListByOrganization(organizationID)
}

func (s *FinancialService) Update(tx *entity.FinancialTransaction) error {
	if tx.Description == "" {
		return errors.New("description is required")
	}
	tx.UpdatedAt = time.Now()
	return s.repo.Update(tx)
}

func (s *FinancialService) Delete(id string) error {
	return s.repo.Delete(id)
}
