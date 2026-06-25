package repository

import (
	"database/sql"
	"fmt"

	"github.com/aralvesandrade/cafeos/internal/domain/entity"
)

type AgriculturalProductRepository struct {
	db *sql.DB
}

func NewAgriculturalProductRepository(db *sql.DB) *AgriculturalProductRepository {
	return &AgriculturalProductRepository{db: db}
}

func (r *AgriculturalProductRepository) Create(p *entity.AgriculturalProduct) error {
	query := `INSERT INTO agricultural_products (id, tenant_id, name, type, unit) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(query, p.ID, p.TenantID, p.Name, p.Type, p.Unit)
	if err != nil {
		return fmt.Errorf("create agricultural product: %w", err)
	}
	return nil
}

func (r *AgriculturalProductRepository) GetByID(id string) (*entity.AgriculturalProduct, error) {
	query := `SELECT id, tenant_id, name, type, unit, created_at FROM agricultural_products WHERE id = $1`
	p := &entity.AgriculturalProduct{}
	err := r.db.QueryRow(query, id).Scan(&p.ID, &p.TenantID, &p.Name, &p.Type, &p.Unit, &p.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("get agricultural product by id: %w", err)
	}
	return p, nil
}

func (r *AgriculturalProductRepository) ListByTenant(tenantID string) ([]*entity.AgriculturalProduct, error) {
	query := `SELECT id, tenant_id, name, type, unit, created_at FROM agricultural_products WHERE tenant_id = $1 ORDER BY name`
	rows, err := r.db.Query(query, tenantID)
	if err != nil {
		return nil, fmt.Errorf("list agricultural products by tenant: %w", err)
	}
	defer rows.Close()

	var products []*entity.AgriculturalProduct
	for rows.Next() {
		p := &entity.AgriculturalProduct{}
		if err := rows.Scan(&p.ID, &p.TenantID, &p.Name, &p.Type, &p.Unit, &p.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan agricultural product: %w", err)
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *AgriculturalProductRepository) Update(p *entity.AgriculturalProduct) error {
	query := `UPDATE agricultural_products SET name=$1, type=$2, unit=$3 WHERE id=$4`
	_, err := r.db.Exec(query, p.Name, p.Type, p.Unit, p.ID)
	if err != nil {
		return fmt.Errorf("update agricultural product: %w", err)
	}
	return nil
}

func (r *AgriculturalProductRepository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM agricultural_products WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete agricultural product: %w", err)
	}
	return nil
}
