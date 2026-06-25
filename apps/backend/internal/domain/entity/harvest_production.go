package entity

import "time"

type HarvestProduction struct {
	ID         string    `json:"id" db:"id"`
	TenantID   string    `json:"tenant_id" db:"tenant_id"`
	HarvestID  string    `json:"harvest_id" db:"harvest_id"`
	PlotID     string    `json:"plot_id" db:"plot_id"`
	Quantity   float64   `json:"quantity" db:"quantity"`
	RecordedAt time.Time `json:"recorded_at" db:"recorded_at"`
	Notes      string    `json:"notes" db:"notes"`
}
