package service

import (
	"testing"

	"github.com/aralvesandrade/cafeos/internal/domain/entity"
)

func TestHarvestService_BuildIndicators_COE_COT_CT(t *testing.T) {
	s := &HarvestService{}
	harvest := &entity.Harvest{ID: "h1", TenantID: "t1"}

	costByGroup := map[entity.CostGroup]float64{
		entity.CostGroupOperacionalEfetivo: 1000,
		entity.CostGroupMaoDeObraFamiliar:  200,
		entity.CostGroupCapitalDepreciacao: 100,
		entity.CostGroupRemuneracaoCapital: 50,
	}

	indicators := s.buildIndicators(harvest, 100 /* totalProduction */, 10 /* totalPlantedArea */, 5000 /* totalCost legacy */, costByGroup)

	byType := make(map[entity.IndicatorType]float64)
	for _, ind := range indicators {
		byType[ind.Type] = ind.Value
	}

	if got := byType[entity.IndCOE]; got != 1000 {
		t.Errorf("expected COE 1000, got %f", got)
	}
	if got := byType[entity.IndCOT]; got != 1300 {
		t.Errorf("expected COT 1300 (COE + mao de obra familiar + capital), got %f", got)
	}
	if got := byType[entity.IndCTProducao]; got != 1350 {
		t.Errorf("expected CT 1350 (COT + remuneracao capital), got %f", got)
	}
	if got := byType[entity.IndCOEPorArea]; got != 100 {
		t.Errorf("expected COE/area 100, got %f", got)
	}
	if got := byType[entity.IndCOTPorSaca]; got != 13 {
		t.Errorf("expected COT/saca 13, got %f", got)
	}
	if got := byType[entity.IndAreaProducao]; got != 10 {
		t.Errorf("expected area_producao 10, got %f", got)
	}
	// legacy indicators must remain untouched by the new classification
	if got := byType[entity.IndCustoTotal]; got != 5000 {
		t.Errorf("expected legacy custo_total 5000, got %f", got)
	}
}

func TestHarvestService_BuildIndicators_UnclassifiedCostsExcluded(t *testing.T) {
	s := &HarvestService{}
	harvest := &entity.Harvest{ID: "h1", TenantID: "t1"}

	// no classified cost groups at all (legacy cost centers)
	indicators := s.buildIndicators(harvest, 100, 10, 5000, map[entity.CostGroup]float64{})

	for _, ind := range indicators {
		switch ind.Type {
		case entity.IndCOE, entity.IndCOT, entity.IndCTProducao:
			if ind.Value != 0 {
				t.Errorf("expected %s to be 0 when no cost center is classified, got %f", ind.Type, ind.Value)
			}
		}
	}
}
