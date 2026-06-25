package service

import (
	"fmt"
	"time"
)

type RuleCondition struct {
	Field    string
	Operator string
	Value    float64
}

type RuleAction struct {
	Type    string
	Message string
}

type Rule struct {
	ID          string
	Name        string
	Description string
	Conditions  []RuleCondition
	Action      RuleAction
	Enabled     bool
}

type RuleEngine struct {
	rules []Rule
}

func NewRuleEngine() *RuleEngine {
	return &RuleEngine{
		rules: defaultRules(),
	}
}

func defaultRules() []Rule {
	return []Rule{
		{
			ID:          "rule-low-productivity",
			Name:        "Alerta de Baixa Produtividade",
			Description: "Gera alerta se produtividade < 25 sacas/ha",
			Conditions: []RuleCondition{
				{Field: "sacas_por_hectare", Operator: "lt", Value: 25},
			},
			Action: RuleAction{
				Type:    "alert",
				Message: "Alerta de baixa produtividade: produtividade inferior a 25 sacas/hectare. Avaliar necessidade de correção de manejo.",
			},
			Enabled: true,
		},
		{
			ID:          "rule-high-cost",
			Name:        "Alerta de Custo Elevado",
			Description: "Gera alerta se custo por saca > 400",
			Conditions: []RuleCondition{
				{Field: "custo_por_saca", Operator: "gt", Value: 400},
			},
			Action: RuleAction{
				Type:    "alert",
				Message: "Alerta de custo elevado: custo por saca superior a R$ 400. Revisar despesas operacionais.",
			},
			Enabled: true,
		},
	}
}

type RuleEvaluation struct {
	Rule      Rule
	Triggered bool
	Alert     *RuleAlert
}

type RuleAlert struct {
	RuleID    string
	Message   string
	Severity  string
	CreatedAt time.Time
}

func (e *RuleEngine) Evaluate(field string, value float64) []RuleEvaluation {
	var results []RuleEvaluation

	for _, rule := range e.rules {
		if !rule.Enabled {
			continue
		}

		matched := false
		for _, cond := range rule.Conditions {
			if cond.Field != field {
				continue
			}
			switch cond.Operator {
			case "lt":
				matched = value < cond.Value
			case "gt":
				matched = value > cond.Value
			case "lte":
				matched = value <= cond.Value
			case "gte":
				matched = value >= cond.Value
			case "eq":
				matched = value == cond.Value
			}
		}

		eval := RuleEvaluation{
			Rule:      rule,
			Triggered: matched,
		}

		if matched {
			eval.Alert = &RuleAlert{
				RuleID:    rule.ID,
				Message:   rule.Action.Message,
				Severity:  "warning",
				CreatedAt: time.Now(),
			}
		}

		results = append(results, eval)
	}

	return results
}

func (e *RuleEngine) AddRule(rule Rule) {
	e.rules = append(e.rules, rule)
}

func (e *RuleEngine) ListRules() []Rule {
	return e.rules
}

func (e *RuleEngine) EvaluateIndicators(indicators map[string]float64) []RuleEvaluation {
	var allResults []RuleEvaluation

	for field, value := range indicators {
		results := e.Evaluate(field, value)
		allResults = append(allResults, results...)
	}

	return allResults
}

func (e *RuleEngine) FormatAlert(alert *RuleAlert) string {
	return fmt.Sprintf("[%s] %s - %s", alert.Severity, alert.RuleID, alert.Message)
}
