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

func NewRouter(db *gorm.DB, eventBus event.Bus, publisher *messaging.Publisher, jwtSecret string, signupOrgSlug string, log *slog.Logger) http.Handler {
	mux := http.NewServeMux()

	farmRepo := infraRepo.NewFarmRepository(db)
	producerRepo := infraRepo.NewProducerRepository(db)
	plotRepo := infraRepo.NewPlotRepository(db)
	opRepo := infraRepo.NewOperationRepository(db)
	harvestRepo := infraRepo.NewHarvestRepository(db)
	hpRepo := infraRepo.NewHarvestProductionRepository(db)
	indicatorRepo := infraRepo.NewIndicatorRepository(db)
	organizationRepo := infraRepo.NewOrganizationRepository(db)
	planRepo := infraRepo.NewPlanRepository(db)
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
	permRepo := infraRepo.NewPermissionRepository(db)
	roleRepo := infraRepo.NewRoleRepository(db)
	moduleRepo := infraRepo.NewModuleRepository(db)
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
	roleSvc := domainSvc.NewRoleService(roleRepo, userRepo)
	moduleSvc := domainSvc.NewModuleService(moduleRepo)
	permSvc := domainSvc.NewPermissionService(permRepo, roleRepo)
	planSvc := domainSvc.NewPlanService(planRepo)

	if err := roleSvc.SeedDefaultsIfMissing(); err != nil {
		log.Error("failed to seed roles", "error", err)
	}
	if err := moduleSvc.SeedDefaultsIfMissing(); err != nil {
		log.Error("failed to seed modules", "error", err)
	}
	if organizations, err := organizationRepo.List(); err != nil {
		log.Error("failed to list organizations for permission backfill", "error", err)
	} else {
		for _, org := range organizations {
			if err := permSvc.SeedDefaultsIfMissing(org.ID); err != nil {
				log.Error("failed to backfill default permissions", "organization_id", org.ID, "error", err)
			}
		}
	}

	userSvc := domainSvc.NewUserService(userRepo, planRepo)
	signupSvc := domainSvc.NewSignupService(transactor, organizationRepo, planRepo, roleSvc, signupOrgSlug)
	farmH := handler.NewFarmHandler(farmSvc, userSvc, roleSvc)
	plotH := handler.NewPlotHandler(plotSvc, farmSvc)
	opH := handler.NewOperationHandler(opSvc, plotSvc, farmSvc)
	harvestH := handler.NewHarvestHandler(harvestSvc, plotSvc, farmSvc, indicatorRepo)
	dashboardH := handler.NewDashboardHandler(harvestRepo, indicatorRepo, opRepo, plotRepo, farmRepo, hpRepo, farmSvc)
	authH := handler.NewAuthHandler(userRepo, organizationRepo, jwtSecret)
	signupH := handler.NewSignupHandler(signupSvc)
	organizationH := handler.NewOrganizationHandler(organizationRepo, permSvc)
	planH := handler.NewPlanHandler(planSvc)
	userH := handler.NewUserHandler(userRepo, roleSvc, userSvc)
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
	permH := handler.NewPermissionHandler(permSvc)
	roleH := handler.NewRoleHandler(roleSvc)
	moduleH := handler.NewModuleHandler(moduleSvc)

	authMw := middleware.Auth(jwtSecret)
	adminMw := middleware.RequireRole(entity.SystemRolePlatformOwner)
	corsMw := middleware.CORS

	// Auth (public)
	mux.HandleFunc("POST /auth/login", authH.Login)
	mux.HandleFunc("POST /auth/register", signupH.Register)

	// Health
	mux.HandleFunc("GET /health", handler.HealthCheck)

	// Public plans (landing page)
	mux.HandleFunc("GET /api/v1/public/plans", planH.ListPublic)

	// Swagger UI
	mux.Handle("GET /swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	// API v1 routes
	chain := func(h http.HandlerFunc) http.Handler {
		return authMw(http.HandlerFunc(h))
	}
	mchain := func(module entity.ModuleKey, access entity.AccessLevel, h http.HandlerFunc) http.Handler {
		return authMw(middleware.RequireModule(permSvc, module, access)(http.HandlerFunc(h)))
	}
	read := entity.AccessRead
	write := entity.AccessWrite

	// Farms
	mux.Handle("POST /api/v1/{organization_id}/farms", mchain(entity.ModuleFarms, write, farmH.Create))
	mux.Handle("GET /api/v1/{organization_id}/farms", mchain(entity.ModuleFarms, read, farmH.List))
	mux.Handle("GET /api/v1/{organization_id}/farms/{id}", mchain(entity.ModuleFarms, read, farmH.GetByID))
	mux.Handle("PUT /api/v1/{organization_id}/farms/{id}", mchain(entity.ModuleFarms, write, farmH.Update))
	mux.Handle("DELETE /api/v1/{organization_id}/farms/{id}", mchain(entity.ModuleFarms, write, farmH.Delete))

	// Plots
	mux.Handle("POST /api/v1/{organization_id}/plots", mchain(entity.ModuleFarms, write, plotH.Create))
	mux.Handle("GET /api/v1/{organization_id}/plots", mchain(entity.ModuleFarms, read, plotH.List))
	mux.Handle("GET /api/v1/{organization_id}/plots/{id}", mchain(entity.ModuleFarms, read, plotH.GetByID))
	mux.Handle("PUT /api/v1/{organization_id}/plots/{id}", mchain(entity.ModuleFarms, write, plotH.Update))
	mux.Handle("DELETE /api/v1/{organization_id}/plots/{id}", mchain(entity.ModuleFarms, write, plotH.Delete))
	mux.Handle("GET /api/v1/{organization_id}/farms/{farm_id}/plots", mchain(entity.ModuleFarms, read, plotH.ListByFarm))

	// Operations
	mux.Handle("POST /api/v1/{organization_id}/operations", mchain(entity.ModuleOperations, write, opH.Create))
	mux.Handle("GET /api/v1/{organization_id}/operations", mchain(entity.ModuleOperations, read, opH.List))
	mux.Handle("GET /api/v1/{organization_id}/operations/recent", mchain(entity.ModuleOperations, read, opH.ListRecent))
	mux.Handle("GET /api/v1/{organization_id}/operations/{id}", mchain(entity.ModuleOperations, read, opH.GetByID))
	mux.Handle("PUT /api/v1/{organization_id}/operations/{id}", mchain(entity.ModuleOperations, write, opH.Update))
	mux.Handle("DELETE /api/v1/{organization_id}/operations/{id}", mchain(entity.ModuleOperations, write, opH.Delete))
	mux.Handle("GET /api/v1/{organization_id}/plots/{plot_id}/operations", mchain(entity.ModuleOperations, read, opH.ListByPlot))

	// Operation Types
	mux.Handle("POST /api/v1/{organization_id}/operation-types", mchain(entity.ModuleFarms, write, otH.Create))
	mux.Handle("GET /api/v1/{organization_id}/operation-types", mchain(entity.ModuleFarms, read, otH.List))
	mux.Handle("GET /api/v1/{organization_id}/operation-types/{id}", mchain(entity.ModuleFarms, read, otH.GetByID))
	mux.Handle("PUT /api/v1/{organization_id}/operation-types/{id}", mchain(entity.ModuleFarms, write, otH.Update))
	mux.Handle("DELETE /api/v1/{organization_id}/operation-types/{id}", mchain(entity.ModuleFarms, write, otH.Delete))

	// Harvests
	mux.Handle("POST /api/v1/{organization_id}/harvests", mchain(entity.ModuleHarvests, write, harvestH.Create))
	mux.Handle("GET /api/v1/{organization_id}/harvests", mchain(entity.ModuleHarvests, read, harvestH.List))
	mux.Handle("GET /api/v1/{organization_id}/harvests/{id}", mchain(entity.ModuleHarvests, read, harvestH.GetByID))
	mux.Handle("PUT /api/v1/{organization_id}/harvests/{id}/finalize", mchain(entity.ModuleHarvests, write, harvestH.Finalize))
	mux.Handle("POST /api/v1/{organization_id}/harvests/{id}/production", mchain(entity.ModuleHarvests, write, harvestH.RecordProduction))
	mux.Handle("GET /api/v1/{organization_id}/harvests/{id}/production", mchain(entity.ModuleHarvests, read, harvestH.GetProduction))
	mux.Handle("GET /api/v1/{organization_id}/harvests/{id}/indicators", mchain(entity.ModuleHarvests, read, harvestH.GetIndicators))

	// Dashboard
	mux.Handle("GET /api/v1/{organization_id}/dashboard", mchain(entity.ModuleDashboard, read, dashboardH.GetDashboard))

	// Agricultural Products
	mux.Handle("GET /api/v1/{organization_id}/agricultural-products", mchain(entity.ModuleFarms, read, prodH.List))

	// Financial (Phase 2)
	mux.Handle("POST /api/v1/{organization_id}/financial", mchain(entity.ModuleFinancial, write, finH.Create))
	mux.Handle("GET /api/v1/{organization_id}/financial", mchain(entity.ModuleFinancial, read, finH.List))
	mux.Handle("GET /api/v1/{organization_id}/financial/{id}", mchain(entity.ModuleFinancial, read, finH.GetByID))
	mux.Handle("PUT /api/v1/{organization_id}/financial/{id}", mchain(entity.ModuleFinancial, write, finH.Update))
	mux.Handle("DELETE /api/v1/{organization_id}/financial/{id}", mchain(entity.ModuleFinancial, write, finH.Delete))

	// Stock (Phase 2)
	mux.Handle("POST /api/v1/{organization_id}/stock/items", mchain(entity.ModuleResources, write, stockH.CreateItem))
	mux.Handle("GET /api/v1/{organization_id}/stock/items", mchain(entity.ModuleResources, read, stockH.ListItems))
	mux.Handle("GET /api/v1/{organization_id}/stock/items/{id}", mchain(entity.ModuleResources, read, stockH.GetItemByID))
	mux.Handle("PUT /api/v1/{organization_id}/stock/items/{id}", mchain(entity.ModuleResources, write, stockH.UpdateItem))
	mux.Handle("DELETE /api/v1/{organization_id}/stock/items/{id}", mchain(entity.ModuleResources, write, stockH.DeleteItem))
	mux.Handle("POST /api/v1/{organization_id}/stock/movements", mchain(entity.ModuleResources, write, stockH.RecordMovement))
	mux.Handle("GET /api/v1/{organization_id}/stock/movements", mchain(entity.ModuleResources, read, stockH.ListMovements))

	// Fleet (Phase 2)
	mux.Handle("POST /api/v1/{organization_id}/fleet/vehicles", mchain(entity.ModuleResources, write, fleetH.CreateVehicle))
	mux.Handle("GET /api/v1/{organization_id}/fleet/vehicles", mchain(entity.ModuleResources, read, fleetH.ListVehicles))
	mux.Handle("GET /api/v1/{organization_id}/fleet/vehicles/{id}", mchain(entity.ModuleResources, read, fleetH.GetVehicleByID))
	mux.Handle("PUT /api/v1/{organization_id}/fleet/vehicles/{id}", mchain(entity.ModuleResources, write, fleetH.UpdateVehicle))
	mux.Handle("DELETE /api/v1/{organization_id}/fleet/vehicles/{id}", mchain(entity.ModuleResources, write, fleetH.DeleteVehicle))
	mux.Handle("POST /api/v1/{organization_id}/fleet/maintenance", mchain(entity.ModuleResources, write, fleetH.CreateMaintenance))
	mux.Handle("GET /api/v1/{organization_id}/fleet/maintenance", mchain(entity.ModuleResources, read, fleetH.ListMaintenance))
	mux.Handle("GET /api/v1/{organization_id}/fleet/maintenance/{id}", mchain(entity.ModuleResources, read, fleetH.GetMaintenanceByID))
	mux.Handle("DELETE /api/v1/{organization_id}/fleet/maintenance/{id}", mchain(entity.ModuleResources, write, fleetH.DeleteMaintenance))

	// Labor — Teams, Workers, Shifts (Phase 2)
	mux.Handle("POST /api/v1/{organization_id}/labor/teams", mchain(entity.ModuleResources, write, laborH.CreateTeam))
	mux.Handle("GET /api/v1/{organization_id}/labor/teams", mchain(entity.ModuleResources, read, laborH.ListTeams))
	mux.Handle("PUT /api/v1/{organization_id}/labor/teams/{id}", mchain(entity.ModuleResources, write, laborH.UpdateTeam))
	mux.Handle("DELETE /api/v1/{organization_id}/labor/teams/{id}", mchain(entity.ModuleResources, write, laborH.DeleteTeam))
	mux.Handle("POST /api/v1/{organization_id}/labor/workers", mchain(entity.ModuleResources, write, laborH.CreateWorker))
	mux.Handle("GET /api/v1/{organization_id}/labor/workers", mchain(entity.ModuleResources, read, laborH.ListWorkers))
	mux.Handle("PUT /api/v1/{organization_id}/labor/workers/{id}", mchain(entity.ModuleResources, write, laborH.UpdateWorker))
	mux.Handle("DELETE /api/v1/{organization_id}/labor/workers/{id}", mchain(entity.ModuleResources, write, laborH.DeleteWorker))
	mux.Handle("POST /api/v1/{organization_id}/labor/shifts", mchain(entity.ModuleResources, write, laborH.CreateWorkShift))
	mux.Handle("GET /api/v1/{organization_id}/labor/shifts", mchain(entity.ModuleResources, read, laborH.ListWorkShifts))
	mux.Handle("DELETE /api/v1/{organization_id}/labor/shifts/{id}", mchain(entity.ModuleResources, write, laborH.DeleteWorkShift))

	// Admin — Organization management (platform_owner only)
	adminChain := func(h http.HandlerFunc) http.Handler {
		return authMw(adminMw(http.HandlerFunc(h)))
	}

	mux.Handle("GET /api/v1/admin/organizations", adminChain(organizationH.List))
	mux.Handle("POST /api/v1/admin/organizations", adminChain(organizationH.Create))
	mux.Handle("GET /api/v1/admin/organizations/{id}", adminChain(organizationH.GetByID))
	mux.Handle("PUT /api/v1/admin/organizations/{id}", adminChain(organizationH.Update))
	mux.Handle("DELETE /api/v1/admin/organizations/{id}", adminChain(organizationH.Delete))

	mux.Handle("GET /api/v1/admin/plans", adminChain(planH.List))
	mux.Handle("POST /api/v1/admin/plans", adminChain(planH.Create))
	mux.Handle("GET /api/v1/admin/plans/{id}", adminChain(planH.GetByID))
	mux.Handle("PUT /api/v1/admin/plans/{id}", adminChain(planH.Update))
	mux.Handle("DELETE /api/v1/admin/plans/{id}", adminChain(planH.Delete))

	// Admin — User management (platform_owner only)
	mux.Handle("GET /api/v1/admin/users", adminChain(userH.List))
	mux.Handle("POST /api/v1/admin/users", adminChain(userH.Create))
	mux.Handle("PUT /api/v1/admin/users/{id}", adminChain(userH.Update))
	mux.Handle("DELETE /api/v1/admin/users/{id}", adminChain(userH.Delete))

	// Cost Centers
	mux.Handle("POST /api/v1/{organization_id}/cost-centers", mchain(entity.ModuleFinancial, write, ccH.Create))
	mux.Handle("GET /api/v1/{organization_id}/cost-centers", mchain(entity.ModuleFinancial, read, ccH.List))
	mux.Handle("GET /api/v1/{organization_id}/cost-centers/senar-categories", mchain(entity.ModuleFinancial, read, ccH.SenarCategories))
	mux.Handle("GET /api/v1/{organization_id}/cost-centers/{id}", mchain(entity.ModuleFinancial, read, ccH.GetByID))
	mux.Handle("PUT /api/v1/{organization_id}/cost-centers/{id}", mchain(entity.ModuleFinancial, write, ccH.Update))
	mux.Handle("DELETE /api/v1/{organization_id}/cost-centers/{id}", mchain(entity.ModuleFinancial, write, ccH.Delete))

	// Budgets
	mux.Handle("POST /api/v1/{organization_id}/budgets", mchain(entity.ModuleFinancial, write, budgetH.Create))
	mux.Handle("GET /api/v1/{organization_id}/harvests/{harvest_id}/budgets", mchain(entity.ModuleFinancial, read, budgetH.ListByHarvest))
	mux.Handle("GET /api/v1/{organization_id}/budgets/{id}", mchain(entity.ModuleFinancial, read, budgetH.GetByID))
	mux.Handle("PUT /api/v1/{organization_id}/budgets/{id}", mchain(entity.ModuleFinancial, write, budgetH.Update))
	mux.Handle("DELETE /api/v1/{organization_id}/budgets/{id}", mchain(entity.ModuleFinancial, write, budgetH.Delete))

	// Cost Allocations
	mux.Handle("POST /api/v1/{organization_id}/cost-allocations", mchain(entity.ModuleFinancial, write, allocH.Create))
	mux.Handle("GET /api/v1/{organization_id}/harvests/{harvest_id}/cost-allocations", mchain(entity.ModuleFinancial, read, allocH.ListByHarvest))
	mux.Handle("GET /api/v1/{organization_id}/cost-allocations/{id}", mchain(entity.ModuleFinancial, read, allocH.GetByID))
	mux.Handle("DELETE /api/v1/{organization_id}/cost-allocations/{id}", mchain(entity.ModuleFinancial, write, allocH.Delete))

	// Alerts
	mux.Handle("GET /api/v1/{organization_id}/alerts", mchain(entity.ModuleDashboard, read, alertH.List))
	mux.Handle("PUT /api/v1/{organization_id}/alerts/{id}", mchain(entity.ModuleDashboard, write, alertH.Update))

	// Permissions
	mux.Handle("GET /api/v1/{organization_id}/permissions", mchain(entity.ModulePermissions, write, permH.List))
	mux.Handle("PUT /api/v1/{organization_id}/permissions", mchain(entity.ModulePermissions, write, permH.Update))
	mux.Handle("GET /api/v1/{organization_id}/permissions/me", chain(permH.Mine))

	// Roles
	mux.Handle("GET /api/v1/{organization_id}/roles", mchain(entity.ModulePermissions, write, roleH.List))
	mux.Handle("POST /api/v1/{organization_id}/roles", mchain(entity.ModulePermissions, write, roleH.Create))
	mux.Handle("PUT /api/v1/{organization_id}/roles/{id}", mchain(entity.ModulePermissions, write, roleH.Update))
	mux.Handle("DELETE /api/v1/{organization_id}/roles/{id}", mchain(entity.ModulePermissions, write, roleH.Delete))

	// Modules (catalog — display metadata only, keys are fixed in code)
	mux.Handle("GET /api/v1/{organization_id}/modules", mchain(entity.ModulePermissions, write, moduleH.List))
	mux.Handle("PUT /api/v1/{organization_id}/modules/{key}", mchain(entity.ModulePermissions, write, moduleH.Update))

	// Users — org-scoped roster (distinct from the cross-tenant /admin/users)
	mux.Handle("POST /api/v1/{organization_id}/users", mchain(entity.ModuleUsers, write, userH.CreateForOrg))
	mux.Handle("GET /api/v1/{organization_id}/users", mchain(entity.ModuleUsers, read, userH.ListForOrg))
	mux.Handle("PUT /api/v1/{organization_id}/users/{id}", mchain(entity.ModuleUsers, write, userH.UpdateForOrg))
	mux.Handle("DELETE /api/v1/{organization_id}/users/{id}", mchain(entity.ModuleUsers, write, userH.DeleteForOrg))
	mux.Handle("PATCH /api/v1/{organization_id}/users/{id}/plan", authMw(http.HandlerFunc(userH.AssignPlanForOrg)))

	// Sync — offline mobile (requires RabbitMQ publisher)
	if publisher != nil {
		syncH := handler.NewSyncHandler(publisher)
		mux.Handle("POST /api/v1/{organization_id}/sync", chain(syncH.Sync))
	}

	loggingMw := middleware.RequestLogging(log)
	return loggingMw(corsMw(mux))
}
