package entity

import "time"

type User struct {
	ID             string       `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	OrganizationID string       `json:"organization_id" gorm:"type:uuid;not null;uniqueIndex:idx_organization_email"`
	Name           string       `json:"name" gorm:"not null"`
	Email          string       `json:"email" gorm:"not null;uniqueIndex:idx_organization_email"`
	PasswordHash   string       `json:"-" gorm:"not null"`
	RoleID         string       `json:"role_id" gorm:"type:uuid;not null"`
	IsActive       bool         `json:"is_active" gorm:"default:true"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
	Organization   Organization `json:"-" gorm:"foreignKey:OrganizationID"`
	Role           Role         `json:"role" gorm:"foreignKey:RoleID"`
}

func (User) TableName() string {
	return "users"
}
