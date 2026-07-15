package repository

import "github.com/aralvesandrade/cafeos/internal/domain/entity"

type PermissionRepository interface {
	ListByOrganization(organizationID string) ([]*entity.RolePermission, error)
	Upsert(permission *entity.RolePermission) error
	CountByOrganization(organizationID string) (int64, error)
}
