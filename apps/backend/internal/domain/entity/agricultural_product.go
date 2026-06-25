package entity

import "time"

type ProductType string

const (
	ProdFertilizante ProductType = "fertilizante"
	ProdDefensivo    ProductType = "defensivo"
	ProdCombustivel  ProductType = "combustivel"
	ProdOutro        ProductType = "outro"
)

type AgriculturalProduct struct {
	ID        string      `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	TenantID  string      `json:"tenant_id" gorm:"type:uuid;not null;index"`
	Name      string      `json:"name" gorm:"not null"`
	Type      ProductType `json:"type" gorm:"default:'outro'"`
	Unit      string      `json:"unit" gorm:"default:''"`
	CreatedAt time.Time   `json:"created_at"`
	Tenant    Tenant      `json:"-" gorm:"foreignKey:TenantID"`
}
