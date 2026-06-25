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
	"log"
	"net/http"
	"os"

	api "github.com/aralvesandrade/cafeos/internal/api"
	"github.com/aralvesandrade/cafeos/internal/event"
	"github.com/aralvesandrade/cafeos/internal/infra/config"
	"github.com/aralvesandrade/cafeos/internal/infra/db/postgres"
	"github.com/aralvesandrade/cafeos/internal/infra/messaging"
)

func main() {
	cfg := config.Load()

	db, err := postgres.NewConnection(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get sql.DB: %v", err)
	}
	defer sqlDB.Close()

	eventBus := event.NewInMemoryBus()
	event.SetupHandlers(eventBus)

	var pub *messaging.Publisher
	if cfg.RabbitMQURL != "" {
		conn, err := messaging.NewConnection(cfg.RabbitMQURL)
		if err != nil {
			log.Printf("[RABBITMQ] connection failed (sync disabled): %v", err)
		} else {
			defer conn.Close()
			pub = messaging.NewPublisher(conn.Channel())
			queues := []string{"sync.operations", "sync.stock", "sync.harvest", "sync.financial", "sync.labor"}
			for _, q := range queues {
				if err := pub.DeclareQueue(q); err != nil {
					log.Printf("[RABBITMQ] queue declare failed: %v", err)
				}
			}
			log.Println("[RABBITMQ] connected, sync enabled")
		}
	}

	router := api.NewRouter(db, eventBus, pub, cfg.JWTSecret)

	port := cfg.ServerPort
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	addr := ":" + port
	log.Printf("CafeOS API starting on %s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
