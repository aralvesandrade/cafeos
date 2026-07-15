package entity

import "time"

type AccessLevel string

const (
	AccessNone  AccessLevel = "none"
	AccessRead  AccessLevel = "read"
	AccessWrite AccessLevel = "write"
)

// Satisfies reports whether this access level meets the required level
// (write satisfies read; read/none do not satisfy write).
func (a AccessLevel) Satisfies(need AccessLevel) bool {
	switch need {
	case AccessRead:
		return a == AccessRead || a == AccessWrite
	case AccessWrite:
		return a == AccessWrite
	default:
		return true
	}
}

type Module string

const (
	ModuleDashboard   Module = "dashboard"
	ModuleFarms       Module = "farms"
	ModuleOperations  Module = "operations"
	ModuleHarvests    Module = "harvests"
	ModuleResources   Module = "resources"
	ModuleFinancial   Module = "financial"
	ModuleUsers       Module = "users"
	ModulePermissions Module = "permissions"
)

var AllModules = []Module{
	ModuleDashboard,
	ModuleFarms,
	ModuleOperations,
	ModuleHarvests,
	ModuleResources,
	ModuleFinancial,
	ModuleUsers,
	ModulePermissions,
}

// RolePermission is the per-organization, per-role, per-module access
// level configured by an organization admin through the Permissions screen.
type RolePermission struct {
	ID             string      `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	OrganizationID string      `json:"organization_id" gorm:"type:uuid;not null;uniqueIndex:idx_org_role_module"`
	Role           UserRole    `json:"role" gorm:"not null;uniqueIndex:idx_org_role_module"`
	Module         Module      `json:"module" gorm:"not null;uniqueIndex:idx_org_role_module"`
	Access         AccessLevel `json:"access" gorm:"not null;default:'none'"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
}
