package repository

import (
	"database/sql"
	"fmt"

	"github.com/aralvesandrade/cafeos/internal/domain/entity"
)

type IndicatorRepository struct {
	db *sql.DB
}

func NewIndicatorRepository(db *sql.DB) *IndicatorRepository {
	return &IndicatorRepository{db: db}
}

func (r *IndicatorRepository) Create(ind *entity.Indicator) error {
	query := `INSERT INTO indicators (id, tenant_id, harvest_id, plot_id, type, value, calculated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.Exec(query, ind.ID, ind.TenantID, ind.HarvestID, ind.PlotID, ind.Type, ind.Value, ind.CalculatedAt)
	if err != nil {
		return fmt.Errorf("create indicator: %w", err)
	}
	return nil
}

func (r *IndicatorRepository) ListByHarvest(harvestID string) ([]*entity.Indicator, error) {
	query := `SELECT id, tenant_id, harvest_id, plot_id, type, value, calculated_at
		FROM indicators WHERE harvest_id = $1 ORDER BY type`
	rows, err := r.db.Query(query, harvestID)
	if err != nil {
		return nil, fmt.Errorf("list indicators by harvest: %w", err)
	}
	defer rows.Close()

	var indicators []*entity.Indicator
	for rows.Next() {
		ind := &entity.Indicator{}
		if err := rows.Scan(&ind.ID, &ind.TenantID, &ind.HarvestID, &ind.PlotID, &ind.Type, &ind.Value, &ind.CalculatedAt); err != nil {
			return nil, fmt.Errorf("scan indicator: %w", err)
		}
		indicators = append(indicators, ind)
	}
	return indicators, nil
}

func (r *IndicatorRepository) ListByTenant(tenantID string) ([]*entity.Indicator, error) {
	query := `SELECT id, tenant_id, harvest_id, plot_id, type, value, calculated_at
		FROM indicators WHERE tenant_id = $1 ORDER BY calculated_at DESC`
	rows, err := r.db.Query(query, tenantID)
	if err != nil {
		return nil, fmt.Errorf("list indicators by tenant: %w", err)
	}
	defer rows.Close()

	var indicators []*entity.Indicator
	for rows.Next() {
		ind := &entity.Indicator{}
		if err := rows.Scan(&ind.ID, &ind.TenantID, &ind.HarvestID, &ind.PlotID, &ind.Type, &ind.Value, &ind.CalculatedAt); err != nil {
			return nil, fmt.Errorf("scan indicator: %w", err)
		}
		indicators = append(indicators, ind)
	}
	return indicators, nil
}

func (r *IndicatorRepository) ListByTenantAndType(tenantID string, indicatorType entity.IndicatorType) ([]*entity.Indicator, error) {
	query := `SELECT id, tenant_id, harvest_id, plot_id, type, value, calculated_at
		FROM indicators WHERE tenant_id = $1 AND type = $2 ORDER BY calculated_at DESC`
	rows, err := r.db.Query(query, tenantID, indicatorType)
	if err != nil {
		return nil, fmt.Errorf("list indicators by tenant and type: %w", err)
	}
	defer rows.Close()

	var indicators []*entity.Indicator
	for rows.Next() {
		ind := &entity.Indicator{}
		if err := rows.Scan(&ind.ID, &ind.TenantID, &ind.HarvestID, &ind.PlotID, &ind.Type, &ind.Value, &ind.CalculatedAt); err != nil {
			return nil, fmt.Errorf("scan indicator: %w", err)
		}
		indicators = append(indicators, ind)
	}
	return indicators, nil
}
