package postgres

import (
	"fmt"
	"log/slog"

	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	infraLogger "github.com/aralvesandrade/cafeos/internal/infra/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewConnection(databaseURL string, log *slog.Logger, gormLevel slog.Level) (*gorm.DB, error) {
	gormLogger := infraLogger.NewGORMLogger(log, gormLevel)
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)

	db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)

	if err := db.AutoMigrate(
		// Global catalog entities — not organization-scoped, shared by the
		// whole platform (plans on offer, fixed application modules, the
		// roles a user can hold).
		&entity.Plan{},
		&entity.Module{},
		&entity.Role{},

		// Organization and everything organization-scoped below. Tenant
		// isolation runs through OrganizationID on every entity in this
		// group; RolePermission is the one exception worth noting — it's
		// the per-organization access matrix over the global Role/Module
		// catalogs above (organization × role × module -> access level).
		&entity.Organization{},
		&entity.RolePermission{},
		&entity.User{},
		&entity.Farm{},
		&entity.Producer{},
		&entity.Plot{},
		&entity.OperationType{},
		&entity.Operation{},
		&entity.Harvest{},
		&entity.HarvestProduction{},
		&entity.Indicator{},
		&entity.AgriculturalProduct{},
		&entity.FinancialTransaction{},
		&entity.StockItem{},
		&entity.StockMovement{},
		&entity.Vehicle{},
		&entity.Maintenance{},
		&entity.Team{},
		&entity.Worker{},
		&entity.WorkShift{},
		&entity.CostCenter{},
		&entity.Budget{},
		&entity.CostAllocation{},
		&entity.CostAllocationItem{},
		&entity.Alert{},
	); err != nil {
		return nil, fmt.Errorf("auto migrate: %w", err)
	}

	return db, nil
}
