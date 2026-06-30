package event

import "log/slog"

func SetupHandlers(bus Bus) {
	bus.Subscribe(func(event interface{}) {
		switch e := event.(type) {
		case OperationRegistered:
			slog.Info("event operation registered",
				"operation_id", e.OperationID,
				"type", e.Type,
				"plot_id", e.PlotID,
				"cost", e.Cost,
			)
		case HarvestFinalized:
			slog.Info("event harvest finalized",
				"harvest_id", e.HarvestID,
				"year", e.Year,
			)
		case IndicatorUpdated:
			slog.Info("event indicator updated",
				"type", e.Type,
				"value", e.Value,
				"harvest_id", e.HarvestID,
			)
		case AlertGenerated:
			slog.Info("event alert generated",
				"severity", e.Severity,
				"rule_id", e.RuleID,
				"message", e.Message,
			)
		}
	})
}
