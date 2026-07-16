package repository

import (
	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	"gorm.io/gorm"
)

type OperationTypeRepository struct {
	db *gorm.DB
}

func NewOperationTypeRepository(db *gorm.DB) *OperationTypeRepository {
	return &OperationTypeRepository{db: db}
}

func (r *OperationTypeRepository) WithTx(tx *gorm.DB) *OperationTypeRepository {
	return &OperationTypeRepository{db: tx}
}

func (r *OperationTypeRepository) Create(ot *entity.OperationType) error {
	return r.db.Create(ot).Error
}

func (r *OperationTypeRepository) GetByID(id string) (*entity.OperationType, error) {
	var ot entity.OperationType
	err := r.db.First(&ot, "id = ?", id).Error
	return &ot, err
}

func (r *OperationTypeRepository) GetByOrganizationAndCode(organizationID, code string) (*entity.OperationType, error) {
	var ot entity.OperationType
	err := r.db.Where("organization_id = ? AND code = ?", organizationID, code).First(&ot).Error
	return &ot, err
}

func (r *OperationTypeRepository) ListByOrganization(organizationID string) ([]*entity.OperationType, error) {
	var ots []*entity.OperationType
	err := r.db.Where("organization_id = ?", organizationID).Order("name").Find(&ots).Error
	return ots, err
}

func (r *OperationTypeRepository) Update(ot *entity.OperationType) error {
	return r.db.Save(ot).Error
}

func (r *OperationTypeRepository) Delete(id string) error {
	return r.db.Delete(&entity.OperationType{}, "id = ?", id).Error
}
