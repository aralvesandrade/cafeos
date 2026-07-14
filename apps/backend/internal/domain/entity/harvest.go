package entity

import "time"

type HarvestStatus string

const (
	HarvestPlanejada   HarvestStatus = "planejada"
	HarvestEmAndamento HarvestStatus = "em_andamento"
	HarvestFinalizada  HarvestStatus = "finalizada"
)

type Harvest struct {
	ID                  string        `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	OrganizationID      string        `json:"organization_id" gorm:"type:uuid;not null;uniqueIndex:idx_organization_year"`
	Year                int           `json:"year" gorm:"not null;uniqueIndex:idx_organization_year"`
	Description         string        `json:"description" gorm:"default:''"`
	EstimatedProduction float64       `json:"estimated_production" gorm:"type:numeric(12,2);default:0"`
	Status              HarvestStatus `json:"status" gorm:"default:'planejada'"`
	StartDate           *time.Time    `json:"start_date"`
	EndDate             *time.Time    `json:"end_date"`
	CreatedAt           time.Time     `json:"created_at"`
	UpdatedAt           time.Time     `json:"updated_at"`
	Organization        Organization  `json:"-" gorm:"foreignKey:OrganizationID"`
}
