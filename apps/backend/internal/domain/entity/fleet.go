package entity

import "time"

type VehicleType string

const (
	VeicTractor      VehicleType = "trator"
	VeicPulverizador VehicleType = "pulverizador"
	VeicCaminhao     VehicleType = "caminhao"
	VeicUtilitario   VehicleType = "utilitario"
	VeicOutro        VehicleType = "outro"
)

type Vehicle struct {
	ID        string      `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	TenantID  string      `json:"tenant_id" gorm:"type:uuid;not null;index"`
	Name      string      `json:"name" gorm:"not null"`
	Type      VehicleType `json:"type" gorm:"default:'outro'"`
	Plate     string      `json:"plate" gorm:"default:''"`
	Brand     string      `json:"brand" gorm:"default:''"`
	Model     string      `json:"model" gorm:"default:''"`
	Year      int         `json:"year" gorm:"default:0"`
	Status    string      `json:"status" gorm:"default:'active'"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	Tenant    Tenant      `json:"-" gorm:"foreignKey:TenantID"`
}

type Maintenance struct {
	ID           string      `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	TenantID     string      `json:"tenant_id" gorm:"type:uuid;not null;index"`
	VehicleID    string      `json:"vehicle_id" gorm:"type:uuid;not null;index"`
	CostCenterID *string     `json:"cost_center_id" gorm:"type:uuid;index"`
	Date         time.Time   `json:"date" gorm:"not null;index"`
	Type         string      `json:"type" gorm:"default:'preventive'"` // "preventive" | "corrective"
	Description  string      `json:"description" gorm:"default:''"`
	Cost         float64     `json:"cost" gorm:"type:numeric(12,2);default:0"`
	Odometer     float64     `json:"odometer" gorm:"type:numeric(12,2);default:0"`
	Notes        string      `json:"notes" gorm:"default:''"`
	CreatedAt    time.Time   `json:"created_at"`
	Tenant       Tenant      `json:"-" gorm:"foreignKey:TenantID"`
	Vehicle      Vehicle     `json:"-" gorm:"foreignKey:VehicleID"`
	CostCenter   *CostCenter `json:"-" gorm:"foreignKey:CostCenterID"`
}
