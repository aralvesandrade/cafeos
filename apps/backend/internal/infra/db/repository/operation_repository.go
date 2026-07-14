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

func (r *OperationRepository) Update(op *entity.Operation) error {
	return r.db.Save(op).Error
}

func populateOperationType(op *entity.Operation) {
	if op.Type != nil {
		op.TypeName = op.Type.Name
		op.TypeColor = op.Type.Color
	}
}

func (r *OperationRepository) GetByID(id string) (*entity.Operation, error) {
	var op entity.Operation
	err := r.db.Preload("Plot").Preload("Type").First(&op, "id = ?", id).Error
	if err == nil {
		op.PlotName = op.Plot.Name
		populateOperationType(&op)
	}
	return &op, err
}

func (r *OperationRepository) ListByPlot(plotID string) ([]*entity.Operation, error) {
	var ops []*entity.Operation
	err := r.db.Preload("Plot").Preload("Type").Where("plot_id = ?", plotID).Order("date DESC").Find(&ops).Error
	for _, op := range ops {
		op.PlotName = op.Plot.Name
		populateOperationType(op)
	}
	return ops, err
}

func (r *OperationRepository) ListByOrganization(organizationID string) ([]*entity.Operation, error) {
	var ops []*entity.Operation
	err := r.db.Preload("Plot").Preload("Type").Where("organization_id = ?", organizationID).Order("date DESC").Find(&ops).Error
	for _, op := range ops {
		op.PlotName = op.Plot.Name
		populateOperationType(op)
	}
	return ops, err
}

func (r *OperationRepository) ListByOrganizationAndPeriod(organizationID string, start, end string) ([]*entity.Operation, error) {
	var ops []*entity.Operation
	err := r.db.Preload("Plot").Preload("Type").Where("organization_id = ? AND date >= ? AND date <= ?", organizationID, start, end).
		Order("date DESC").Find(&ops).Error
	for _, op := range ops {
		op.PlotName = op.Plot.Name
		populateOperationType(op)
	}
	return ops, err
}

func (r *OperationRepository) Delete(id string) error {
	return r.db.Delete(&entity.Operation{}, "id = ?", id).Error
}
