package repository

import (
	"database/sql"
	"fmt"

	"github.com/aralvesandrade/cafeos/internal/domain/entity"
)

type TenantRepository struct {
	db *sql.DB
}

func NewTenantRepository(db *sql.DB) *TenantRepository {
	return &TenantRepository{db: db}
}

func (r *TenantRepository) Create(t *entity.Tenant) error {
	query := `INSERT INTO tenants (id, name, slug, brand_name, logo_url, primary_color, plan, domain)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.Exec(query, t.ID, t.Name, t.Slug, t.BrandName, t.LogoURL, t.PrimaryColor, t.Plan, t.Domain)
	if err != nil {
		return fmt.Errorf("create tenant: %w", err)
	}
	return nil
}

func (r *TenantRepository) GetByID(id string) (*entity.Tenant, error) {
	query := `SELECT id, name, slug, brand_name, logo_url, primary_color, plan, domain, created_at, updated_at
		FROM tenants WHERE id = $1`
	t := &entity.Tenant{}
	err := r.db.QueryRow(query, id).Scan(&t.ID, &t.Name, &t.Slug, &t.BrandName, &t.LogoURL, &t.PrimaryColor, &t.Plan, &t.Domain, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("get tenant by id: %w", err)
	}
	return t, nil
}

func (r *TenantRepository) GetBySlug(slug string) (*entity.Tenant, error) {
	query := `SELECT id, name, slug, brand_name, logo_url, primary_color, plan, domain, created_at, updated_at
		FROM tenants WHERE slug = $1`
	t := &entity.Tenant{}
	err := r.db.QueryRow(query, slug).Scan(&t.ID, &t.Name, &t.Slug, &t.BrandName, &t.LogoURL, &t.PrimaryColor, &t.Plan, &t.Domain, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("get tenant by slug: %w", err)
	}
	return t, nil
}

func (r *TenantRepository) List() ([]*entity.Tenant, error) {
	query := `SELECT id, name, slug, brand_name, logo_url, primary_color, plan, domain, created_at, updated_at
		FROM tenants ORDER BY name`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("list tenants: %w", err)
	}
	defer rows.Close()

	var tenants []*entity.Tenant
	for rows.Next() {
		t := &entity.Tenant{}
		if err := rows.Scan(&t.ID, &t.Name, &t.Slug, &t.BrandName, &t.LogoURL, &t.PrimaryColor, &t.Plan, &t.Domain, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan tenant: %w", err)
		}
		tenants = append(tenants, t)
	}
	return tenants, nil
}

func (r *TenantRepository) Update(t *entity.Tenant) error {
	query := `UPDATE tenants SET name=$1, slug=$2, brand_name=$3, logo_url=$4, primary_color=$5, plan=$6, domain=$7, updated_at=NOW()
		WHERE id=$8`
	_, err := r.db.Exec(query, t.Name, t.Slug, t.BrandName, t.LogoURL, t.PrimaryColor, t.Plan, t.Domain, t.ID)
	if err != nil {
		return fmt.Errorf("update tenant: %w", err)
	}
	return nil
}
