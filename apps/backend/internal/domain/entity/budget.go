package entity

import "time"

type Budget struct {
	ID            string      `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	TenantID      string      `json:"tenant_id" gorm:"type:uuid;not null;uniqueIndex:idx_budget_harvest_cc"`
	HarvestID     string      `json:"harvest_id" gorm:"type:uuid;not null;uniqueIndex:idx_budget_harvest_cc"`
	CostCenterID  string      `json:"cost_center_id" gorm:"type:uuid;not null;uniqueIndex:idx_budget_harvest_cc"`
	PlannedAmount float64     `json:"planned_amount" gorm:"type:numeric(12,2);default:0"`
	Description   string      `json:"description" gorm:"default:''"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
	Tenant        Tenant      `json:"-" gorm:"foreignKey:TenantID"`
	Harvest       Harvest     `json:"-" gorm:"foreignKey:HarvestID"`
	CostCenter    CostCenter  `json:"-" gorm:"foreignKey:CostCenterID"`
}

func (Budget) TableName() string {
	return "budgets"
}
