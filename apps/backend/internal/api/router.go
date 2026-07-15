package api

import (
	"log/slog"
	"net/http"

	"github.com/aralvesandrade/cafeos/internal/api/handler"
	"github.com/aralvesandrade/cafeos/internal/api/middleware"
	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	domainSvc "github.com/aralvesandrade/cafeos/internal/domain/service"
	"github.com/aralvesandrade/cafeos/internal/event"
	"github.com/aralvesandrade/cafeos/internal/infra/db/postgres"
	infraRepo "github.com/aralvesandrade/cafeos/internal/infra/db/repository"
	"github.com/aralvesandrade/cafeos/internal/infra/messaging"
	"gorm.io/gorm"

	_ "github.com/aralvesandrade/cafeos/docs"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func NewRouter(db *gorm.DB, eventBus event.Bus, publisher *messaging.Publisher, jwtSecret string, log *slog.Logger) http.Handler {
	mux := http.NewServeMux()

	farmRepo := infraRepo.NewFarmRepository(db)
	producerRepo := infraRepo.NewProducerRepository(db)
	plotRepo := infraRepo.NewPlotRepository(db)
	opRepo := infraRepo.NewOperationRepository(db)
	harvestRepo := infraRepo.NewHarvestRepository(db)
	hpRepo := infraRepo.NewHarvestProductionRepository(db)
	indicatorRepo := infraRepo.NewIndicatorRepository(db)
	organizationRepo := infraRepo.NewOrganizationRepository(db)
	userRepo := infraRepo.NewUserRepository(db)
	finRepo := infraRepo.NewFinancialRepository(db)
	stockItemRepo := infraRepo.NewStockItemRepository(db)
	stockMovRepo := infraRepo.NewStockMovementRepository(db)
	vehRepo := infraRepo.NewVehicleRepository(db)
	maintRepo := infraRepo.NewMaintenanceRepository(db)
	teamRepo := infraRepo.NewTeamRepository(db)
	workerRepo := infraRepo.NewWorkerRepository(db)
	shiftRepo := infraRepo.NewWorkShiftRepository(db)
	prodRepo := infraRepo.NewAgriculturalProductRepository(db)
	ccRepo := infraRepo.NewCostCenterRepository(db)
	otRepo := infraRepo.NewOperationTypeRepository(db)
	budgetRepo := infraRepo.NewBudgetRepository(db)
	allocRepo := infraRepo.NewCostAllocationRepository(db)
	alertRepo := infraRepo.NewAlertRepository(db)
	transactor := postgres.NewTransactor(db)
	ruleEngine := domainSvc.NewRuleEngine()

	farmSvc := domainSvc.NewFarmService(farmRepo, producerRepo)
	plotSvc := domainSvc.NewPlotService(plotRepo)
	opSvc := domainSvc.NewOperationService(opRepo, eventBus)
	harvestSvc := domainSvc.NewHarvestService(harvestRepo, hpRepo, indicatorRepo, plotRepo, opRepo, maintRepo, shiftRepo, finRepo, allocRepo, transactor, eventBus, ruleEngine, alertRepo)
	finSvc := domainSvc.NewFinancialService(finRepo)
	stockSvc := domainSvc.NewStockService(stockItemRepo, stockMovRepo)
	fleetSvc := domainSvc.NewFleetService(vehRepo, maintRepo)
	laborSvc := domainSvc.NewLaborService(teamRepo, workerRepo, shiftRepo)
	ccSvc := domainSvc.NewCostCenterService(ccRepo)
	otSvc := domainSvc.NewOperationTypeService(otRepo)
	budgetSvc := domainSvc.NewBudgetService(budgetRepo, harvestRepo, opRepo, maintRepo, shiftRepo, finRepo, allocRepo)
	allocSvc := domainSvc.NewCostAllocationService(allocRepo, plotRepo)

	farmH := handler.NewFarmHandler(farmSvc)
	plotH := handler.NewPlotHandler(plotSvc, farmSvc)
	opH := handler.NewOperationHandler(opSvc, plotSvc, farmSvc)
	harvestH := handler.NewHarvestHandler(harvestSvc, plotSvc, farmSvc, indicatorRepo)
	dashboardH := handler.NewDashboardHandler(harvestRepo, indicatorRepo, opRepo, plotRepo, farmRepo, hpRepo, farmSvc)
	authH := handler.NewAuthHandler(userRepo, organizationRepo, jwtSecret)
	organizationH := handler.NewOrganizationHandler(organizationRepo)
	userH := handler.NewUserHandler(userRepo)
	finH := handler.NewFinancialHandler(finSvc, farmSvc)
	stockH := handler.NewStockHandler(stockSvc, farmSvc)
	fleetH := handler.NewFleetHandler(fleetSvc, farmSvc)
	laborH := handler.NewLaborHandler(laborSvc)
	prodH := handler.NewProductHandler(prodRepo)
	ccH := handler.NewCostCenterHandler(ccSvc)
	otH := handler.NewOperationTypeHandler(otSvc)
	budgetH := handler.NewBudgetHandler(budgetSvc)
	allocH := handler.NewCostAllocationHandler(allocSvc)
	alertH := handler.NewAlertHandler(alertRepo)

	authMw := middleware.Auth(jwtSecret)
	adminMw := middleware.RequireRole(entity.RolePlatformOwner)
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
	mux.Handle("POST /api/v1/{organization_id}/farms", chain(farmH.Create))
	mux.Handle("GET /api/v1/{organization_id}/farms", chain(farmH.List))
	mux.Handle("GET /api/v1/{organization_id}/farms/{id}", chain(farmH.GetByID))
	mux.Handle("PUT /api/v1/{organization_id}/farms/{id}", chain(farmH.Update))
	mux.Handle("DELETE /api/v1/{organization_id}/farms/{id}", chain(farmH.Delete))

	// Plots
	mux.Handle("POST /api/v1/{organization_id}/plots", chain(plotH.Create))
	mux.Handle("GET /api/v1/{organization_id}/plots", chain(plotH.List))
	mux.Handle("GET /api/v1/{organization_id}/plots/{id}", chain(plotH.GetByID))
	mux.Handle("PUT /api/v1/{organization_id}/plots/{id}", chain(plotH.Update))
	mux.Handle("DELETE /api/v1/{organization_id}/plots/{id}", chain(plotH.Delete))
	mux.Handle("GET /api/v1/{organization_id}/farms/{farm_id}/plots", chain(plotH.ListByFarm))

	// Operations
	mux.Handle("POST /api/v1/{organization_id}/operations", chain(opH.Create))
	mux.Handle("GET /api/v1/{organization_id}/operations", chain(opH.List))
	mux.Handle("GET /api/v1/{organization_id}/operations/recent", chain(opH.ListRecent))
	mux.Handle("GET /api/v1/{organization_id}/operations/{id}", chain(opH.GetByID))
	mux.Handle("PUT /api/v1/{organization_id}/operations/{id}", chain(opH.Update))
	mux.Handle("DELETE /api/v1/{organization_id}/operations/{id}", chain(opH.Delete))
	mux.Handle("GET /api/v1/{organization_id}/plots/{plot_id}/operations", chain(opH.ListByPlot))

	// Operation Types
	mux.Handle("POST /api/v1/{organization_id}/operation-types", chain(otH.Create))
	mux.Handle("GET /api/v1/{organization_id}/operation-types", chain(otH.List))
	mux.Handle("GET /api/v1/{organization_id}/operation-types/{id}", chain(otH.GetByID))
	mux.Handle("PUT /api/v1/{organization_id}/operation-types/{id}", chain(otH.Update))
	mux.Handle("DELETE /api/v1/{organization_id}/operation-types/{id}", chain(otH.Delete))

	// Harvests
	mux.Handle("POST /api/v1/{organization_id}/harvests", chain(harvestH.Create))
	mux.Handle("GET /api/v1/{organization_id}/harvests", chain(harvestH.List))
	mux.Handle("GET /api/v1/{organization_id}/harvests/{id}", chain(harvestH.GetByID))
	mux.Handle("PUT /api/v1/{organization_id}/harvests/{id}/finalize", chain(harvestH.Finalize))
	mux.Handle("POST /api/v1/{organization_id}/harvests/{id}/production", chain(harvestH.RecordProduction))
	mux.Handle("GET /api/v1/{organization_id}/harvests/{id}/production", chain(harvestH.GetProduction))
	mux.Handle("GET /api/v1/{organization_id}/harvests/{id}/indicators", chain(harvestH.GetIndicators))

	// Dashboard
	mux.Handle("GET /api/v1/{organization_id}/dashboard", chain(dashboardH.GetDashboard))

	// Agricultural Products
	mux.Handle("GET /api/v1/{organization_id}/agricultural-products", chain(prodH.List))

	// Financial (Phase 2)
	mux.Handle("POST /api/v1/{organization_id}/financial", chain(finH.Create))
	mux.Handle("GET /api/v1/{organization_id}/financial", chain(finH.List))
	mux.Handle("GET /api/v1/{organization_id}/financial/{id}", chain(finH.GetByID))
	mux.Handle("PUT /api/v1/{organization_id}/financial/{id}", chain(finH.Update))
	mux.Handle("DELETE /api/v1/{organization_id}/financial/{id}", chain(finH.Delete))

	// Stock (Phase 2)
	mux.Handle("POST /api/v1/{organization_id}/stock/items", chain(stockH.CreateItem))
	mux.Handle("GET /api/v1/{organization_id}/stock/items", chain(stockH.ListItems))
	mux.Handle("GET /api/v1/{organization_id}/stock/items/{id}", chain(stockH.GetItemByID))
	mux.Handle("PUT /api/v1/{organization_id}/stock/items/{id}", chain(stockH.UpdateItem))
	mux.Handle("DELETE /api/v1/{organization_id}/stock/items/{id}", chain(stockH.DeleteItem))
	mux.Handle("POST /api/v1/{organization_id}/stock/movements", chain(stockH.RecordMovement))
	mux.Handle("GET /api/v1/{organization_id}/stock/movements", chain(stockH.ListMovements))

	// Fleet (Phase 2)
	mux.Handle("POST /api/v1/{organization_id}/fleet/vehicles", chain(fleetH.CreateVehicle))
	mux.Handle("GET /api/v1/{organization_id}/fleet/vehicles", chain(fleetH.ListVehicles))
	mux.Handle("GET /api/v1/{organization_id}/fleet/vehicles/{id}", chain(fleetH.GetVehicleByID))
	mux.Handle("PUT /api/v1/{organization_id}/fleet/vehicles/{id}", chain(fleetH.UpdateVehicle))
	mux.Handle("DELETE /api/v1/{organization_id}/fleet/vehicles/{id}", chain(fleetH.DeleteVehicle))
	mux.Handle("POST /api/v1/{organization_id}/fleet/maintenance", chain(fleetH.CreateMaintenance))
	mux.Handle("GET /api/v1/{organization_id}/fleet/maintenance", chain(fleetH.ListMaintenance))
	mux.Handle("GET /api/v1/{organization_id}/fleet/maintenance/{id}", chain(fleetH.GetMaintenanceByID))
	mux.Handle("DELETE /api/v1/{organization_id}/fleet/maintenance/{id}", chain(fleetH.DeleteMaintenance))

	// Labor — Teams, Workers, Shifts (Phase 2)
	mux.Handle("POST /api/v1/{organization_id}/labor/teams", chain(laborH.CreateTeam))
	mux.Handle("GET /api/v1/{organization_id}/labor/teams", chain(laborH.ListTeams))
	mux.Handle("PUT /api/v1/{organization_id}/labor/teams/{id}", chain(laborH.UpdateTeam))
	mux.Handle("DELETE /api/v1/{organization_id}/labor/teams/{id}", chain(laborH.DeleteTeam))
	mux.Handle("POST /api/v1/{organization_id}/labor/workers", chain(laborH.CreateWorker))
	mux.Handle("GET /api/v1/{organization_id}/labor/workers", chain(laborH.ListWorkers))
	mux.Handle("PUT /api/v1/{organization_id}/labor/workers/{id}", chain(laborH.UpdateWorker))
	mux.Handle("DELETE /api/v1/{organization_id}/labor/workers/{id}", chain(laborH.DeleteWorker))
	mux.Handle("POST /api/v1/{organization_id}/labor/shifts", chain(laborH.CreateWorkShift))
	mux.Handle("GET /api/v1/{organization_id}/labor/shifts", chain(laborH.ListWorkShifts))
	mux.Handle("DELETE /api/v1/{organization_id}/labor/shifts/{id}", chain(laborH.DeleteWorkShift))

	// Admin — Organization management (platform_owner only)
	adminChain := func(h http.HandlerFunc) http.Handler {
		return authMw(adminMw(http.HandlerFunc(h)))
	}

	mux.Handle("GET /api/v1/admin/organizations", adminChain(organizationH.List))
	mux.Handle("POST /api/v1/admin/organizations", adminChain(organizationH.Create))
	mux.Handle("GET /api/v1/admin/organizations/{id}", adminChain(organizationH.GetByID))
	mux.Handle("PUT /api/v1/admin/organizations/{id}", adminChain(organizationH.Update))
	mux.Handle("DELETE /api/v1/admin/organizations/{id}", adminChain(organizationH.Delete))

	// Admin — User management (platform_owner only)
	mux.Handle("GET /api/v1/admin/users", adminChain(userH.List))
	mux.Handle("POST /api/v1/admin/users", adminChain(userH.Create))
	mux.Handle("PUT /api/v1/admin/users/{id}", adminChain(userH.Update))
	mux.Handle("DELETE /api/v1/admin/users/{id}", adminChain(userH.Delete))

	// Cost Centers
	mux.Handle("POST /api/v1/{organization_id}/cost-centers", chain(ccH.Create))
	mux.Handle("GET /api/v1/{organization_id}/cost-centers", chain(ccH.List))
	mux.Handle("GET /api/v1/{organization_id}/cost-centers/senar-categories", chain(ccH.SenarCategories))
	mux.Handle("GET /api/v1/{organization_id}/cost-centers/{id}", chain(ccH.GetByID))
	mux.Handle("PUT /api/v1/{organization_id}/cost-centers/{id}", chain(ccH.Update))
	mux.Handle("DELETE /api/v1/{organization_id}/cost-centers/{id}", chain(ccH.Delete))

	// Budgets
	mux.Handle("POST /api/v1/{organization_id}/budgets", chain(budgetH.Create))
	mux.Handle("GET /api/v1/{organization_id}/harvests/{harvest_id}/budgets", chain(budgetH.ListByHarvest))
	mux.Handle("GET /api/v1/{organization_id}/budgets/{id}", chain(budgetH.GetByID))
	mux.Handle("PUT /api/v1/{organization_id}/budgets/{id}", chain(budgetH.Update))
	mux.Handle("DELETE /api/v1/{organization_id}/budgets/{id}", chain(budgetH.Delete))

	// Cost Allocations
	mux.Handle("POST /api/v1/{organization_id}/cost-allocations", chain(allocH.Create))
	mux.Handle("GET /api/v1/{organization_id}/harvests/{harvest_id}/cost-allocations", chain(allocH.ListByHarvest))
	mux.Handle("GET /api/v1/{organization_id}/cost-allocations/{id}", chain(allocH.GetByID))
	mux.Handle("DELETE /api/v1/{organization_id}/cost-allocations/{id}", chain(allocH.Delete))

	// Alerts
	mux.Handle("GET /api/v1/{organization_id}/alerts", chain(alertH.List))
	mux.Handle("PUT /api/v1/{organization_id}/alerts/{id}", chain(alertH.Update))

	// Sync — offline mobile (requires RabbitMQ publisher)
	if publisher != nil {
		syncH := handler.NewSyncHandler(publisher)
		mux.Handle("POST /api/v1/{organization_id}/sync", chain(syncH.Sync))
	}

	loggingMw := middleware.RequestLogging(log)
	return loggingMw(corsMw(mux))
}
