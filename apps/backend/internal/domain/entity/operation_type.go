package entity

import "time"

type OperationType struct {
	ID             string       `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	OrganizationID string       `json:"organization_id" gorm:"type:uuid;not null;index"`
	Name           string       `json:"name" gorm:"not null"`
	Code           string       `json:"code" gorm:"uniqueIndex:idx_operation_types_organization_code;not null"`
	Color          string       `json:"color" gorm:"default:'default'"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
	Organization   Organization `json:"-" gorm:"foreignKey:OrganizationID"`
}

func (OperationType) TableName() string {
	return "operation_types"
}
