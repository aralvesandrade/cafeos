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
		&entity.Organization{},
		&entity.User{},
		&entity.Farm{},
		&entity.Producer{},
		&entity.Plot{},
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
	); err != nil {
		return nil, fmt.Errorf("auto migrate: %w", err)
	}

	return db, nil
}
