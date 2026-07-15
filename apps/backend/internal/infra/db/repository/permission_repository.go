package repository

import (
	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PermissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) *PermissionRepository {
	return &PermissionRepository{db: db}
}

func (r *PermissionRepository) ListByOrganization(organizationID string) ([]*entity.RolePermission, error) {
	var permissions []*entity.RolePermission
	err := r.db.Where("organization_id = ?", organizationID).Find(&permissions).Error
	return permissions, err
}

func (r *PermissionRepository) Upsert(permission *entity.RolePermission) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "organization_id"}, {Name: "role"}, {Name: "module"}},
		DoUpdates: clause.AssignmentColumns([]string{"access", "updated_at"}),
	}).Create(permission).Error
}

func (r *PermissionRepository) CountByOrganization(organizationID string) (int64, error) {
	var count int64
	err := r.db.Model(&entity.RolePermission{}).Where("organization_id = ?", organizationID).Count(&count).Error
	return count, err
}
