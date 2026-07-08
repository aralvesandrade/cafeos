package service

import (
	"errors"
	"time"

	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	"github.com/aralvesandrade/cafeos/internal/domain/repository"
	"github.com/google/uuid"
)

type FarmService struct {
	repo         repository.FarmRepository
	producerRepo repository.ProducerRepository
}

func NewFarmService(repo repository.FarmRepository, producerRepo repository.ProducerRepository) *FarmService {
	return &FarmService{repo: repo, producerRepo: producerRepo}
}

func validateFarm(farm *entity.Farm) error {
	if farm.Name == "" {
		return errors.New("farm name is required")
	}
	if farm.TotalAreaHA <= 0 {
		return errors.New("total area must be greater than zero")
	}
	if farm.PlantedAreaHA > farm.TotalAreaHA {
		return errors.New("planted area cannot exceed total area")
	}
	return nil
}

// Create registers a new farm and, optionally, its responsible producer.
// producer may be nil when no producer data is provided yet.
func (s *FarmService) Create(farm *entity.Farm, producer *entity.Producer) (*entity.Farm, error) {
	if err := validateFarm(farm); err != nil {
		return nil, err
	}

	farm.ID = uuid.New().String()
	farm.CreatedAt = time.Now()
	farm.UpdatedAt = time.Now()

	if err := s.repo.Create(farm); err != nil {
		return nil, err
	}

	if producer != nil {
		if producer.Name == "" {
			return nil, errors.New("producer name is required")
		}
		producer.ID = uuid.New().String()
		producer.TenantID = farm.TenantID
		producer.FarmID = farm.ID
		producer.CreatedAt = time.Now()
		producer.UpdatedAt = time.Now()

		if err := s.producerRepo.Create(producer); err != nil {
			_ = s.repo.Delete(farm.ID)
			return nil, err
		}
		farm.Producer = producer
	}

	return farm, nil
}

func (s *FarmService) GetByID(id string) (*entity.Farm, error) {
	return s.repo.GetByID(id)
}

func (s *FarmService) ListByTenant(tenantID string) ([]*entity.Farm, error) {
	return s.repo.ListByTenant(tenantID)
}

func (s *FarmService) Update(farm *entity.Farm) error {
	if err := validateFarm(farm); err != nil {
		return err
	}
	return s.repo.Update(farm)
}

// UpsertProducer creates or updates the producer responsible for a farm.
func (s *FarmService) UpsertProducer(farmID string, producer *entity.Producer) (*entity.Producer, error) {
	if producer.Name == "" {
		return nil, errors.New("producer name is required")
	}

	farm, err := s.repo.GetByID(farmID)
	if err != nil {
		return nil, errors.New("farm not found")
	}

	existing, err := s.producerRepo.GetByFarmID(farmID)
	if err != nil {
		producer.ID = uuid.New().String()
		producer.TenantID = farm.TenantID
		producer.FarmID = farmID
		producer.CreatedAt = time.Now()
		producer.UpdatedAt = time.Now()
		if err := s.producerRepo.Create(producer); err != nil {
			return nil, err
		}
		return producer, nil
	}

	producer.ID = existing.ID
	producer.TenantID = existing.TenantID
	producer.FarmID = existing.FarmID
	producer.CreatedAt = existing.CreatedAt
	producer.UpdatedAt = time.Now()
	if err := s.producerRepo.Update(producer); err != nil {
		return nil, err
	}
	return producer, nil
}

func (s *FarmService) Delete(id string) error {
	if err := s.producerRepo.DeleteByFarmID(id); err != nil {
		return err
	}
	return s.repo.Delete(id)
}
