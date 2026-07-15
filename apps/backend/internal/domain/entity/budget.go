package entity

import "time"

type Budget struct {
	ID             string       `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	OrganizationID string       `json:"organization_id" gorm:"type:uuid;not null;uniqueIndex:idx_budget_harvest_cc"`
	HarvestID      string       `json:"harvest_id" gorm:"type:uuid;not null;uniqueIndex:idx_budget_harvest_cc"`
	CostCenterID   string       `json:"cost_center_id" gorm:"type:uuid;not null;uniqueIndex:idx_budget_harvest_cc"`
	PlannedAmount  float64      `json:"planned_amount" gorm:"type:numeric(12,2);default:0"`
	Description    string       `json:"description" gorm:"default:''"`
	CostCenterName string       `json:"cost_center_name" gorm:"->;-:migration"`
	RealizedAmount float64      `json:"realized_amount" gorm:"->;-:migration"`
	Variance       float64      `json:"variance" gorm:"->;-:migration"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
	Organization   Organization `json:"-" gorm:"foreignKey:OrganizationID"`
	Harvest        Harvest      `json:"-" gorm:"foreignKey:HarvestID"`
	CostCenter     CostCenter   `json:"-" gorm:"foreignKey:CostCenterID"`
}

func (Budget) TableName() string {
	return "budgets"
}
