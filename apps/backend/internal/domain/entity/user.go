package entity

import "time"

type UserRole string

const (
	RolePlatformOwner  UserRole = "platform_owner"
	RoleTenantAdmin    UserRole = "tenant_admin"
	RoleProprietario   UserRole = "proprietario"
	RoleGerente        UserRole = "gerente_agricola"
	RoleEngenheiro     UserRole = "engenheiro_agronomo"
	RoleTecnico        UserRole = "tecnico_agricola"
	RoleOperador       UserRole = "operador_campo"
	RoleFinanceiro     UserRole = "financeiro"
	RoleConsultor      UserRole = "consultor_externo"
	RoleAuditor        UserRole = "auditor"
)

type User struct {
	ID           string    `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	TenantID     string    `json:"tenant_id" gorm:"type:uuid;not null;uniqueIndex:idx_tenant_email"`
	Name         string    `json:"name" gorm:"not null"`
	Email        string    `json:"email" gorm:"not null;uniqueIndex:idx_tenant_email"`
	PasswordHash string    `json:"-" gorm:"not null"`
	Role         UserRole  `json:"role" gorm:"default:'operador_campo'"`
	IsActive     bool      `json:"is_active" gorm:"default:true"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Tenant       Tenant    `json:"-" gorm:"foreignKey:TenantID"`
}
