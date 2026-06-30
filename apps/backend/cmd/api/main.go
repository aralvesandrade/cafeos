// @title CafeOS API
// @version 0.1.0
// @description Plataforma SaaS especializada em cafeicultura para gestão operacional, produtiva, financeira e analítica de propriedades cafeeiras.
// @host localhost:5001
// @BasePath /api/v1/{tenant_id}
// @schemes http
// @license.name Proprietary
// @securityDefinitions.apiKey BearerAuth
// @in header
// @name Authorization
// @description JWT token. Use: Bearer {token}
package main

import (
	"log/slog"
	"net/http"
	"os"

	api "github.com/aralvesandrade/cafeos/internal/api"
	"github.com/aralvesandrade/cafeos/internal/event"
	"github.com/aralvesandrade/cafeos/internal/infra/config"
	"github.com/aralvesandrade/cafeos/internal/infra/db/postgres"
	infraLogger "github.com/aralvesandrade/cafeos/internal/infra/logger"
	"github.com/aralvesandrade/cafeos/internal/infra/messaging"
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

	eventBus := event.NewInMemoryBus()
	event.SetupHandlers(eventBus)

	var pub *messaging.Publisher
	if cfg.RabbitMQURL != "" {
		conn, err := messaging.NewConnection(cfg.RabbitMQURL)
		if err != nil {
			log.Warn("rabbitmq connection failed, sync disabled", "error", err)
		} else {
			defer conn.Close()
			pub = messaging.NewPublisher(conn.Channel())
			queues := []string{"sync.operations", "sync.stock", "sync.harvest", "sync.financial", "sync.labor"}
			for _, q := range queues {
				if err := pub.DeclareQueue(q); err != nil {
					log.Warn("rabbitmq queue declare failed", "queue", q, "error", err)
				}
			}
			log.Info("rabbitmq connected, sync enabled")
		}
	}

	router := api.NewRouter(db, eventBus, pub, cfg.JWTSecret, log)

	port := cfg.ServerPort
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	addr := ":" + port
	log.Info("starting cafeos api", "addr", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Error("server failed", "error", err)
		os.Exit(1)
	}
}
