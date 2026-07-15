package repository

import "github.com/aralvesandrade/cafeos/internal/domain/entity"

type RoleRepository interface {
	// ListForOrganization returns every role usable within an organization:
	// the global system roles (organization_id NULL) plus the organization's
	// own rows.
	ListForOrganization(organizationID string) ([]*entity.Role, error)
	GetByID(id string) (*entity.Role, error)
	// FindByKey resolves a role key within an organization, preferring an
	// org-scoped row and falling back to a global system role.
	FindByKey(organizationID, key string) (*entity.Role, error)
	Create(role *entity.Role) error
	Update(role *entity.Role) error
	Delete(id string) error
	CountByOrganization(organizationID string) (int64, error)
}
