package entity

import "time"

type Tenant struct {
	ID           string    `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name         string    `json:"name" gorm:"not null"`
	Slug         string    `json:"slug" gorm:"uniqueIndex;not null"`
	BrandName    string    `json:"brand_name" gorm:"default:''"`
	LogoURL      string    `json:"logo_url" gorm:"default:''"`
	PrimaryColor string    `json:"primary_color" gorm:"default:'#2E7D32'"`
	Plan         string    `json:"plan" gorm:"default:'free'"`
	Domain       string    `json:"domain" gorm:"default:''"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
