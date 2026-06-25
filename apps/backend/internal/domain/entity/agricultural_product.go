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
	ID        string      `json:"id" db:"id"`
	TenantID  string      `json:"tenant_id" db:"tenant_id"`
	Name      string      `json:"name" db:"name"`
	Type      ProductType `json:"type" db:"type"`
	Unit      string      `json:"unit" db:"unit"`
	CreatedAt time.Time   `json:"created_at" db:"created_at"`
}
