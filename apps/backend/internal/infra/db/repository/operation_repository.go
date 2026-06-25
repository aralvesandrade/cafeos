package repository

import (
	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	"gorm.io/gorm"
)

type OperationRepository struct {
	db *gorm.DB
}

func NewOperationRepository(db *gorm.DB) *OperationRepository {
	return &OperationRepository{db: db}
}

func (r *OperationRepository) WithTx(tx *gorm.DB) *OperationRepository {
	return &OperationRepository{db: tx}
}

func (r *OperationRepository) Create(op *entity.Operation) error {
	return r.db.Create(op).Error
}

func (r *OperationRepository) GetByID(id string) (*entity.Operation, error) {
	var op entity.Operation
	err := r.db.Preload("Plot").First(&op, "id = ?", id).Error
	if err == nil {
		op.PlotName = op.Plot.Name
	}
	return &op, err
}

func (r *OperationRepository) ListByPlot(plotID string) ([]*entity.Operation, error) {
	var ops []*entity.Operation
	err := r.db.Preload("Plot").Where("plot_id = ?", plotID).Order("date DESC").Find(&ops).Error
	for _, op := range ops {
		op.PlotName = op.Plot.Name
	}
	return ops, err
}

func (r *OperationRepository) ListByTenant(tenantID string) ([]*entity.Operation, error) {
	var ops []*entity.Operation
	err := r.db.Preload("Plot").Where("tenant_id = ?", tenantID).Order("date DESC").Find(&ops).Error
	for _, op := range ops {
		op.PlotName = op.Plot.Name
	}
	return ops, err
}

func (r *OperationRepository) ListByTenantAndPeriod(tenantID string, start, end string) ([]*entity.Operation, error) {
	var ops []*entity.Operation
	err := r.db.Preload("Plot").Where("tenant_id = ? AND date >= ? AND date <= ?", tenantID, start, end).
		Order("date DESC").Find(&ops).Error
	for _, op := range ops {
		op.PlotName = op.Plot.Name
	}
	return ops, err
}

func (r *OperationRepository) Delete(id string) error {
	return r.db.Delete(&entity.Operation{}, "id = ?", id).Error
}
