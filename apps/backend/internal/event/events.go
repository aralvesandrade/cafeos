package event

import "time"

type OperationRegistered struct {
	OperationID    string
	OrganizationID string
	PlotID         string
	Type           string
	Cost           float64
	Date           time.Time
}

type HarvestFinalized struct {
	HarvestID      string
	OrganizationID string
	Year           int
}

type IndicatorUpdated struct {
	OrganizationID string
	HarvestID      string
	IndicatorID    string
	Type           string
	Value          float64
}

type AlertGenerated struct {
	OrganizationID string
	RuleID         string
	Message        string
	Severity       string
	CreatedAt      time.Time
}
