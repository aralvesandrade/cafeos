package service

import (
	"errors"
	"time"

	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	"github.com/aralvesandrade/cafeos/internal/domain/repository"
	"github.com/google/uuid"
)

type FleetService struct {
	vehRepo   repository.VehicleRepository
	maintRepo repository.MaintenanceRepository
}

func NewFleetService(vehRepo repository.VehicleRepository, maintRepo repository.MaintenanceRepository) *FleetService {
	return &FleetService{vehRepo: vehRepo, maintRepo: maintRepo}
}

func (s *FleetService) CreateVehicle(organizationID, name, vehType, plate, brand, model string, year int) (*entity.Vehicle, error) {
	if name == "" {
		return nil, errors.New("vehicle name is required")
	}
	v := &entity.Vehicle{
		ID:             uuid.New().String(),
		OrganizationID: organizationID,
		Name:           name,
		Type:           entity.VehicleType(vehType),
		Plate:          plate,
		Brand:          brand,
		Model:          model,
		Year:           year,
		Status:         "active",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	if err := s.vehRepo.Create(v); err != nil {
		return nil, err
	}
	return v, nil
}

func (s *FleetService) GetVehicleByID(id string) (*entity.Vehicle, error) {
	return s.vehRepo.GetByID(id)
}

func (s *FleetService) ListVehicles(organizationID string) ([]*entity.Vehicle, error) {
	return s.vehRepo.ListByOrganization(organizationID)
}

func (s *FleetService) UpdateVehicle(v *entity.Vehicle) error {
	v.UpdatedAt = time.Now()
	return s.vehRepo.Update(v)
}

func (s *FleetService) DeleteVehicle(id string) error {
	return s.vehRepo.Delete(id)
}

func (s *FleetService) CreateMaintenance(organizationID, vehicleID, maintType, description, notes string, cost, odometer float64, date time.Time) (*entity.Maintenance, error) {
	if vehicleID == "" {
		return nil, errors.New("vehicle is required")
	}
	m := &entity.Maintenance{
		ID:             uuid.New().String(),
		OrganizationID: organizationID,
		VehicleID:      vehicleID,
		Date:           date,
		Type:           maintType,
		Description:    description,
		Cost:           cost,
		Odometer:       odometer,
		Notes:          notes,
		CreatedAt:      time.Now(),
	}
	if err := s.maintRepo.Create(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *FleetService) GetMaintenanceByID(id string) (*entity.Maintenance, error) {
	return s.maintRepo.GetByID(id)
}

func (s *FleetService) ListMaintenance(organizationID string) ([]*entity.Maintenance, error) {
	return s.maintRepo.ListByOrganization(organizationID)
}

func (s *FleetService) ListMaintenanceByVehicle(vehicleID string) ([]*entity.Maintenance, error) {
	return s.maintRepo.ListByVehicle(vehicleID)
}

func (s *FleetService) UpdateMaintenance(m *entity.Maintenance) error {
	return s.maintRepo.Update(m)
}

func (s *FleetService) DeleteMaintenance(id string) error {
	return s.maintRepo.Delete(id)
}
