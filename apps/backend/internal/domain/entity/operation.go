package entity

import "time"

type OperationType string

const (
	OpAdubacao     OperationType = "adubacao"
	OpPulverizacao OperationType = "pulverizacao"
	OpIrrigacao    OperationType = "irrigacao"
	OpPoda         OperationType = "poda"
	OpColheita     OperationType = "colheita"
)

type Operation struct {
	ID          string        `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	TenantID    string        `json:"tenant_id" gorm:"type:uuid;not null;index"`
	PlotID      string        `json:"plot_id" gorm:"type:uuid;not null;index"`
	PlotName    string        `json:"plot_name" gorm:"->;-:migration"`
	Type        OperationType `json:"type" gorm:"not null"`
	Date        time.Time     `json:"date" gorm:"not null;index"`
	Responsible string        `json:"responsible" gorm:"default:''"`
	ProductUsed string        `json:"product_used" gorm:"default:''"`
	Quantity    float64       `json:"quantity" gorm:"type:numeric(12,2);default:0"`
	Cost        float64       `json:"cost" gorm:"type:numeric(12,2);default:0"`
	Notes       string        `json:"notes" gorm:"default:''"`
	CreatedAt   time.Time     `json:"created_at"`
	Tenant      Tenant        `json:"-" gorm:"foreignKey:TenantID"`
	Plot        Plot          `json:"-" gorm:"foreignKey:PlotID"`
}
