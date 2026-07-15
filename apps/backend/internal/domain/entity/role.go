package entity

import "time"

// SystemRole keys are not editable or deletable — platform_owner is checked
// directly by RequireRole for /admin routes, and organization_admin is the
// highest-privilege role within a tenant. Every other role is an ordinary,
// editable row in the roles table.
const (
	SystemRolePlatformOwner     = "platform_owner"
	SystemRoleOrganizationAdmin = "organization_admin"
)

var SystemRoleKeys = []string{SystemRolePlatformOwner, SystemRoleOrganizationAdmin}

// RoleKeyProprietario is referenced by name in ownership-scoping logic (e.g.
// restrictedOwnerID in api/handler/access.go), since "proprietario" is the
// only role whose visibility is scoped to farms it owns rather than the
// whole organization. It may still be renamed or deleted like any other
// non-system role — the scoping check simply no-ops if it's gone.
const RoleKeyProprietario = "proprietario"

// DefaultRoleKeys/DefaultRoleNames seed the roles table on first boot,
// alongside the two system roles above.
var DefaultRoleKeys = []string{
	RoleKeyProprietario,
	"gerente_agricola",
	"engenheiro_agronomo",
	"tecnico_agricola",
	"operador_campo",
	"financeiro",
	"consultor_externo",
	"auditor",
}

var DefaultRoleNames = map[string]string{
	"proprietario":        "Proprietário",
	"gerente_agricola":    "Gerente Agrícola",
	"engenheiro_agronomo": "Engenheiro Agrônomo",
	"tecnico_agricola":    "Técnico Agrícola",
	"operador_campo":      "Operador de Campo",
	"financeiro":          "Financeiro",
	"consultor_externo":   "Consultor Externo",
	"auditor":             "Auditor",
}

// Role is a papel (persona) a User can hold. Roles are global — shared by
// every organization, like Module — not organization-scoped. What varies
// per organization is access (see RolePermission), not the set of roles
// available. System roles (platform_owner, organization_admin) cannot be
// edited or deleted; every other role can be renamed, created, or removed
// (if unused) through the Roles screen.
type Role struct {
	ID        string    `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Key       string    `json:"key" gorm:"not null;uniqueIndex"`
	Name      string    `json:"name" gorm:"not null"`
	IsSystem  bool      `json:"is_system" gorm:"not null;default:false"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
