package repository

import (
	"database/sql"
	"fmt"

	"github.com/aralvesandrade/cafeos/internal/domain/entity"
)

type OperationRepository struct {
	db *sql.DB
}

func NewOperationRepository(db *sql.DB) *OperationRepository {
	return &OperationRepository{db: db}
}

func (r *OperationRepository) Create(op *entity.Operation) error {
	query := `INSERT INTO operations (id, tenant_id, plot_id, type, date, responsible, product_used, quantity, cost, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := r.db.Exec(query, op.ID, op.TenantID, op.PlotID, op.Type, op.Date, op.Responsible, op.ProductUsed, op.Quantity, op.Cost, op.Notes)
	if err != nil {
		return fmt.Errorf("create operation: %w", err)
	}
	return nil
}

func (r *OperationRepository) GetByID(id string) (*entity.Operation, error) {
	query := `SELECT id, tenant_id, plot_id, type, date, responsible, product_used, quantity, cost, notes, created_at
		FROM operations WHERE id = $1`
	op := &entity.Operation{}
	err := r.db.QueryRow(query, id).Scan(&op.ID, &op.TenantID, &op.PlotID, &op.Type, &op.Date, &op.Responsible, &op.ProductUsed, &op.Quantity, &op.Cost, &op.Notes, &op.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("get operation by id: %w", err)
	}
	return op, nil
}

func (r *OperationRepository) ListByPlot(plotID string) ([]*entity.Operation, error) {
	query := `SELECT id, tenant_id, plot_id, type, date, responsible, product_used, quantity, cost, notes, created_at
		FROM operations WHERE plot_id = $1 ORDER BY date DESC`
	rows, err := r.db.Query(query, plotID)
	if err != nil {
		return nil, fmt.Errorf("list operations by plot: %w", err)
	}
	defer rows.Close()

	var ops []*entity.Operation
	for rows.Next() {
		op := &entity.Operation{}
		if err := rows.Scan(&op.ID, &op.TenantID, &op.PlotID, &op.Type, &op.Date, &op.Responsible, &op.ProductUsed, &op.Quantity, &op.Cost, &op.Notes, &op.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan operation: %w", err)
		}
		ops = append(ops, op)
	}
	return ops, nil
}

func (r *OperationRepository) ListByTenant(tenantID string) ([]*entity.Operation, error) {
	query := `SELECT id, tenant_id, plot_id, type, date, responsible, product_used, quantity, cost, notes, created_at
		FROM operations WHERE tenant_id = $1 ORDER BY date DESC`
	rows, err := r.db.Query(query, tenantID)
	if err != nil {
		return nil, fmt.Errorf("list operations by tenant: %w", err)
	}
	defer rows.Close()

	var ops []*entity.Operation
	for rows.Next() {
		op := &entity.Operation{}
		if err := rows.Scan(&op.ID, &op.TenantID, &op.PlotID, &op.Type, &op.Date, &op.Responsible, &op.ProductUsed, &op.Quantity, &op.Cost, &op.Notes, &op.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan operation: %w", err)
		}
		ops = append(ops, op)
	}
	return ops, nil
}

func (r *OperationRepository) ListByTenantAndPeriod(tenantID string, start, end string) ([]*entity.Operation, error) {
	query := `SELECT id, tenant_id, plot_id, type, date, responsible, product_used, quantity, cost, notes, created_at
		FROM operations WHERE tenant_id = $1 AND date >= $2 AND date <= $3 ORDER BY date DESC`
	rows, err := r.db.Query(query, tenantID, start, end)
	if err != nil {
		return nil, fmt.Errorf("list operations by tenant and period: %w", err)
	}
	defer rows.Close()

	var ops []*entity.Operation
	for rows.Next() {
		op := &entity.Operation{}
		if err := rows.Scan(&op.ID, &op.TenantID, &op.PlotID, &op.Type, &op.Date, &op.Responsible, &op.ProductUsed, &op.Quantity, &op.Cost, &op.Notes, &op.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan operation: %w", err)
		}
		ops = append(ops, op)
	}
	return ops, nil
}

func (r *OperationRepository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM operations WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete operation: %w", err)
	}
	return nil
}
