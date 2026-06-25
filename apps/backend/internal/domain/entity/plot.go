package entity

import "time"

type Plot struct {
	ID           string    `json:"id" db:"id"`
	TenantID     string    `json:"tenant_id" db:"tenant_id"`
	FarmID       string    `json:"farm_id" db:"farm_id"`
	Name         string    `json:"name" db:"name"`
	AreaHA       float64   `json:"area_ha" db:"area_ha"`
	Cultivar     string    `json:"cultivar" db:"cultivar"`
	PlantingYear int       `json:"planting_year" db:"planting_year"`
	Altitude     int       `json:"altitude" db:"altitude"`
	SoilType     string    `json:"soil_type" db:"soil_type"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}
