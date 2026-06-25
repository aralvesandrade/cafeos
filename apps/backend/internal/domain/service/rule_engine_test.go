package service

import (
	"testing"
)

func TestRuleEngine_Evaluate(t *testing.T) {
	engine := NewRuleEngine()

	results := engine.Evaluate("sacas_por_hectare", 20)
	triggered := false
	for _, r := range results {
		if r.Rule.ID == "rule-low-productivity" && r.Triggered {
			triggered = true
			break
		}
	}
	if !triggered {
		t.Error("expected low productivity rule to trigger for value 20")
	}

	results = engine.Evaluate("sacas_por_hectare", 30)
	for _, r := range results {
		if r.Rule.ID == "rule-low-productivity" && r.Triggered {
			t.Error("expected low productivity rule NOT to trigger for value 30")
		}
	}
}

func TestRuleEngine_HighCostAlert(t *testing.T) {
	engine := NewRuleEngine()

	results := engine.Evaluate("custo_por_saca", 500)
	triggered := false
	for _, r := range results {
		if r.Rule.ID == "rule-high-cost" && r.Triggered {
			triggered = true
			break
		}
	}
	if !triggered {
		t.Error("expected high cost rule to trigger for value 500")
	}
}

func TestRuleEngine_DisabledRule(t *testing.T) {
	engine := NewRuleEngine()
	engine.AddRule(Rule{
		ID:      "test-disabled",
		Name:    "Disabled Rule",
		Enabled: false,
		Conditions: []RuleCondition{
			{Field: "test", Operator: "gt", Value: 0},
		},
		Action: RuleAction{Type: "alert", Message: "test"},
	})

	results := engine.Evaluate("test", 100)
	for _, r := range results {
		if r.Rule.ID == "test-disabled" && r.Triggered {
			t.Error("disabled rule should not trigger")
		}
	}
}

func TestRuleEngine_EvaluateIndicators(t *testing.T) {
	engine := NewRuleEngine()

	indicators := map[string]float64{
		"sacas_por_hectare": 20,
		"custo_por_saca":    500,
	}

	results := engine.EvaluateIndicators(indicators)
	if len(results) == 0 {
		t.Error("expected results from indicator evaluation")
	}

	alertCount := 0
	for _, r := range results {
		if r.Triggered {
			alertCount++
		}
	}
	if alertCount != 2 {
		t.Errorf("expected 2 triggered alerts, got %d", alertCount)
	}
}

func TestRuleEngine_FormatAlert(t *testing.T) {
	engine := NewRuleEngine()
	alert := &RuleAlert{
		RuleID:   "rule-test",
		Message:  "Test message",
		Severity: "warning",
	}

	formatted := engine.FormatAlert(alert)
	if formatted == "" {
		t.Error("expected non-empty formatted alert")
	}
}

func TestRuleEngine_AddRule(t *testing.T) {
	engine := NewRuleEngine()
	initialCount := len(engine.ListRules())

	engine.AddRule(Rule{
		ID:      "custom-rule",
		Name:    "Custom Rule",
		Enabled: true,
		Conditions: []RuleCondition{
			{Field: "test", Operator: "gt", Value: 50},
		},
		Action: RuleAction{Type: "alert", Message: "Custom alert"},
	})

	if len(engine.ListRules()) != initialCount+1 {
		t.Errorf("expected %d rules, got %d", initialCount+1, len(engine.ListRules()))
	}

	results := engine.Evaluate("test", 100)
	found := false
	for _, r := range results {
		if r.Rule.ID == "custom-rule" && r.Triggered {
			found = true
			break
		}
	}
	if !found {
		t.Error("custom rule should have triggered")
	}
}

func TestRuleEngine_Operators(t *testing.T) {
	engine := NewRuleEngine()

	engine.AddRule(Rule{
		ID: "test-lte", Enabled: true,
		Conditions: []RuleCondition{{Field: "val", Operator: "lte", Value: 10}},
		Action:     RuleAction{Type: "alert", Message: "lte"},
	})
	engine.AddRule(Rule{
		ID: "test-gte", Enabled: true,
		Conditions: []RuleCondition{{Field: "val", Operator: "gte", Value: 10}},
		Action:     RuleAction{Type: "alert", Message: "gte"},
	})
	engine.AddRule(Rule{
		ID: "test-eq", Enabled: true,
		Conditions: []RuleCondition{{Field: "val", Operator: "eq", Value: 10}},
		Action:     RuleAction{Type: "alert", Message: "eq"},
	})

	tests := []struct {
		value    float64
		lteMatch bool
		gteMatch bool
		eqMatch  bool
	}{
		{5, true, false, false},
		{10, true, true, true},
		{15, false, true, false},
	}

	for _, tt := range tests {
		results := engine.Evaluate("val", tt.value)
		for _, r := range results {
			switch r.Rule.ID {
			case "test-lte":
				if r.Triggered != tt.lteMatch {
					t.Errorf("lte(%f): expected %v, got %v", tt.value, tt.lteMatch, r.Triggered)
				}
			case "test-gte":
				if r.Triggered != tt.gteMatch {
					t.Errorf("gte(%f): expected %v, got %v", tt.value, tt.gteMatch, r.Triggered)
				}
			case "test-eq":
				if r.Triggered != tt.eqMatch {
					t.Errorf("eq(%f): expected %v, got %v", tt.value, tt.eqMatch, r.Triggered)
				}
			}
		}
	}
}
