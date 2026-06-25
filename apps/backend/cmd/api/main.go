// @title CafeOS API
// @version 0.1.0
// @description Plataforma SaaS especializada em cafeicultura para gestão operacional, produtiva, financeira e analítica de propriedades cafeeiras.
// @host localhost:8080
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
)

func main() {
	cfg := config.Load()

	db, err := postgres.NewConnection(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	eventBus := event.NewInMemoryBus()
	event.SetupHandlers(eventBus)

	router := api.NewRouter(db, eventBus, cfg.JWTSecret)

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
