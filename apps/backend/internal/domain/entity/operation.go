package entity

import "time"

type OperationType string

const (
	OpAdubacao    OperationType = "adubacao"
	OpPulverizacao OperationType = "pulverizacao"
	OpIrrigacao   OperationType = "irrigacao"
	OpPoda        OperationType = "poda"
	OpColheita    OperationType = "colheita"
)

type Operation struct {
	ID            string        `json:"id" db:"id"`
	TenantID      string        `json:"tenant_id" db:"tenant_id"`
	PlotID        string        `json:"plot_id" db:"plot_id"`
	Type          OperationType `json:"type" db:"type"`
	Date          time.Time     `json:"date" db:"date"`
	Responsible   string        `json:"responsible" db:"responsible"`
	ProductUsed   string        `json:"product_used" db:"product_used"`
	Quantity      float64       `json:"quantity" db:"quantity"`
	Cost          float64       `json:"cost" db:"cost"`
	Notes         string        `json:"notes" db:"notes"`
	CreatedAt     time.Time     `json:"created_at" db:"created_at"`
}
