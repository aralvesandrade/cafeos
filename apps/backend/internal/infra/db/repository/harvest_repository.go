package repository

import (
	"database/sql"
	"fmt"

	"github.com/aralvesandrade/cafeos/internal/domain/entity"
)

type HarvestRepository struct {
	db *sql.DB
}

func NewHarvestRepository(db *sql.DB) *HarvestRepository {
	return &HarvestRepository{db: db}
}

func (r *HarvestRepository) Create(h *entity.Harvest) error {
	query := `INSERT INTO harvests (id, tenant_id, year, description, estimated_production, status)
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.Exec(query, h.ID, h.TenantID, h.Year, h.Description, h.EstimatedProduction, h.Status)
	if err != nil {
		return fmt.Errorf("create harvest: %w", err)
	}
	return nil
}

func (r *HarvestRepository) GetByID(id string) (*entity.Harvest, error) {
	query := `SELECT id, tenant_id, year, description, estimated_production, status, created_at, updated_at
		FROM harvests WHERE id = $1`
	h := &entity.Harvest{}
	err := r.db.QueryRow(query, id).Scan(&h.ID, &h.TenantID, &h.Year, &h.Description, &h.EstimatedProduction, &h.Status, &h.CreatedAt, &h.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("get harvest by id: %w", err)
	}
	return h, nil
}

func (r *HarvestRepository) ListByTenant(tenantID string) ([]*entity.Harvest, error) {
	query := `SELECT id, tenant_id, year, description, estimated_production, status, created_at, updated_at
		FROM harvests WHERE tenant_id = $1 ORDER BY year DESC`
	rows, err := r.db.Query(query, tenantID)
	if err != nil {
		return nil, fmt.Errorf("list harvests by tenant: %w", err)
	}
	defer rows.Close()

	var harvests []*entity.Harvest
	for rows.Next() {
		h := &entity.Harvest{}
		if err := rows.Scan(&h.ID, &h.TenantID, &h.Year, &h.Description, &h.EstimatedProduction, &h.Status, &h.CreatedAt, &h.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan harvest: %w", err)
		}
		harvests = append(harvests, h)
	}
	return harvests, nil
}

func (r *HarvestRepository) Update(h *entity.Harvest) error {
	query := `UPDATE harvests SET year=$1, description=$2, estimated_production=$3, status=$4, updated_at=NOW() WHERE id=$5`
	_, err := r.db.Exec(query, h.Year, h.Description, h.EstimatedProduction, h.Status, h.ID)
	if err != nil {
		return fmt.Errorf("update harvest: %w", err)
	}
	return nil
}

type HarvestProductionRepository struct {
	db *sql.DB
}

func NewHarvestProductionRepository(db *sql.DB) *HarvestProductionRepository {
	return &HarvestProductionRepository{db: db}
}

func (r *HarvestProductionRepository) Create(hp *entity.HarvestProduction) error {
	query := `INSERT INTO harvest_productions (id, tenant_id, harvest_id, plot_id, quantity, recorded_at, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.Exec(query, hp.ID, hp.TenantID, hp.HarvestID, hp.PlotID, hp.Quantity, hp.RecordedAt, hp.Notes)
	if err != nil {
		return fmt.Errorf("create harvest production: %w", err)
	}
	return nil
}

func (r *HarvestProductionRepository) GetByID(id string) (*entity.HarvestProduction, error) {
	query := `SELECT id, tenant_id, harvest_id, plot_id, quantity, recorded_at, notes
		FROM harvest_productions WHERE id = $1`
	hp := &entity.HarvestProduction{}
	err := r.db.QueryRow(query, id).Scan(&hp.ID, &hp.TenantID, &hp.HarvestID, &hp.PlotID, &hp.Quantity, &hp.RecordedAt, &hp.Notes)
	if err != nil {
		return nil, fmt.Errorf("get harvest production by id: %w", err)
	}
	return hp, nil
}

func (r *HarvestProductionRepository) ListByHarvest(harvestID string) ([]*entity.HarvestProduction, error) {
	query := `SELECT id, tenant_id, harvest_id, plot_id, quantity, recorded_at, notes
		FROM harvest_productions WHERE harvest_id = $1 ORDER BY plot_id`
	rows, err := r.db.Query(query, harvestID)
	if err != nil {
		return nil, fmt.Errorf("list harvest productions by harvest: %w", err)
	}
	defer rows.Close()

	var prods []*entity.HarvestProduction
	for rows.Next() {
		hp := &entity.HarvestProduction{}
		if err := rows.Scan(&hp.ID, &hp.TenantID, &hp.HarvestID, &hp.PlotID, &hp.Quantity, &hp.RecordedAt, &hp.Notes); err != nil {
			return nil, fmt.Errorf("scan harvest production: %w", err)
		}
		prods = append(prods, hp)
	}
	return prods, nil
}

func (r *HarvestProductionRepository) ListByPlot(plotID string) ([]*entity.HarvestProduction, error) {
	query := `SELECT id, tenant_id, harvest_id, plot_id, quantity, recorded_at, notes
		FROM harvest_productions WHERE plot_id = $1 ORDER BY recorded_at DESC`
	rows, err := r.db.Query(query, plotID)
	if err != nil {
		return nil, fmt.Errorf("list harvest productions by plot: %w", err)
	}
	defer rows.Close()

	var prods []*entity.HarvestProduction
	for rows.Next() {
		hp := &entity.HarvestProduction{}
		if err := rows.Scan(&hp.ID, &hp.TenantID, &hp.HarvestID, &hp.PlotID, &hp.Quantity, &hp.RecordedAt, &hp.Notes); err != nil {
			return nil, fmt.Errorf("scan harvest production: %w", err)
		}
		prods = append(prods, hp)
	}
	return prods, nil
}

func (r *HarvestProductionRepository) Update(hp *entity.HarvestProduction) error {
	query := `UPDATE harvest_productions SET quantity=$1, notes=$2 WHERE id=$3`
	_, err := r.db.Exec(query, hp.Quantity, hp.Notes, hp.ID)
	if err != nil {
		return fmt.Errorf("update harvest production: %w", err)
	}
	return nil
}
