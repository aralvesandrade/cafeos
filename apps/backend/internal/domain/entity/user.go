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
	ID           string    `json:"id" db:"id"`
	TenantID     string    `json:"tenant_id" db:"tenant_id"`
	Name         string    `json:"name" db:"name"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	Role         UserRole  `json:"role" db:"role"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}
