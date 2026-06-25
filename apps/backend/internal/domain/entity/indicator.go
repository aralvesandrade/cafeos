package entity

import "time"

type IndicatorType string

const (
	IndSacasHA       IndicatorType = "sacas_por_hectare"
	IndCustoSaca     IndicatorType = "custo_por_saca"
	IndRentabilidade IndicatorType = "rentabilidade"
	IndBienalidade   IndicatorType = "bienalidade"
	IndProducaoTotal IndicatorType = "producao_total"
	IndCustoTotal    IndicatorType = "custo_total"
)

type Indicator struct {
	ID          string        `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	TenantID    string        `json:"tenant_id" gorm:"type:uuid;not null;index"`
	HarvestID   string        `json:"harvest_id" gorm:"type:uuid;not null;index"`
	PlotID      *string       `json:"plot_id,omitempty" gorm:"type:uuid"`
	Type        IndicatorType `json:"type" gorm:"not null;index"`
	Value       float64       `json:"value" gorm:"type:numeric(14,2);default:0"`
	CalculatedAt time.Time    `json:"calculated_at"`
	Tenant      Tenant        `json:"-" gorm:"foreignKey:TenantID"`
	Harvest     Harvest       `json:"-" gorm:"foreignKey:HarvestID"`
}
