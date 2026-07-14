package entity

import "time"

type HarvestProduction struct {
	ID             string       `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	OrganizationID string       `json:"organization_id" gorm:"type:uuid;not null;index"`
	HarvestID      string       `json:"harvest_id" gorm:"type:uuid;not null;index"`
	PlotID         string       `json:"plot_id" gorm:"type:uuid;not null;index"`
	Quantity       float64      `json:"quantity" gorm:"type:numeric(12,2);default:0"`
	RecordedAt     time.Time    `json:"recorded_at"`
	Notes          string       `json:"notes" gorm:"default:''"`
	Organization   Organization `json:"-" gorm:"foreignKey:OrganizationID"`
	Harvest        Harvest      `json:"-" gorm:"foreignKey:HarvestID"`
	Plot           Plot         `json:"-" gorm:"foreignKey:PlotID"`
}
