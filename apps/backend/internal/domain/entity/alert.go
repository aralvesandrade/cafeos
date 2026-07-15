package entity

import "time"

type Alert struct {
	ID             string       `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	OrganizationID string       `json:"organization_id" gorm:"type:uuid;not null;index"`
	HarvestID      string       `json:"harvest_id" gorm:"type:uuid;not null;index"`
	RuleID         string       `json:"rule_id" gorm:"not null"`
	Message        string       `json:"message" gorm:"not null"`
	Severity       string       `json:"severity" gorm:"default:'warning'"`
	Status         string       `json:"status" gorm:"default:'aberto';index"`
	CreatedAt      time.Time    `json:"created_at"`
	ResolvedAt     *time.Time   `json:"resolved_at"`
	Organization   Organization `json:"-" gorm:"foreignKey:OrganizationID"`
	Harvest        Harvest      `json:"-" gorm:"foreignKey:HarvestID"`
}

func (Alert) TableName() string {
	return "alerts"
}
