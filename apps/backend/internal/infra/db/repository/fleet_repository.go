package repository

import (
	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	"gorm.io/gorm"
)

type VehicleRepository struct {
	db *gorm.DB
}

func NewVehicleRepository(db *gorm.DB) *VehicleRepository {
	return &VehicleRepository{db: db}
}

func (r *VehicleRepository) WithTx(tx *gorm.DB) *VehicleRepository {
	return &VehicleRepository{db: tx}
}

func (r *VehicleRepository) Create(v *entity.Vehicle) error {
	return r.db.Create(v).Error
}

func (r *VehicleRepository) GetByID(id string) (*entity.Vehicle, error) {
	var v entity.Vehicle
	err := r.db.First(&v, "id = ?", id).Error
	return &v, err
}

func (r *VehicleRepository) ListByTenant(tenantID string) ([]*entity.Vehicle, error) {
	var items []*entity.Vehicle
	err := r.db.Where("tenant_id = ?", tenantID).Order("name").Find(&items).Error
	return items, err
}

func (r *VehicleRepository) Update(v *entity.Vehicle) error {
	return r.db.Save(v).Error
}

func (r *VehicleRepository) Delete(id string) error {
	return r.db.Delete(&entity.Vehicle{}, "id = ?", id).Error
}

type MaintenanceRepository struct {
	db *gorm.DB
}

func NewMaintenanceRepository(db *gorm.DB) *MaintenanceRepository {
	return &MaintenanceRepository{db: db}
}

func (r *MaintenanceRepository) WithTx(tx *gorm.DB) *MaintenanceRepository {
	return &MaintenanceRepository{db: tx}
}

func (r *MaintenanceRepository) Create(m *entity.Maintenance) error {
	return r.db.Create(m).Error
}

func (r *MaintenanceRepository) GetByID(id string) (*entity.Maintenance, error) {
	var m entity.Maintenance
	err := r.db.Preload("Vehicle").First(&m, "id = ?", id).Error
	return &m, err
}

func (r *MaintenanceRepository) ListByTenant(tenantID string) ([]*entity.Maintenance, error) {
	var items []*entity.Maintenance
	err := r.db.Preload("Vehicle").Where("tenant_id = ?", tenantID).Order("date DESC").Find(&items).Error
	return items, err
}

func (r *MaintenanceRepository) ListByVehicle(vehicleID string) ([]*entity.Maintenance, error) {
	var items []*entity.Maintenance
	err := r.db.Where("vehicle_id = ?", vehicleID).Order("date DESC").Find(&items).Error
	return items, err
}

func (r *MaintenanceRepository) Update(m *entity.Maintenance) error {
	return r.db.Save(m).Error
}

func (r *MaintenanceRepository) Delete(id string) error {
	return r.db.Delete(&entity.Maintenance{}, "id = ?", id).Error
}
