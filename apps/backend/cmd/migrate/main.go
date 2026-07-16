package main

import (
	"log/slog"
	"os"

	"github.com/aralvesandrade/cafeos/internal/infra/config"
	"github.com/aralvesandrade/cafeos/internal/infra/db/postgres"
	infraLogger "github.com/aralvesandrade/cafeos/internal/infra/logger"
)

func main() {
	cfg := config.Load()

	log := infraLogger.New(infraLogger.Config{
		Level:  cfg.LogLevel,
		Format: cfg.LogFormat,
	})
	slog.SetDefault(log)

	db, err := postgres.NewConnection(cfg.DatabaseURL, log, slog.LevelInfo)
	if err != nil {
		log.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Error("failed to get sql.DB", "error", err)
		os.Exit(1)
	}
	defer sqlDB.Close()

	log.Info("schema migrated successfully")
}
