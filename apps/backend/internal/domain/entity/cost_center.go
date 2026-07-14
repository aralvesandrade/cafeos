package entity

import "time"

type CostCenterType string

const (
	CCReceita CostCenterType = "receita"
	CCDespesa CostCenterType = "despesa"
)

// CostGroup classifies a despesa CostCenter into the SENAR/CEPEA cost
// hierarchy, used to compute COE (Custo Operacional Efetivo), COT (Custo
// Operacional Total) and CT (Custo Total):
//
//	COE = sum(OperacionalEfetivo)
//	COT = COE + sum(MaoDeObraFamiliar) + sum(CapitalDepreciacao)
//	CT  = COT + sum(RemuneracaoCapital)
//
// A CostCenter with an empty CostGroup is not counted in any of the three
// (e.g. legacy cost centers created before this classification existed).
type CostGroup string

const (
	CostGroupOperacionalEfetivo CostGroup = "operacional_efetivo"
	CostGroupMaoDeObraFamiliar  CostGroup = "mao_de_obra_familiar"
	CostGroupCapitalDepreciacao CostGroup = "capital_depreciacao"
	CostGroupRemuneracaoCapital CostGroup = "remuneracao_capital"
)

type CostCenter struct {
	ID             string         `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	OrganizationID string         `json:"organization_id" gorm:"type:uuid;not null;index"`
	Name           string         `json:"name" gorm:"not null"`
	Code           string         `json:"code" gorm:"uniqueIndex:idx_cost_centers_organization_code;not null"`
	Type           CostCenterType `json:"type" gorm:"not null"`
	CostGroup      CostGroup      `json:"cost_group" gorm:"default:''"`
	Description    string         `json:"description" gorm:"default:''"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	Organization   Organization   `json:"-" gorm:"foreignKey:OrganizationID"`
}

func (CostCenter) TableName() string {
	return "cost_centers"
}
