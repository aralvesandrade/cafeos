package entity

import "time"

type Team struct {
	ID             string       `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	OrganizationID string       `json:"organization_id" gorm:"type:uuid;not null;index"`
	Name           string       `json:"name" gorm:"not null"`
	Leader         string       `json:"leader" gorm:"default:''"`
	Description    string       `json:"description" gorm:"default:''"`
	CreatedAt      time.Time    `json:"created_at"`
	Organization   Organization `json:"-" gorm:"foreignKey:OrganizationID"`
}

type Worker struct {
	ID             string       `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	OrganizationID string       `json:"organization_id" gorm:"type:uuid;not null;index"`
	TeamID         string       `json:"team_id" gorm:"type:uuid;index"`
	Name           string       `json:"name" gorm:"not null"`
	Role           string       `json:"role" gorm:"default:''"`
	Phone          string       `json:"phone" gorm:"default:''"`
	HourlyRate     float64      `json:"hourly_rate" gorm:"type:numeric(10,2);default:0"`
	IsActive       bool         `json:"is_active" gorm:"default:true"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
	Organization   Organization `json:"-" gorm:"foreignKey:OrganizationID"`
	Team           Team         `json:"-" gorm:"foreignKey:TeamID"`
}

type WorkShift struct {
	ID             string       `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	OrganizationID string       `json:"organization_id" gorm:"type:uuid;not null;index"`
	WorkerID       string       `json:"worker_id" gorm:"type:uuid;not null;index"`
	OperationID    *string      `json:"operation_id" gorm:"type:uuid;index"`
	CostCenterID   *string      `json:"cost_center_id" gorm:"type:uuid;index"`
	Date           time.Time    `json:"date" gorm:"not null;index"`
	Hours          float64      `json:"hours" gorm:"type:numeric(5,2);not null"`
	Cost           float64      `json:"cost" gorm:"type:numeric(10,2);default:0"`
	Notes          string       `json:"notes" gorm:"default:''"`
	CreatedAt      time.Time    `json:"created_at"`
	Organization   Organization `json:"-" gorm:"foreignKey:OrganizationID"`
	Worker         Worker       `json:"-" gorm:"foreignKey:WorkerID"`
	CostCenter     *CostCenter  `json:"-" gorm:"foreignKey:CostCenterID"`
}
