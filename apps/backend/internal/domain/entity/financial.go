package entity

import "time"

type TransactionType string

const (
	TranReceita TransactionType = "receita"
	TranDespesa TransactionType = "despesa"
)

type FinancialTransaction struct {
	ID            string          `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	TenantID      string          `json:"tenant_id" gorm:"type:uuid;not null;index"`
	Type          TransactionType `json:"type" gorm:"not null"`
	Category      string          `json:"category" gorm:"default:''"`
	Description   string          `json:"description" gorm:"not null"`
	Amount        float64         `json:"amount" gorm:"type:numeric(12,2);not null"`
	Date          time.Time       `json:"date" gorm:"not null;index"`
	DueDate       time.Time       `json:"due_date"`
	PaymentDate   *time.Time      `json:"payment_date"`
	Status        string          `json:"status" gorm:"default:'pending'"`
	PaymentMethod string          `json:"payment_method" gorm:"default:''"`
	Notes         string          `json:"notes" gorm:"default:''"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
	Tenant        Tenant          `json:"-" gorm:"foreignKey:TenantID"`
}
