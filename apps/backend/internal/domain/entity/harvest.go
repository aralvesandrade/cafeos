package entity

import "time"

type HarvestStatus string

const (
	HarvestPlanejada   HarvestStatus = "planejada"
	HarvestEmAndamento HarvestStatus = "em_andamento"
	HarvestFinalizada  HarvestStatus = "finalizada"
)

type Harvest struct {
	ID                 string        `json:"id" db:"id"`
	TenantID           string        `json:"tenant_id" db:"tenant_id"`
	Year               int           `json:"year" db:"year"`
	Description        string        `json:"description" db:"description"`
	EstimatedProduction float64      `json:"estimated_production" db:"estimated_production"`
	Status             HarvestStatus `json:"status" db:"status"`
	CreatedAt          time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time     `json:"updated_at" db:"updated_at"`
}
