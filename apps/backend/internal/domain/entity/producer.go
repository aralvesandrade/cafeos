package entity

import "time"

// Producer is a Farm↔User link with a role scoped to that farm (e.g. one
// farm can have a "proprietario", an "operador_campo" and a "financeiro",
// each a real registered User). It is NOT organization↔user membership —
// User.RoleID (global, system-wide) is untouched by this; RoleID here is
// purely descriptive of the user's function on this specific farm.
//
// The legal/registration fields (CPF, RG, etc.) remain for compliance
// purposes (CAR/INCRA registration data) — they're only meaningful when
// the link's role is "proprietario", but nothing enforces that at the
// database level; it's a UI convention.
type Producer struct {
	ID             string `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	OrganizationID string `json:"organization_id" gorm:"type:uuid;not null;index"`
	FarmID         string `json:"farm_id" gorm:"type:uuid;not null;uniqueIndex:idx_producers_farm_user"`
	UserID         string `json:"user_id" gorm:"type:uuid;not null;uniqueIndex:idx_producers_farm_user"`
	RoleID         string `json:"role_id" gorm:"type:uuid;not null"`

	CPF           string     `json:"cpf" gorm:"default:''"`
	Name          string     `json:"name" gorm:"not null"`
	RG            string     `json:"rg" gorm:"default:''"`
	IssuingBody   string     `json:"issuing_body" gorm:"default:''"`
	Gender        string     `json:"gender" gorm:"default:''"`
	BirthDate     *time.Time `json:"birth_date"`
	MaritalStatus string     `json:"marital_status" gorm:"default:''"`
	Phone         string     `json:"phone" gorm:"default:''"`
	Email         string     `json:"email" gorm:"default:''"`
	Education     string     `json:"education" gorm:"default:''"`

	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
	Organization Organization `json:"-" gorm:"foreignKey:OrganizationID"`
	Farm         Farm         `json:"-" gorm:"foreignKey:FarmID"`
	User         User         `json:"-" gorm:"foreignKey:UserID"`
	Role         Role         `json:"-" gorm:"foreignKey:RoleID"`
}

func (Producer) TableName() string {
	return "producers"
}
