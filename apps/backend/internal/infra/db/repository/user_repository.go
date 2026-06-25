package repository

import (
	"database/sql"
	"fmt"

	"github.com/aralvesandrade/cafeos/internal/domain/entity"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(u *entity.User) error {
	query := `INSERT INTO users (id, tenant_id, name, email, password_hash, role, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.Exec(query, u.ID, u.TenantID, u.Name, u.Email, u.PasswordHash, u.Role, u.IsActive)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (r *UserRepository) GetByID(id string) (*entity.User, error) {
	query := `SELECT id, tenant_id, name, email, password_hash, role, is_active, created_at, updated_at
		FROM users WHERE id = $1`
	u := &entity.User{}
	err := r.db.QueryRow(query, id).Scan(&u.ID, &u.TenantID, &u.Name, &u.Email, &u.PasswordHash, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return u, nil
}

func (r *UserRepository) GetByEmail(email string) (*entity.User, error) {
	query := `SELECT id, tenant_id, name, email, password_hash, role, is_active, created_at, updated_at
		FROM users WHERE email = $1`
	u := &entity.User{}
	err := r.db.QueryRow(query, email).Scan(&u.ID, &u.TenantID, &u.Name, &u.Email, &u.PasswordHash, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("get user by email: %w", err)
	}
	return u, nil
}

func (r *UserRepository) ListByTenant(tenantID string) ([]*entity.User, error) {
	query := `SELECT id, tenant_id, name, email, password_hash, role, is_active, created_at, updated_at
		FROM users WHERE tenant_id = $1 ORDER BY name`
	rows, err := r.db.Query(query, tenantID)
	if err != nil {
		return nil, fmt.Errorf("list users by tenant: %w", err)
	}
	defer rows.Close()

	var users []*entity.User
	for rows.Next() {
		u := &entity.User{}
		if err := rows.Scan(&u.ID, &u.TenantID, &u.Name, &u.Email, &u.PasswordHash, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *UserRepository) Update(u *entity.User) error {
	query := `UPDATE users SET name=$1, email=$2, role=$3, is_active=$4, updated_at=NOW() WHERE id=$5`
	_, err := r.db.Exec(query, u.Name, u.Email, u.Role, u.IsActive, u.ID)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	return nil
}
