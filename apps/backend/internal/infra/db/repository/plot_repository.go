package repository

import (
	"database/sql"
	"fmt"

	"github.com/aralvesandrade/cafeos/internal/domain/entity"
)

type PlotRepository struct {
	db *sql.DB
}

func NewPlotRepository(db *sql.DB) *PlotRepository {
	return &PlotRepository{db: db}
}

func (r *PlotRepository) Create(p *entity.Plot) error {
	query := `INSERT INTO plots (id, tenant_id, farm_id, name, area_ha, cultivar, planting_year, altitude, soil_type)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.db.Exec(query, p.ID, p.TenantID, p.FarmID, p.Name, p.AreaHA, p.Cultivar, p.PlantingYear, p.Altitude, p.SoilType)
	if err != nil {
		return fmt.Errorf("create plot: %w", err)
	}
	return nil
}

func (r *PlotRepository) GetByID(id string) (*entity.Plot, error) {
	query := `SELECT id, tenant_id, farm_id, name, area_ha, cultivar, planting_year, altitude, soil_type, created_at, updated_at
		FROM plots WHERE id = $1`
	p := &entity.Plot{}
	err := r.db.QueryRow(query, id).Scan(&p.ID, &p.TenantID, &p.FarmID, &p.Name, &p.AreaHA, &p.Cultivar, &p.PlantingYear, &p.Altitude, &p.SoilType, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("get plot by id: %w", err)
	}
	return p, nil
}

func (r *PlotRepository) ListByFarm(farmID string) ([]*entity.Plot, error) {
	query := `SELECT id, tenant_id, farm_id, name, area_ha, cultivar, planting_year, altitude, soil_type, created_at, updated_at
		FROM plots WHERE farm_id = $1 ORDER BY name`
	rows, err := r.db.Query(query, farmID)
	if err != nil {
		return nil, fmt.Errorf("list plots by farm: %w", err)
	}
	defer rows.Close()

	var plots []*entity.Plot
	for rows.Next() {
		p := &entity.Plot{}
		if err := rows.Scan(&p.ID, &p.TenantID, &p.FarmID, &p.Name, &p.AreaHA, &p.Cultivar, &p.PlantingYear, &p.Altitude, &p.SoilType, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan plot: %w", err)
		}
		plots = append(plots, p)
	}
	return plots, nil
}

func (r *PlotRepository) ListByTenant(tenantID string) ([]*entity.Plot, error) {
	query := `SELECT id, tenant_id, farm_id, name, area_ha, cultivar, planting_year, altitude, soil_type, created_at, updated_at
		FROM plots WHERE tenant_id = $1 ORDER BY name`
	rows, err := r.db.Query(query, tenantID)
	if err != nil {
		return nil, fmt.Errorf("list plots by tenant: %w", err)
	}
	defer rows.Close()

	var plots []*entity.Plot
	for rows.Next() {
		p := &entity.Plot{}
		if err := rows.Scan(&p.ID, &p.TenantID, &p.FarmID, &p.Name, &p.AreaHA, &p.Cultivar, &p.PlantingYear, &p.Altitude, &p.SoilType, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan plot: %w", err)
		}
		plots = append(plots, p)
	}
	return plots, nil
}

func (r *PlotRepository) Update(p *entity.Plot) error {
	query := `UPDATE plots SET name=$1, area_ha=$2, cultivar=$3, planting_year=$4, altitude=$5, soil_type=$6, updated_at=NOW() WHERE id=$7`
	_, err := r.db.Exec(query, p.Name, p.AreaHA, p.Cultivar, p.PlantingYear, p.Altitude, p.SoilType, p.ID)
	if err != nil {
		return fmt.Errorf("update plot: %w", err)
	}
	return nil
}

func (r *PlotRepository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM plots WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete plot: %w", err)
	}
	return nil
}
