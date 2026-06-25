package event

import "time"

type OperationRegistered struct {
	OperationID string
	TenantID    string
	PlotID      string
	Type        string
	Cost        float64
	Date        time.Time
}

type HarvestFinalized struct {
	HarvestID string
	TenantID  string
	Year      int
}

type IndicatorUpdated struct {
	TenantID    string
	HarvestID   string
	IndicatorID string
	Type        string
	Value       float64
}

type AlertGenerated struct {
	TenantID  string
	RuleID    string
	Message   string
	Severity  string
	CreatedAt time.Time
}
