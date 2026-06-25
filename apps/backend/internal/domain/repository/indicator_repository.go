package repository

import "github.com/aralvesandrade/cafeos/internal/domain/entity"

type IndicatorRepository interface {
	Create(indicator *entity.Indicator) error
	ListByHarvest(harvestID string) ([]*entity.Indicator, error)
	ListByTenant(tenantID string) ([]*entity.Indicator, error)
	ListByTenantAndType(tenantID string, indicatorType entity.IndicatorType) ([]*entity.Indicator, error)
}
