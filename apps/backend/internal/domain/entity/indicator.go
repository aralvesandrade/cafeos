package entity

import "time"

type IndicatorType string

const (
	IndSacasHA       IndicatorType = "sacas_por_hectare"
	IndCustoSaca     IndicatorType = "custo_por_saca"
	IndProducaoTotal IndicatorType = "producao_total"
	IndCustoTotal    IndicatorType = "custo_total"

	// SENAR/CEPEA cost hierarchy indicators. These require CostCenter.CostGroup
	// classification (see entity.CostGroup) — cost centers without a
	// classification are not counted in any of the three.
	IndAreaProducao IndicatorType = "area_producao"

	IndCOE        IndicatorType = "coe"
	IndCOEPorArea IndicatorType = "coe_por_area"
	IndCOEPorSaca IndicatorType = "coe_por_saca"

	IndCOT        IndicatorType = "cot"
	IndCOTPorArea IndicatorType = "cot_por_area"
	IndCOTPorSaca IndicatorType = "cot_por_saca"

	IndCTProducao        IndicatorType = "ct_producao"
	IndCTProducaoPorArea IndicatorType = "ct_producao_por_area"
	IndCTProducaoPorSaca IndicatorType = "ct_producao_por_saca"
)

type Indicator struct {
	ID             string        `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	OrganizationID string        `json:"organization_id" gorm:"type:uuid;not null;index"`
	HarvestID      string        `json:"harvest_id" gorm:"type:uuid;not null;index"`
	PlotID         *string       `json:"plot_id,omitempty" gorm:"type:uuid"`
	Type           IndicatorType `json:"type" gorm:"not null;index"`
	Value          float64       `json:"value" gorm:"type:numeric(14,2);default:0"`
	CalculatedAt   time.Time     `json:"calculated_at"`
	Organization   Organization  `json:"-" gorm:"foreignKey:OrganizationID"`
	Harvest        Harvest       `json:"-" gorm:"foreignKey:HarvestID"`
}

func (Indicator) TableName() string {
	return "indicators"
}
