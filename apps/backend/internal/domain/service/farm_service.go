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
		producer.OrganizationID = farm.OrganizationID
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

func (s *FarmService) ListByOrganization(organizationID string) ([]*entity.Farm, error) {
	return s.repo.ListByOrganization(organizationID)
}

// ListByOwner returns only the farms whose linked producer is the given user
// (used to scope the proprietario role to the farms it actually owns).
func (s *FarmService) ListByOwner(organizationID, userID string) ([]*entity.Farm, error) {
	farms, err := s.repo.ListByOrganization(organizationID)
	if err != nil {
		return nil, err
	}
	owned := make([]*entity.Farm, 0, len(farms))
	for _, f := range farms {
		if f.Producer != nil && f.Producer.UserID != nil && *f.Producer.UserID == userID {
			owned = append(owned, f)
		}
	}
	return owned, nil
}

// OwnedFarmIDs returns the set of farm IDs owned by the given user, for
// cross-referencing plot/operation/financial/etc. records to that user's farms.
func (s *FarmService) OwnedFarmIDs(organizationID, userID string) (map[string]bool, error) {
	farms, err := s.ListByOwner(organizationID, userID)
	if err != nil {
		return nil, err
	}
	ids := make(map[string]bool, len(farms))
	for _, f := range farms {
		ids[f.ID] = true
	}
	return ids, nil
}

// IsOwner reports whether the given user is the producer linked to the farm.
func (s *FarmService) IsOwner(farmID, userID string) bool {
	producer, err := s.producerRepo.GetByFarmID(farmID)
	if err != nil {
		return false
	}
	return producer.UserID != nil && *producer.UserID == userID
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
		producer.OrganizationID = farm.OrganizationID
		producer.FarmID = farmID
		producer.CreatedAt = time.Now()
		producer.UpdatedAt = time.Now()
		if err := s.producerRepo.Create(producer); err != nil {
			return nil, err
		}
		return producer, nil
	}

	producer.ID = existing.ID
	producer.OrganizationID = existing.OrganizationID
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
