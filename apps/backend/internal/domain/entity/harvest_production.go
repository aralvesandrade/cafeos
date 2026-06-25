package entity

import "time"

type HarvestProduction struct {
	ID         string    `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	TenantID   string    `json:"tenant_id" gorm:"type:uuid;not null;index"`
	HarvestID  string    `json:"harvest_id" gorm:"type:uuid;not null;index"`
	PlotID     string    `json:"plot_id" gorm:"type:uuid;not null;index"`
	Quantity   float64   `json:"quantity" gorm:"type:numeric(12,2);default:0"`
	RecordedAt time.Time `json:"recorded_at"`
	Notes      string    `json:"notes" gorm:"default:''"`
	Tenant     Tenant   `json:"-" gorm:"foreignKey:TenantID"`
	Harvest    Harvest   `json:"-" gorm:"foreignKey:HarvestID"`
	Plot       Plot      `json:"-" gorm:"foreignKey:PlotID"`
}
