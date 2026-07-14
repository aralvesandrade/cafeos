package entity

import "time"

type StockItem struct {
	ID             string              `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	OrganizationID string              `json:"organization_id" gorm:"type:uuid;not null;index"`
	ProductID      string              `json:"product_id" gorm:"type:uuid;not null;index"`
	Quantity       float64             `json:"quantity" gorm:"type:numeric(12,2);default:0"`
	Unit           string              `json:"unit" gorm:"default:''"`
	Batch          string              `json:"batch" gorm:"default:''"`
	ExpiryDate     *time.Time          `json:"expiry_date"`
	MinStock       float64             `json:"min_stock" gorm:"type:numeric(12,2);default:0"`
	Location       string              `json:"location" gorm:"default:''"`
	Notes          string              `json:"notes" gorm:"default:''"`
	CreatedAt      time.Time           `json:"created_at"`
	UpdatedAt      time.Time           `json:"updated_at"`
	Organization   Organization        `json:"-" gorm:"foreignKey:OrganizationID"`
	Product        AgriculturalProduct `json:"-" gorm:"foreignKey:ProductID"`
}

type StockMovement struct {
	ID             string       `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	OrganizationID string       `json:"organization_id" gorm:"type:uuid;not null;index"`
	ItemID         string       `json:"item_id" gorm:"type:uuid;not null;index"`
	Type           string       `json:"type" gorm:"not null"` // "in" | "out"
	Quantity       float64      `json:"quantity" gorm:"type:numeric(12,2);not null"`
	Date           time.Time    `json:"date" gorm:"not null;index"`
	Reference      string       `json:"reference" gorm:"default:''"`
	Notes          string       `json:"notes" gorm:"default:''"`
	CreatedAt      time.Time    `json:"created_at"`
	Organization   Organization `json:"-" gorm:"foreignKey:OrganizationID"`
	Item           StockItem    `json:"-" gorm:"foreignKey:ItemID"`
}
