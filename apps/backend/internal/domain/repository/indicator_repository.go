package repository

import "github.com/aralvesandrade/cafeos/internal/domain/entity"

type IndicatorRepository interface {
	Create(indicator *entity.Indicator) error
	ListByHarvest(harvestID string) ([]*entity.Indicator, error)
	ListByOrganization(organizationID string) ([]*entity.Indicator, error)
	ListByOrganizationAndType(organizationID string, indicatorType entity.IndicatorType) ([]*entity.Indicator, error)
}
