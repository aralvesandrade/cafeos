package api

import (
	"net/http"

	"github.com/aralvesandrade/cafeos/internal/api/handler"
	"github.com/aralvesandrade/cafeos/internal/api/middleware"
	domainSvc "github.com/aralvesandrade/cafeos/internal/domain/service"
	"github.com/aralvesandrade/cafeos/internal/event"
	infraRepo "github.com/aralvesandrade/cafeos/internal/infra/db/repository"
	"github.com/aralvesandrade/cafeos/internal/infra/db/postgres"
	"gorm.io/gorm"

	_ "github.com/aralvesandrade/cafeos/docs"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func NewRouter(db *gorm.DB, eventBus event.Bus, jwtSecret string) http.Handler {
	mux := http.NewServeMux()

	farmRepo := infraRepo.NewFarmRepository(db)
	plotRepo := infraRepo.NewPlotRepository(db)
	opRepo := infraRepo.NewOperationRepository(db)
	harvestRepo := infraRepo.NewHarvestRepository(db)
	hpRepo := infraRepo.NewHarvestProductionRepository(db)
	indicatorRepo := infraRepo.NewIndicatorRepository(db)
	tenantRepo := infraRepo.NewTenantRepository(db)
	userRepo := infraRepo.NewUserRepository(db)
	transactor := postgres.NewTransactor(db)

	farmSvc := domainSvc.NewFarmService(farmRepo)
	plotSvc := domainSvc.NewPlotService(plotRepo)
	opSvc := domainSvc.NewOperationService(opRepo, eventBus)
	harvestSvc := domainSvc.NewHarvestService(harvestRepo, hpRepo, indicatorRepo, plotRepo, opRepo, transactor, eventBus)

	farmH := handler.NewFarmHandler(farmSvc)
	plotH := handler.NewPlotHandler(plotSvc)
	opH := handler.NewOperationHandler(opSvc)
	harvestH := handler.NewHarvestHandler(harvestSvc)
	dashboardH := handler.NewDashboardHandler(harvestRepo, indicatorRepo, opRepo, plotRepo, farmRepo)
	authH := handler.NewAuthHandler(userRepo, tenantRepo, jwtSecret)

	authMw := middleware.Auth(jwtSecret)
	corsMw := middleware.CORS

	// Auth (public)
	mux.HandleFunc("POST /auth/login", authH.Login)

	// Health
	mux.HandleFunc("GET /health", handler.HealthCheck)

	// Swagger UI
	mux.Handle("GET /swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	// API v1 routes
	chain := func(h http.HandlerFunc) http.Handler {
		return authMw(http.HandlerFunc(h))
	}

	// Farms
	mux.Handle("POST /api/v1/{tenant_id}/farms", chain(farmH.Create))
	mux.Handle("GET /api/v1/{tenant_id}/farms", chain(farmH.List))
	mux.Handle("GET /api/v1/{tenant_id}/farms/{id}", chain(farmH.GetByID))
	mux.Handle("PUT /api/v1/{tenant_id}/farms/{id}", chain(farmH.Update))
	mux.Handle("DELETE /api/v1/{tenant_id}/farms/{id}", chain(farmH.Delete))

	// Plots
	mux.Handle("POST /api/v1/{tenant_id}/plots", chain(plotH.Create))
	mux.Handle("GET /api/v1/{tenant_id}/plots", chain(plotH.List))
	mux.Handle("GET /api/v1/{tenant_id}/plots/{id}", chain(plotH.GetByID))
	mux.Handle("PUT /api/v1/{tenant_id}/plots/{id}", chain(plotH.Update))
	mux.Handle("DELETE /api/v1/{tenant_id}/plots/{id}", chain(plotH.Delete))
	mux.Handle("GET /api/v1/{tenant_id}/farms/{farm_id}/plots", chain(plotH.ListByFarm))

	// Operations
	mux.Handle("POST /api/v1/{tenant_id}/operations", chain(opH.Create))
	mux.Handle("GET /api/v1/{tenant_id}/operations", chain(opH.List))
	mux.Handle("GET /api/v1/{tenant_id}/operations/recent", chain(opH.ListRecent))
	mux.Handle("GET /api/v1/{tenant_id}/operations/{id}", chain(opH.GetByID))
	mux.Handle("DELETE /api/v1/{tenant_id}/operations/{id}", chain(opH.Delete))
	mux.Handle("GET /api/v1/{tenant_id}/plots/{plot_id}/operations", chain(opH.ListByPlot))

	// Harvests
	mux.Handle("POST /api/v1/{tenant_id}/harvests", chain(harvestH.Create))
	mux.Handle("GET /api/v1/{tenant_id}/harvests", chain(harvestH.List))
	mux.Handle("GET /api/v1/{tenant_id}/harvests/{id}", chain(harvestH.GetByID))
	mux.Handle("PUT /api/v1/{tenant_id}/harvests/{id}/finalize", chain(harvestH.Finalize))
	mux.Handle("POST /api/v1/{tenant_id}/harvests/{id}/production", chain(harvestH.RecordProduction))
	mux.Handle("GET /api/v1/{tenant_id}/harvests/{id}/production", chain(harvestH.GetProduction))

	// Dashboard
	mux.Handle("GET /api/v1/{tenant_id}/dashboard", chain(dashboardH.GetDashboard))

	return corsMw(mux)
}
