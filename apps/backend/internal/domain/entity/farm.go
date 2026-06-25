package entity

import "time"

type Farm struct {
	ID           string    `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	TenantID     string    `json:"tenant_id" gorm:"type:uuid;not null;index"`
	Name         string    `json:"name" gorm:"not null"`
	Owner        string    `json:"owner" gorm:"default:''"`
	Location     string    `json:"location" gorm:"default:''"`
	TotalAreaHA  float64   `json:"total_area_ha" gorm:"type:numeric(12,2);default:0"`
	PlantedAreaHA float64  `json:"planted_area_ha" gorm:"type:numeric(12,2);default:0"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Tenant       Tenant    `json:"-" gorm:"foreignKey:TenantID"`
	Plots        []Plot    `json:"plots,omitempty" gorm:"foreignKey:FarmID"`
}
