package entity

import "time"

type Tenant struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Slug        string    `json:"slug" db:"slug"`
	BrandName   string    `json:"brand_name" db:"brand_name"`
	LogoURL     string    `json:"logo_url" db:"logo_url"`
	PrimaryColor string   `json:"primary_color" db:"primary_color"`
	Plan        string    `json:"plan" db:"plan"`
	Domain      string    `json:"domain" db:"domain"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
