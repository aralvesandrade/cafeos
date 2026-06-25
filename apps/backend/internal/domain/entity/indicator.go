package entity

import "time"

type IndicatorType string

const (
	IndSacasHA     IndicatorType = "sacas_por_hectare"
	IndCustoSaca   IndicatorType = "custo_por_saca"
	IndRentabilidade IndicatorType = "rentabilidade"
	IndBienalidade IndicatorType = "bienalidade"
	IndProducaoTotal IndicatorType = "producao_total"
	IndCustoTotal  IndicatorType = "custo_total"
)

type Indicator struct {
	ID         string        `json:"id" db:"id"`
	TenantID   string        `json:"tenant_id" db:"tenant_id"`
	HarvestID  string        `json:"harvest_id" db:"harvest_id"`
	PlotID     *string       `json:"plot_id,omitempty" db:"plot_id"`
	Type       IndicatorType `json:"type" db:"type"`
	Value      float64       `json:"value" db:"value"`
	CalculatedAt time.Time   `json:"calculated_at" db:"calculated_at"`
}
