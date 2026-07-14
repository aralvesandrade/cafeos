package entity

import "time"

type AllocationMethod string

const (
	AllocAreaProportional AllocationMethod = "area_proportional"
	AllocCustomPercentage AllocationMethod = "custom_percentage"
)

type CostAllocation struct {
	ID             string               `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	OrganizationID string               `json:"organization_id" gorm:"type:uuid;not null;index"`
	HarvestID      string               `json:"harvest_id" gorm:"type:uuid;not null;index"`
	CostCenterID   string               `json:"cost_center_id" gorm:"type:uuid;not null;index"`
	Description    string               `json:"description" gorm:"not null"`
	TotalAmount    float64              `json:"total_amount" gorm:"type:numeric(12,2);not null"`
	Method         AllocationMethod     `json:"method" gorm:"not null"`
	Date           time.Time            `json:"date" gorm:"not null;index"`
	CreatedAt      time.Time            `json:"created_at"`
	Organization   Organization         `json:"-" gorm:"foreignKey:OrganizationID"`
	Harvest        Harvest              `json:"-" gorm:"foreignKey:HarvestID"`
	CostCenter     CostCenter           `json:"-" gorm:"foreignKey:CostCenterID"`
	Items          []CostAllocationItem `json:"items,omitempty" gorm:"foreignKey:AllocationID"`
}

func (CostAllocation) TableName() string {
	return "cost_allocations"
}

type CostAllocationItem struct {
	ID           string         `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	AllocationID string         `json:"allocation_id" gorm:"type:uuid;not null;index"`
	PlotID       string         `json:"plot_id" gorm:"type:uuid;not null;index"`
	Amount       float64        `json:"amount" gorm:"type:numeric(12,2);not null"`
	Percentage   float64        `json:"percentage" gorm:"type:numeric(5,2);default:0"`
	Allocation   CostAllocation `json:"-" gorm:"foreignKey:AllocationID"`
	Plot         Plot           `json:"-" gorm:"foreignKey:PlotID"`
}

func (CostAllocationItem) TableName() string {
	return "cost_allocation_items"
}
