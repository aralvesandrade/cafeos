package entity

import "time"

type CostCenterType string

const (
	CCReceita CostCenterType = "receita"
	CCDespesa CostCenterType = "despesa"
)

type CostCenter struct {
	ID          string         `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	TenantID    string         `json:"tenant_id" gorm:"type:uuid;not null;index"`
	Name        string         `json:"name" gorm:"not null"`
	Code        string         `json:"code" gorm:"uniqueIndex:idx_cost_centers_tenant_code;not null"`
	Type        CostCenterType `json:"type" gorm:"not null"`
	Description string         `json:"description" gorm:"default:''"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	Tenant      Tenant         `json:"-" gorm:"foreignKey:TenantID"`
}

func (CostCenter) TableName() string {
	return "cost_centers"
}
