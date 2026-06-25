package repository

import (
	"database/sql"
	"fmt"

	"github.com/aralvesandrade/cafeos/internal/domain/entity"
)

type FarmRepository struct {
	db *sql.DB
}

func NewFarmRepository(db *sql.DB) *FarmRepository {
	return &FarmRepository{db: db}
}

func (r *FarmRepository) Create(f *entity.Farm) error {
	query := `INSERT INTO farms (id, tenant_id, name, owner, location, total_area_ha, planted_area_ha)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.Exec(query, f.ID, f.TenantID, f.Name, f.Owner, f.Location, f.TotalAreaHA, f.PlantedAreaHA)
	if err != nil {
		return fmt.Errorf("create farm: %w", err)
	}
	return nil
}

func (r *FarmRepository) GetByID(id string) (*entity.Farm, error) {
	query := `SELECT id, tenant_id, name, owner, location, total_area_ha, planted_area_ha, created_at, updated_at
		FROM farms WHERE id = $1`
	f := &entity.Farm{}
	err := r.db.QueryRow(query, id).Scan(&f.ID, &f.TenantID, &f.Name, &f.Owner, &f.Location, &f.TotalAreaHA, &f.PlantedAreaHA, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("get farm by id: %w", err)
	}
	return f, nil
}

func (r *FarmRepository) ListByTenant(tenantID string) ([]*entity.Farm, error) {
	query := `SELECT id, tenant_id, name, owner, location, total_area_ha, planted_area_ha, created_at, updated_at
		FROM farms WHERE tenant_id = $1 ORDER BY name`
	rows, err := r.db.Query(query, tenantID)
	if err != nil {
		return nil, fmt.Errorf("list farms by tenant: %w", err)
	}
	defer rows.Close()

	var farms []*entity.Farm
	for rows.Next() {
		f := &entity.Farm{}
		if err := rows.Scan(&f.ID, &f.TenantID, &f.Name, &f.Owner, &f.Location, &f.TotalAreaHA, &f.PlantedAreaHA, &f.CreatedAt, &f.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan farm: %w", err)
		}
		farms = append(farms, f)
	}
	return farms, nil
}

func (r *FarmRepository) Update(f *entity.Farm) error {
	query := `UPDATE farms SET name=$1, owner=$2, location=$3, total_area_ha=$4, planted_area_ha=$5, updated_at=NOW() WHERE id=$6`
	_, err := r.db.Exec(query, f.Name, f.Owner, f.Location, f.TotalAreaHA, f.PlantedAreaHA, f.ID)
	if err != nil {
		return fmt.Errorf("update farm: %w", err)
	}
	return nil
}

func (r *FarmRepository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM farms WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete farm: %w", err)
	}
	return nil
}
