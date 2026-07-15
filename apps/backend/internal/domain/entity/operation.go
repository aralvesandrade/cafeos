package entity

import "time"

type Operation struct {
	ID             string         `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	OrganizationID string         `json:"organization_id" gorm:"type:uuid;not null;index"`
	PlotID         string         `json:"plot_id" gorm:"type:uuid;not null;index"`
	PlotName       string         `json:"plot_name" gorm:"->;-:migration"`
	HarvestID      *string        `json:"harvest_id" gorm:"type:uuid;index"`
	CostCenterID   *string        `json:"cost_center_id" gorm:"type:uuid;index"`
	TypeID         string         `json:"type_id" gorm:"type:uuid;not null;index"`
	TypeName       string         `json:"type_name" gorm:"->;-:migration"`
	TypeColor      string         `json:"type_color" gorm:"->;-:migration"`
	Date           time.Time      `json:"date" gorm:"not null;index"`
	Responsible    string         `json:"responsible" gorm:"default:''"`
	ProductUsed    string         `json:"product_used" gorm:"default:''"`
	Quantity       float64        `json:"quantity" gorm:"type:numeric(12,2);default:0"`
	Cost           float64        `json:"cost" gorm:"type:numeric(12,2);default:0"`
	Notes          string         `json:"notes" gorm:"default:''"`
	CreatedAt      time.Time      `json:"created_at"`
	Organization   Organization   `json:"-" gorm:"foreignKey:OrganizationID"`
	Plot           Plot           `json:"-" gorm:"foreignKey:PlotID"`
	Harvest        *Harvest       `json:"-" gorm:"foreignKey:HarvestID"`
	CostCenter     *CostCenter    `json:"-" gorm:"foreignKey:CostCenterID"`
	Type           *OperationType `json:"-" gorm:"foreignKey:TypeID"`
}

func (Operation) TableName() string {
	return "operations"
}
