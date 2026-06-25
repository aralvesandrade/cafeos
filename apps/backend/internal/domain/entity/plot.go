package entity

import "time"

type Plot struct {
	ID           string    `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	TenantID     string    `json:"tenant_id" gorm:"type:uuid;not null;index"`
	FarmID       string    `json:"farm_id" gorm:"type:uuid;not null;index"`
	Name         string    `json:"name" gorm:"not null"`
	AreaHA       float64   `json:"area_ha" gorm:"type:numeric(12,2);default:0"`
	Cultivar     string    `json:"cultivar" gorm:"default:''"`
	PlantingYear int       `json:"planting_year" gorm:"default:0"`
	Altitude     int       `json:"altitude" gorm:"default:0"`
	SoilType     string    `json:"soil_type" gorm:"default:''"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Tenant       Tenant    `json:"-" gorm:"foreignKey:TenantID"`
	Farm         Farm      `json:"-" gorm:"foreignKey:FarmID"`
}
