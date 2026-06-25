package event

import "log"

func SetupHandlers(bus Bus) {
	bus.Subscribe(func(event interface{}) {
		switch e := event.(type) {
		case OperationRegistered:
			log.Printf("[EVENT] OperationRegistered: %s - %s on plot %s (cost: %.2f)", e.OperationID, e.Type, e.PlotID, e.Cost)
		case HarvestFinalized:
			log.Printf("[EVENT] HarvestFinalized: %s - year %d", e.HarvestID, e.Year)
		case IndicatorUpdated:
			log.Printf("[EVENT] IndicatorUpdated: %s = %.2f for harvest %s", e.Type, e.Value, e.HarvestID)
		case AlertGenerated:
			log.Printf("[EVENT] AlertGenerated: [%s] %s - %s", e.Severity, e.RuleID, e.Message)
		}
	})
}
