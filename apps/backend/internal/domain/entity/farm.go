package entity

import "time"

type Farm struct {
	ID           string    `json:"id" db:"id"`
	TenantID     string    `json:"tenant_id" db:"tenant_id"`
	Name         string    `json:"name" db:"name"`
	Owner        string    `json:"owner" db:"owner"`
	Location     string    `json:"location" db:"location"`
	TotalAreaHA  float64   `json:"total_area_ha" db:"total_area_ha"`
	PlantedAreaHA float64  `json:"planted_area_ha" db:"planted_area_ha"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}
