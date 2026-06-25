package repository

import "github.com/aralvesandrade/cafeos/internal/domain/entity"

type VehicleRepository interface {
	Create(v *entity.Vehicle) error
	GetByID(id string) (*entity.Vehicle, error)
	ListByTenant(tenantID string) ([]*entity.Vehicle, error)
	Update(v *entity.Vehicle) error
	Delete(id string) error
}

type MaintenanceRepository interface {
	Create(m *entity.Maintenance) error
	GetByID(id string) (*entity.Maintenance, error)
	ListByTenant(tenantID string) ([]*entity.Maintenance, error)
	ListByVehicle(vehicleID string) ([]*entity.Maintenance, error)
	Update(m *entity.Maintenance) error
	Delete(id string) error
}
