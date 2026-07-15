package entity

import "time"

// ModuleKey identifies a fixed area of the application. Route wiring in
// router.go references these constants directly — adding a module still
// requires code (new routes, new screen), so the set is closed. The
// modules table only stores display metadata (name, order) for these same
// keys, giving the admin UI a single source of truth instead of a second
// hardcoded label map.
type ModuleKey string

const (
	ModuleDashboard   ModuleKey = "dashboard"
	ModuleFarms       ModuleKey = "farms"
	ModuleOperations  ModuleKey = "operations"
	ModuleHarvests    ModuleKey = "harvests"
	ModuleResources   ModuleKey = "resources"
	ModuleFinancial   ModuleKey = "financial"
	ModuleUsers       ModuleKey = "users"
	ModulePermissions ModuleKey = "permissions"
)

var AllModules = []ModuleKey{
	ModuleDashboard,
	ModuleFarms,
	ModuleOperations,
	ModuleHarvests,
	ModuleResources,
	ModuleFinancial,
	ModuleUsers,
	ModulePermissions,
}

// DefaultModuleNames/DefaultModuleOrder seed the modules table. Order is a
// simple ascending sort key for the admin UI.
var DefaultModuleNames = map[ModuleKey]string{
	ModuleDashboard:   "Dashboard e Alertas",
	ModuleFarms:       "Fazendas, Talhões e Tipos de Operação",
	ModuleOperations:  "Operações",
	ModuleHarvests:    "Safras",
	ModuleResources:   "Estoque, Frota e Equipes",
	ModuleFinancial:   "Financeiro, Centros de Custo, Orçamento e Rateio",
	ModuleUsers:       "Usuários da Organização",
	ModulePermissions: "Permissões",
}

var DefaultModuleOrder = map[ModuleKey]int{
	ModuleDashboard:   1,
	ModuleFarms:       2,
	ModuleOperations:  3,
	ModuleHarvests:    4,
	ModuleResources:   5,
	ModuleFinancial:   6,
	ModuleUsers:       7,
	ModulePermissions: 8,
}

// Module holds the display metadata (name, order) for a ModuleKey. It is a
// global catalog — not organization-scoped — since a module is a property
// of the application, not of a tenant.
type Module struct {
	Key       ModuleKey `json:"key" gorm:"type:varchar(50);primaryKey"`
	Name      string    `json:"name" gorm:"not null"`
	Order     int       `json:"order" gorm:"not null;default:0"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
