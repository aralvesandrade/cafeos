package entity

import "time"

// SystemRole keys — global, shared by every organization, not editable or
// deletable. Every other role is an organization-scoped row in the roles
// table (either seeded as a starter kit or created by an org admin).
const (
	SystemRolePlatformOwner     = "platform_owner"
	SystemRoleOrganizationAdmin = "organization_admin"
)

var SystemRoleKeys = []string{SystemRolePlatformOwner, SystemRoleOrganizationAdmin}

// DefaultOrgRoleKeys are seeded as editable, per-organization rows when an
// organization is created (or backfilled for organizations that predate
// this table). They are a starting kit, not system roles — an organization
// admin may rename or delete them once unused.
// RoleKeyProprietario is referenced by name in ownership-scoping logic (e.g.
// restrictedOwnerID in api/handler/access.go), since "proprietario" is the
// only role whose visibility is scoped to farms it owns rather than the
// whole organization. Organizations may still rename or delete this row
// like any other seeded default — the scoping check simply no-ops if it's
// gone.
const RoleKeyProprietario = "proprietario"

var DefaultOrgRoleKeys = []string{
	RoleKeyProprietario,
	"gerente_agricola",
	"engenheiro_agronomo",
	"tecnico_agricola",
	"operador_campo",
	"financeiro",
	"consultor_externo",
	"auditor",
}

var DefaultOrgRoleNames = map[string]string{
	"proprietario":         "Proprietário",
	"gerente_agricola":     "Gerente Agrícola",
	"engenheiro_agronomo":  "Engenheiro Agrônomo",
	"tecnico_agricola":     "Técnico Agrícola",
	"operador_campo":       "Operador de Campo",
	"financeiro":           "Financeiro",
	"consultor_externo":    "Consultor Externo",
	"auditor":              "Auditor",
}

// Role is a papel (persona) a User can hold. System roles
// (platform_owner, organization_admin) have OrganizationID == nil, are
// shared by every organization, and cannot be edited or deleted. All other
// roles belong to a single organization and are managed by that
// organization's admins through the Roles screen.
type Role struct {
	ID             string    `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	OrganizationID *string   `json:"organization_id" gorm:"type:uuid"`
	Key            string    `json:"key" gorm:"not null"`
	Name           string    `json:"name" gorm:"not null"`
	IsSystem       bool      `json:"is_system" gorm:"not null;default:false"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
