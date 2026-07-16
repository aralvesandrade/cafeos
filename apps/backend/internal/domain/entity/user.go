package entity

import "time"

type User struct {
	ID             string    `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	OrganizationID string    `json:"organization_id" gorm:"type:uuid;not null;uniqueIndex:idx_organization_email"`
	Name           string    `json:"name" gorm:"not null"`
	Email          string    `json:"email" gorm:"not null;uniqueIndex:idx_organization_email"`
	PasswordHash   string    `json:"-" gorm:"not null"`
	RoleID         string    `json:"role_id" gorm:"type:uuid;not null"`
	IsActive       bool      `json:"is_active" gorm:"default:true"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	PlanID *string `json:"plan_id" gorm:"type:uuid"`
	// ManagedByUserID nulo marca o usuário como "principal": dono de um
	// grupo de fazendas, responsável pelo próprio Plan, e o único que pode
	// criar sub-usuários (que recebem este campo apontando para ele).
	// Sub-usuário nunca cria outro usuário nem tem PlanID próprio.
	ManagedByUserID *string `json:"managed_by_user_id" gorm:"type:uuid;index"`

	Organization Organization `json:"-" gorm:"foreignKey:OrganizationID"`

	Role          Role  `json:"role" gorm:"foreignKey:RoleID"`
	Plan          *Plan `json:"plan,omitempty" gorm:"foreignKey:PlanID"`
	ManagedByUser *User `json:"managed_by_user,omitempty" gorm:"foreignKey:ManagedByUserID"`
}

func (User) TableName() string {
	return "users"
}
