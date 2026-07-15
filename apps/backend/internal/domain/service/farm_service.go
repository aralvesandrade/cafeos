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

func validateProducer(producer *entity.Producer) error {
	if producer.Name == "" {
		return errors.New("producer name is required")
	}
	if producer.UserID == "" {
		return errors.New("producer user_id is required")
	}
	if producer.RoleID == "" {
		return errors.New("producer role_id is required")
	}
	return nil
}

// Create registers a new farm and, optionally, the users linked to it (each
// with a role scoped to this farm — see entity.Producer). producers may be
// empty when no links are provided yet.
func (s *FarmService) Create(farm *entity.Farm, producers []*entity.Producer) (*entity.Farm, error) {
	if err := validateFarm(farm); err != nil {
		return nil, err
	}

	farm.ID = uuid.New().String()
	farm.CreatedAt = time.Now()
	farm.UpdatedAt = time.Now()

	if err := s.repo.Create(farm); err != nil {
		return nil, err
	}

	for _, producer := range producers {
		if err := validateProducer(producer); err != nil {
			_ = s.repo.Delete(farm.ID)
			return nil, err
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
	}
	if len(producers) > 0 {
		farm.Producers = make([]entity.Producer, len(producers))
		for i, p := range producers {
			farm.Producers[i] = *p
		}
	}

	return farm, nil
}

func (s *FarmService) GetByID(id string) (*entity.Farm, error) {
	return s.repo.GetByID(id)
}

func (s *FarmService) ListByOrganization(organizationID string) ([]*entity.Farm, error) {
	return s.repo.ListByOrganization(organizationID)
}

// ListByOwner returns only the farms the given user is linked to (used to
// scope the proprietario role to the farms it actually has a link to).
func (s *FarmService) ListByOwner(organizationID, userID string) ([]*entity.Farm, error) {
	links, err := s.producerRepo.ListByUserID(userID)
	if err != nil {
		return nil, err
	}
	owned := make([]*entity.Farm, 0, len(links))
	for _, link := range links {
		if link.OrganizationID != organizationID {
			continue
		}
		farm, err := s.repo.GetByID(link.FarmID)
		if err != nil {
			continue
		}
		owned = append(owned, farm)
	}
	return owned, nil
}

// OwnedFarmIDs returns the set of farm IDs the given user is linked to, for
// cross-referencing plot/operation/financial/etc. records to that user's
// farms.
func (s *FarmService) OwnedFarmIDs(organizationID, userID string) (map[string]bool, error) {
	links, err := s.producerRepo.ListByUserID(userID)
	if err != nil {
		return nil, err
	}
	ids := make(map[string]bool, len(links))
	for _, link := range links {
		if link.OrganizationID == organizationID {
			ids[link.FarmID] = true
		}
	}
	return ids, nil
}

// IsOwner reports whether the given user has a link to the farm.
func (s *FarmService) IsOwner(farmID, userID string) bool {
	exists, err := s.producerRepo.ExistsByFarmAndUser(farmID, userID)
	return err == nil && exists
}

func (s *FarmService) Update(farm *entity.Farm) error {
	if err := validateFarm(farm); err != nil {
		return err
	}
	return s.repo.Update(farm)
}

// SetProducers replaces every user-link of a farm with the given list — the
// same "PUT replaces state" pattern used for the permissions matrix.
func (s *FarmService) SetProducers(farmID string, producers []*entity.Producer) ([]*entity.Producer, error) {
	farm, err := s.repo.GetByID(farmID)
	if err != nil {
		return nil, errors.New("farm not found")
	}

	for _, producer := range producers {
		if err := validateProducer(producer); err != nil {
			return nil, err
		}
	}

	if err := s.producerRepo.DeleteByFarmID(farmID); err != nil {
		return nil, err
	}

	for _, producer := range producers {
		producer.ID = uuid.New().String()
		producer.OrganizationID = farm.OrganizationID
		producer.FarmID = farmID
		producer.CreatedAt = time.Now()
		producer.UpdatedAt = time.Now()
		if err := s.producerRepo.Create(producer); err != nil {
			return nil, err
		}
	}
	return producers, nil
}

func (s *FarmService) Delete(id string) error {
	if err := s.producerRepo.DeleteByFarmID(id); err != nil {
		return err
	}
	return s.repo.Delete(id)
}
