package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	domainSvc "github.com/aralvesandrade/cafeos/internal/domain/service"
	"github.com/aralvesandrade/cafeos/internal/infra/config"
	"github.com/aralvesandrade/cafeos/internal/infra/db/postgres"
	infraRepo "github.com/aralvesandrade/cafeos/internal/infra/db/repository"
	infraLogger "github.com/aralvesandrade/cafeos/internal/infra/logger"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	adminEmail  = "admin@cafeos.com.br"
	adminPass   = "admin123"
	defaultPass = "123456"
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

	if err := seed(db); err != nil {
		log.Error("seed failed", "error", err)
		os.Exit(1)
	}

	fmt.Println("Seed concluído com sucesso!")
}

func hash(password string) string {
	h, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(h)
}

// Phase 1 — system catalog. Idempotent via SeedDefaultsIfMissing.
// Phase 2 — tenant + demo data. Each section checks if data exists
// for the tenant before inserting, making the seed re-runnable.
func seed(db *gorm.DB) error {
	// ── Phase 1: System Catalog ──────────────────────────────────────

	// Plans
	var planCount int64
	planIDBySlug := make(map[string]string)
	db.Model(&entity.Plan{}).Count(&planCount)
	if planCount == 0 {
		plans := []entity.Plan{
			{
				Name: "Essencial", Slug: "essencial", Description: "Para pequenos produtores que estão começando",
				PriceCents: 2900, BillingInterval: "monthly", MaxFarms: 1, MaxUsers: 2,
				Features: entity.PlanFeatureList{
					{Label: "1 fazenda", Included: true},
					{Label: "5 talhões", Included: true},
					{Label: "50 operações/mês", Included: true},
					{Label: "3 safras", Included: true},
					{Label: "2 usuários", Included: true},
					{Label: "Dashboard básico", Included: true},
					{Label: "Relatórios CSV/PDF", Included: false},
					{Label: "White Label", Included: false},
					{Label: "Suporte prioritário", Included: false},
				},
				Active: true, Featured: false, DisplayOrder: 1,
			},
			{
				Name: "Pro", Slug: "pro", Description: "Para médios produtores com operações mais complexas",
				PriceCents: 9700, BillingInterval: "monthly", MaxFarms: 10, MaxUsers: 10,
				Features: entity.PlanFeatureList{
					{Label: "10 fazendas", Included: true},
					{Label: "50 talhões", Included: true},
					{Label: "Operações ilimitadas", Included: true},
					{Label: "Safras ilimitadas", Included: true},
					{Label: "10 usuários", Included: true},
					{Label: "Dashboard avançado", Included: true},
					{Label: "Relatórios CSV/PDF", Included: true},
					{Label: "White Label", Included: false},
					{Label: "Suporte por email", Included: true},
				},
				Active: true, Featured: true, DisplayOrder: 2,
			},
			{
				Name: "Cooperativa", Slug: "cooperativa", Description: "Para cooperativas que gerenciam múltiplos associados",
				PriceCents: 29700, BillingInterval: "monthly", MaxFarms: 0, MaxUsers: 50,
				Features: entity.PlanFeatureList{
					{Label: "Fazendas ilimitadas", Included: true},
					{Label: "Talhões ilimitados", Included: true},
					{Label: "Operações ilimitadas", Included: true},
					{Label: "Safras ilimitadas", Included: true},
					{Label: "50 usuários", Included: true},
					{Label: "Dashboard consolidado", Included: true},
					{Label: "Benchmarking associados", Included: true},
					{Label: "White Label", Included: false},
					{Label: "Suporte prioritário", Included: true},
				},
				Active: true, Featured: false, DisplayOrder: 3,
			},
			{
				Name: "Consultoria", Slug: "consultoria", Description: "Para consultorias que atendem múltiplos clientes",
				PriceCents: 49700, BillingInterval: "monthly", MaxFarms: 0, MaxUsers: 30,
				Features: entity.PlanFeatureList{
					{Label: "Multi-cliente", Included: true},
					{Label: "Talhões ilimitados", Included: true},
					{Label: "Operações ilimitadas", Included: true},
					{Label: "Safras ilimitadas", Included: true},
					{Label: "30 usuários", Included: true},
					{Label: "Dashboard por cliente", Included: true},
					{Label: "Relatórios técnicos", Included: true},
					{Label: "White Label", Included: true},
					{Label: "Suporte dedicado", Included: true},
				},
				Active: true, Featured: false, DisplayOrder: 4,
			},
		}
		for i := range plans {
			if err := db.Create(&plans[i]).Error; err != nil {
				return fmt.Errorf("create plan %s: %w", plans[i].Slug, err)
			}
			planIDBySlug[plans[i].Slug] = plans[i].ID
		}
		fmt.Println("  ✓ Plans: essencial, pro, cooperativa, consultoria")
	} else {
		// Load existing plans for ID references
		var existingPlans []entity.Plan
		db.Find(&existingPlans)
		planIDBySlug = make(map[string]string, len(existingPlans))
		for _, p := range existingPlans {
			planIDBySlug[p.Slug] = p.ID
		}
	}

	// Organizations
	organization := entity.Organization{}
	if err := db.Where("slug = ?", "cafeos").First(&organization).Error; err != nil {
		organization = entity.Organization{
			Name:         "CafeOS Padrão",
			Slug:         "cafeos",
			BrandName:    "CafeOS",
			PrimaryColor: "#2E7D32",
		}
		if err := db.Create(&organization).Error; err != nil {
			return fmt.Errorf("create organization: %w", err)
		}
		fmt.Println("  ✓ Organization: CafeOS Padrão (cafeos)")
	} else {
		fmt.Println("  ✓ Organization: CafeOS Padrão (cafeos) — already exists")
	}

	// Modules
	moduleSvc := domainSvc.NewModuleService(infraRepo.NewModuleRepository(db))
	if err := moduleSvc.SeedDefaultsIfMissing(); err != nil {
		return fmt.Errorf("seed modules: %w", err)
	}
	fmt.Println("  ✓ Módulos seedados")

	// Rules
	roleRepo := infraRepo.NewRoleRepository(db)
	userRepo := infraRepo.NewUserRepository(db)
	roleSvc := domainSvc.NewRoleService(roleRepo, userRepo)
	if err := roleSvc.SeedDefaultsIfMissing(); err != nil {
		return fmt.Errorf("seed roles: %w", err)
	}
	fmt.Println("  ✓ Papéis padrão seedados")

	// Role Permissions
	permSvc := domainSvc.NewPermissionService(infraRepo.NewPermissionRepository(db), roleRepo)
	if err := permSvc.SeedDefaultsIfMissing(organization.ID); err != nil {
		return fmt.Errorf("seed default permissions: %w", err)
	}
	fmt.Println("  ✓ Permissões padrão seedadas")

	// ── Phase 2: Tenant + Demo Data ──────────────────────────────────

	roleIDByKey := make(map[string]string)
	allRoles, err := roleRepo.List()
	if err != nil {
		return fmt.Errorf("list roles: %w", err)
	}
	for _, role := range allRoles {
		roleIDByKey[role.Key] = role.ID
	}

	// now := time.Now()

	// ── Users (idempotent: first-or-create by email) ──────────────
	type userDef struct {
		Name, Email, Password, RoleKey string
		PlanSlug                       string
	}
	userDefs := []userDef{
		// SystemRoleKeys
		{Name: "Administrador", Email: adminEmail, Password: adminPass, RoleKey: entity.SystemRolePlatformOwner},
		{Name: "Fernanda Lima", Email: "fernanda@cafeos.com.br", Password: defaultPass, RoleKey: entity.SystemRoleOrganizationAdmin},
		// DefaultRoleKeys with Plan
		{Name: "João Silva", Email: "joao@cafeos.com.br", Password: defaultPass, RoleKey: entity.RoleKeyProprietario, PlanSlug: "pro"},
		{Name: "Maria Oliveira", Email: "maria@cafeos.com.br", Password: defaultPass, RoleKey: entity.RoleKeyProprietario, PlanSlug: "essencial"},
		// DefaultRoleKeys
		{Name: "Carlos Santos", Email: "carlos@cafeos.com.br", Password: defaultPass, RoleKey: "engenheiro_agronomo"},
		{Name: "Ana Costa", Email: "ana@cafeos.com.br", Password: defaultPass, RoleKey: "operador_campo"},
		{Name: "Rodrigo Alves", Email: "rodrigo@cafeos.com.br", Password: defaultPass, RoleKey: "consultor_externo"},
	}
	userByEmail := make(map[string]*entity.User, len(userDefs))
	for _, d := range userDefs {
		u := entity.User{}
		if err := db.Where("organization_id = ? AND email = ?", organization.ID, d.Email).First(&u).Error; err != nil {
			u = entity.User{
				OrganizationID: organization.ID,
				Name:           d.Name,
				Email:          d.Email,
				PasswordHash:   hash(d.Password),
				RoleID:         roleIDByKey[d.RoleKey],
				IsActive:       true,
			}
			if d.PlanSlug != "" {
				if pid, ok := planIDBySlug[d.PlanSlug]; ok {
					u.PlanID = &pid
				}
			}
			if err := db.Create(&u).Error; err != nil {
				return fmt.Errorf("create user %s: %w", d.Email, err)
			}
		}
		userByEmail[d.Email] = &u
	}
	fmt.Printf("  ✓ Usuários: %d\n", len(userByEmail))

	joaoUser := userByEmail["joao@cafeos.com.br"]
	mariaUser := userByEmail["maria@cafeos.com.br"]
	// anaUser := userByEmail["ana@cafeos.com.br"]

	// Managed-by hierarchy
	type mgmt struct {
		UserEmail string
		ManagerID string
	}
	mgmtChain := []mgmt{
		{UserEmail: "carlos@cafeos.com.br", ManagerID: joaoUser.ID},
		{UserEmail: "ana@cafeos.com.br", ManagerID: joaoUser.ID},
		{UserEmail: "rodrigo@cafeos.com.br", ManagerID: mariaUser.ID},
	}
	for _, m := range mgmtChain {
		db.Model(&entity.User{}).Where("organization_id = ? AND email = ?", organization.ID, m.UserEmail).Update("managed_by_user_id", m.ManagerID)
	}
	fmt.Println("  ✓ Hierarquia principal/sub-usuário vinculada")

	// ── Agricultural Products (idempotent) ─────────────────────────
	// products := make([]entity.AgriculturalProduct, 0)
	// var prodCount int64
	// db.Model(&entity.AgriculturalProduct{}).Where("organization_id = ?", organization.ID).Count(&prodCount)
	// if prodCount == 0 {
	// 	productDefs := []entity.AgriculturalProduct{
	// 		{OrganizationID: organization.ID, Name: "NPK 20-05-20", Type: entity.ProdFertilizante, Unit: "kg"},
	// 		{OrganizationID: organization.ID, Name: "Calcário Dolomítico", Type: entity.ProdFertilizante, Unit: "kg"},
	// 		{OrganizationID: organization.ID, Name: "Glyphosate", Type: entity.ProdDefensivo, Unit: "l"},
	// 		{OrganizationID: organization.ID, Name: "Óleo Mineral", Type: entity.ProdDefensivo, Unit: "l"},
	// 		{OrganizationID: organization.ID, Name: "Óleo Diesel", Type: entity.ProdCombustivel, Unit: "l"},
	// 	}
	// 	for i := range productDefs {
	// 		if err := db.Create(&productDefs[i]).Error; err != nil {
	// 			return fmt.Errorf("create product %s: %w", productDefs[i].Name, err)
	// 		}
	// 		products = append(products, productDefs[i])
	// 	}
	// 	fmt.Printf("  ✓ Produtos Agrícolas: %d\n", len(products))
	// } else {
	// 	db.Where("organization_id = ?", organization.ID).Find(&products)
	// 	fmt.Printf("  ✓ Produtos Agrícolas: %d (loaded)\n", len(products))
	// }

	// ── Farms + Plots + Operation Types (idempotent) ──────────────
	// farms := make([]entity.Farm, 0)
	// plots := make([]entity.Plot, 0)
	operationTypes := make([]entity.OperationType, 0)
	// ccIDByCode := make(map[string]string)

	// db.Where("organization_id = ?", organization.ID).Find(&farms)
	// if len(farms) == 0 {
	// 	farmDefs := []entity.Farm{
	// 		{OrganizationID: organization.ID, Name: "Fazenda Recanto Verde", Owner: "João Silva", Location: "Alfenas - MG", TotalAreaHA: 120, PlantedAreaHA: 95},
	// 		{OrganizationID: organization.ID, Name: "Sítio Boa Esperança", Owner: "João Silva", Location: "Machado - MG", TotalAreaHA: 45, PlantedAreaHA: 40},
	// 		{OrganizationID: organization.ID, Name: "Fazenda Monte Alegre", Owner: "Maria Oliveira", Location: "Poços de Caldas - MG", TotalAreaHA: 200, PlantedAreaHA: 160},
	// 	}
	// 	for i := range farmDefs {
	// 		if err := db.Create(&farmDefs[i]).Error; err != nil {
	// 			return fmt.Errorf("create farm %s: %w", farmDefs[i].Name, err)
	// 		}
	// 		farms = append(farms, farmDefs[i])
	// 	}
	// 	fmt.Printf("  ✓ Fazendas: %d\n", len(farms))
	// } else {
	// 	fmt.Printf("  ✓ Fazendas: %d (loaded)\n", len(farms))
	// }

	// Producers (recria se vazio — vincula User↔Farm)
	// var prodExists int64
	// db.Model(&entity.Producer{}).Where("organization_id = ?", organization.ID).Count(&prodExists)
	// if prodExists == 0 {
	// 	producers := []entity.Producer{
	// 		{OrganizationID: organization.ID, FarmID: farms[0].ID, UserID: joaoUser.ID, RoleID: roleIDByKey[entity.RoleKeyProprietario], CPF: "123.456.789-00", Name: "João Silva", Phone: "(35) 99999-0001", Email: "joao@cafeos.com.br"},
	// 		{OrganizationID: organization.ID, FarmID: farms[0].ID, UserID: anaUser.ID, RoleID: roleIDByKey["operador_campo"], Name: "Ana Costa", Phone: "(35) 99999-0004", Email: "ana@cafeos.com.br"},
	// 		{OrganizationID: organization.ID, FarmID: farms[1].ID, UserID: anaUser.ID, RoleID: roleIDByKey["operador_campo"], Name: "Ana Costa", Phone: "(35) 99999-0004", Email: "ana@cafeos.com.br"},
	// 		{OrganizationID: organization.ID, FarmID: farms[1].ID, UserID: joaoUser.ID, RoleID: roleIDByKey[entity.RoleKeyProprietario], CPF: "123.456.789-00", Name: "João Silva", Phone: "(35) 99999-0001", Email: "joao@cafeos.com.br"},
	// 		{OrganizationID: organization.ID, FarmID: farms[2].ID, UserID: mariaUser.ID, RoleID: roleIDByKey[entity.RoleKeyProprietario], CPF: "987.654.321-00", Name: "Maria Oliveira", Phone: "(35) 99999-0002", Email: "maria@cafeos.com.br"},
	// 	}
	// 	for i := range producers {
	// 		if err := db.Create(&producers[i]).Error; err != nil {
	// 			return fmt.Errorf("create producer: %w", err)
	// 		}
	// 	}
	// 	fmt.Printf("  ✓ Produtores: %d\n", len(producers))
	// } else {
	// 	fmt.Printf("  ✓ Produtores: %d (loaded)\n", prodExists)
	// }

	// db.Where("organization_id = ?", organization.ID).Find(&plots)
	// if len(plots) == 0 {
	// 	plotDefs := []entity.Plot{
	// 		{OrganizationID: organization.ID, FarmID: farms[0].ID, Name: "Talhão A-1", AreaHA: 30, Cultivar: "Catuaí Vermelho", PlantingYear: 2018, Altitude: 950, SoilType: "argiloso"},
	// 		{OrganizationID: organization.ID, FarmID: farms[0].ID, Name: "Talhão A-2", AreaHA: 35, Cultivar: "Mundo Novo", PlantingYear: 2019, Altitude: 920, SoilType: "argiloso"},
	// 		{OrganizationID: organization.ID, FarmID: farms[0].ID, Name: "Talhão B-1", AreaHA: 30, Cultivar: "Catuaí Amarelo", PlantingYear: 2020, Altitude: 980, SoilType: "arenoso"},
	// 		{OrganizationID: organization.ID, FarmID: farms[1].ID, Name: "Talhão Único", AreaHA: 40, Cultivar: "Bourbon Amarelo", PlantingYear: 2017, Altitude: 1050, SoilType: "organico"},
	// 		{OrganizationID: organization.ID, FarmID: farms[2].ID, Name: "Talhão Sul", AreaHA: 80, Cultivar: "Catuaí Vermelho", PlantingYear: 2016, Altitude: 1100, SoilType: "argiloso"},
	// 		{OrganizationID: organization.ID, FarmID: farms[2].ID, Name: "Talhão Norte", AreaHA: 80, Cultivar: "Acauã", PlantingYear: 2018, Altitude: 1080, SoilType: "siltoso"},
	// 	}
	// 	for i := range plotDefs {
	// 		if err := db.Create(&plotDefs[i]).Error; err != nil {
	// 			return fmt.Errorf("create plot %s: %w", plotDefs[i].Name, err)
	// 		}
	// 		plots = append(plots, plotDefs[i])
	// 	}
	// 	fmt.Printf("  ✓ Talhões: %d\n", len(plots))
	// } else {
	// 	fmt.Printf("  ✓ Talhões: %d (loaded)\n", len(plots))
	// }

	// Operation Types — codes are the unique key per org
	{
		otDefs := []entity.OperationType{
			{OrganizationID: organization.ID, Name: "Adubação", Code: "adubacao", Color: "info"},
			{OrganizationID: organization.ID, Name: "Pulverização", Code: "pulverizacao", Color: "warning"},
			{OrganizationID: organization.ID, Name: "Irrigação", Code: "irrigacao", Color: "success"},
			{OrganizationID: organization.ID, Name: "Poda", Code: "poda", Color: "default"},
			{OrganizationID: organization.ID, Name: "Colheita", Code: "colheita", Color: "danger"},
		}
		for i := range otDefs {
			existing := entity.OperationType{}
			if err := db.Where("organization_id = ? AND code = ?", organization.ID, otDefs[i].Code).First(&existing).Error; err != nil {
				if err := db.Create(&otDefs[i]).Error; err != nil {
					return fmt.Errorf("create operation type %s: %w", otDefs[i].Name, err)
				}
				operationTypes = append(operationTypes, otDefs[i])
			} else {
				operationTypes = append(operationTypes, existing)
			}
		}
		otMap := make(map[string]*entity.OperationType, len(operationTypes))
		for i := range operationTypes {
			otMap[operationTypes[i].Code] = &operationTypes[i]
		}
		_ = otMap
		fmt.Printf("  ✓ Tipos de Operação: %d\n", len(operationTypes))
	}

	// ── Operations ────────────────────────────────────────────────
	// var opCount int64
	// db.Model(&entity.Operation{}).Where("organization_id = ?", organization.ID).Count(&opCount)
	// if opCount == 0 {
	// 	ccAdubos := ""
	// 	ccDefensivos := ""
	// 	ccMaoObra := ""
	// 	ccIrrigacao := ""
	// 	{
	// 		var cc entity.CostCenter
	// 		if err := db.Where("organization_id = ? AND code = ?", organization.ID, "DESP_ADUBOS").First(&cc).Error; err == nil {
	// 			ccAdubos = cc.ID
	// 		}
	// 		if err := db.Where("organization_id = ? AND code = ?", organization.ID, "DESP_DEFENSIVOS").First(&cc).Error; err == nil {
	// 			ccDefensivos = cc.ID
	// 		}
	// 		if err := db.Where("organization_id = ? AND code = ?", organization.ID, "DESP_MAO_OBRA").First(&cc).Error; err == nil {
	// 			ccMaoObra = cc.ID
	// 		}
	// 		if err := db.Where("organization_id = ? AND code = ?", organization.ID, "DESP_IRRIGACAO").First(&cc).Error; err == nil {
	// 			ccIrrigacao = cc.ID
	// 		}
	// 	}
	// 	ccAdubosP := &ccAdubos
	// 	if ccAdubos == "" {
	// 		ccAdubosP = nil
	// 	}
	// 	ccDefensivosP := &ccDefensivos
	// 	if ccDefensivos == "" {
	// 		ccDefensivosP = nil
	// 	}
	// 	ccMaoObraP := &ccMaoObra
	// 	if ccMaoObra == "" {
	// 		ccMaoObraP = nil
	// 	}
	// 	ccIrrigacaoP := &ccIrrigacao
	// 	if ccIrrigacao == "" {
	// 		ccIrrigacaoP = nil
	// 	}

	// 	operations := []entity.Operation{
	// 		{OrganizationID: organization.ID, PlotID: plots[0].ID, TypeID: operationTypes[0].ID, CostCenterID: ccAdubosP, Date: now.AddDate(0, -3, 0), Responsible: "Ana Costa", ProductUsed: "NPK 20-05-20", Quantity: 600, Cost: 4800, Notes: "Adubação de cobertura"},
	// 		{OrganizationID: organization.ID, PlotID: plots[1].ID, TypeID: operationTypes[0].ID, CostCenterID: ccAdubosP, Date: now.AddDate(0, -3, -2), Responsible: "Ana Costa", ProductUsed: "NPK 20-05-20", Quantity: 700, Cost: 5600, Notes: ""},
	// 		{OrganizationID: organization.ID, PlotID: plots[0].ID, TypeID: operationTypes[1].ID, CostCenterID: ccDefensivosP, Date: now.AddDate(0, -2, 0), Responsible: "Carlos Santos", ProductUsed: "Glyphosate", Quantity: 30, Cost: 1200, Notes: "Controle de plantas daninhas"},
	// 		{OrganizationID: organization.ID, PlotID: plots[3].ID, TypeID: operationTypes[3].ID, CostCenterID: ccMaoObraP, Date: now.AddDate(0, -1, -15), Responsible: "Maria Oliveira", ProductUsed: "", Quantity: 0, Cost: 2500, Notes: "Poda de formação"},
	// 		{OrganizationID: organization.ID, PlotID: plots[2].ID, TypeID: operationTypes[2].ID, CostCenterID: ccIrrigacaoP, Date: now.AddDate(0, -1, -5), Responsible: "Ana Costa", ProductUsed: "", Quantity: 0, Cost: 1800, Notes: "Irrigação de salvamento"},
	// 		{OrganizationID: organization.ID, PlotID: plots[4].ID, TypeID: operationTypes[1].ID, CostCenterID: ccDefensivosP, Date: now.AddDate(0, 0, -10), Responsible: "Carlos Santos", ProductUsed: "Óleo Mineral", Quantity: 50, Cost: 2000, Notes: "Controle de ácaros"},
	// 		{OrganizationID: organization.ID, PlotID: plots[0].ID, TypeID: operationTypes[0].ID, CostCenterID: ccAdubosP, Date: now.AddDate(0, 0, -5), Responsible: "Ana Costa", ProductUsed: "Calcário Dolomítico", Quantity: 1500, Cost: 3000, Notes: "Calagem"},
	// 	}
	// 	for i := range operations {
	// 		if err := db.Create(&operations[i]).Error; err != nil {
	// 			return fmt.Errorf("create operation: %w", err)
	// 		}
	// 	}
	// 	fmt.Printf("  ✓ Operações: %d\n", len(operations))
	// } else {
	// 	fmt.Printf("  ✓ Operações: %d (loaded)\n", opCount)
	// }

	// ── Harvests (idempotent by org+year) ─────────────────────────
	// harvests := make([]entity.Harvest, 0)
	// for _, y := range []int{2024, 2025} {
	// 	h := entity.Harvest{}
	// 	if err := db.Where("organization_id = ? AND year = ?", organization.ID, y).First(&h).Error; err != nil {
	// 		status := entity.HarvestPlanejada
	// 		if y == 2024 {
	// 			status = entity.HarvestFinalizada
	// 		} else if y == 2025 {
	// 			status = entity.HarvestEmAndamento
	// 		}
	// 		h = entity.Harvest{
	// 			OrganizationID:      organization.ID,
	// 			Year:                y,
	// 			Description:         fmt.Sprintf("Safra %d", y),
	// 			EstimatedProduction: 3000 + float64(y-2024)*500,
	// 			Status:              status,
	// 		}
	// 		if err := db.Create(&h).Error; err != nil {
	// 			return fmt.Errorf("create harvest %d: %w", y, err)
	// 		}
	// 	}
	// 	harvests = append(harvests, h)
	// }
	// fmt.Printf("  ✓ Safras: %d\n", len(harvests))

	// ── Harvest Productions ───────────────────────────────────────
	// var hpCount int64
	// db.Model(&entity.HarvestProduction{}).Where("organization_id = ?", organization.ID).Count(&hpCount)
	// if hpCount == 0 {
	// 	productions := []entity.HarvestProduction{
	// 		{OrganizationID: organization.ID, HarvestID: harvests[0].ID, PlotID: plots[0].ID, Quantity: 800, RecordedAt: now.AddDate(0, -6, 0), Notes: "Lote 1"},
	// 		{OrganizationID: organization.ID, HarvestID: harvests[0].ID, PlotID: plots[1].ID, Quantity: 950, RecordedAt: now.AddDate(0, -6, -5), Notes: "Lote 2"},
	// 		{OrganizationID: organization.ID, HarvestID: harvests[0].ID, PlotID: plots[2].ID, Quantity: 750, RecordedAt: now.AddDate(0, -5, -10), Notes: "Lote 3"},
	// 		{OrganizationID: organization.ID, HarvestID: harvests[0].ID, PlotID: plots[3].ID, Quantity: 1100, RecordedAt: now.AddDate(0, -5, -15), Notes: ""},
	// 		{OrganizationID: organization.ID, HarvestID: harvests[1].ID, PlotID: plots[0].ID, Quantity: 400, RecordedAt: now.AddDate(0, -1, 0), Notes: "Parcial 1"},
	// 	}
	// 	for i := range productions {
	// 		if err := db.Create(&productions[i]).Error; err != nil {
	// 			return fmt.Errorf("create harvest production: %w", err)
	// 		}
	// 	}
	// 	fmt.Printf("  ✓ Produções: %d\n", len(productions))
	// } else {
	// 	fmt.Printf("  ✓ Produções: %d (loaded)\n", hpCount)
	// }

	// ── Cost Centers (idempotent by org+code) ─────────────────────
	// {
	// 	ccDefs := []entity.CostCenter{
	// 		{OrganizationID: organization.ID, Name: "Adubos", Code: "DESP_ADUBOS", Type: entity.CCDespesa, Description: "Fertilizantes e corretivos"},
	// 		{OrganizationID: organization.ID, Name: "Defensivos", Code: "DESP_DEFENSIVOS", Type: entity.CCDespesa, Description: "Agrotóxicos e defensivos agrícolas"},
	// 		{OrganizationID: organization.ID, Name: "Combustíveis", Code: "DESP_COMBUSTIVEL", Type: entity.CCDespesa, Description: "Diesel, gasolina e lubrificantes"},
	// 		{OrganizationID: organization.ID, Name: "Mão de Obra", Code: "DESP_MAO_OBRA", Type: entity.CCDespesa, Description: "Salários e encargos trabalhistas"},
	// 		{OrganizationID: organization.ID, Name: "Frete", Code: "DESP_FRETE", Type: entity.CCDespesa, Description: "Transporte de insumos e produção"},
	// 		{OrganizationID: organization.ID, Name: "Manutenção", Code: "DESP_MANUTENCAO", Type: entity.CCDespesa, Description: "Manutenção de máquinas e equipamentos"},
	// 		{OrganizationID: organization.ID, Name: "Irrigação", Code: "DESP_IRRIGACAO", Type: entity.CCDespesa, Description: "Custo de irrigação e água"},
	// 		{OrganizationID: organization.ID, Name: "Análise de Solo", Code: "DESP_ANALISE_SOLO", Type: entity.CCDespesa, Description: "Análises laboratoriais de solo e folha"},
	// 		{OrganizationID: organization.ID, Name: "Outros Insumos", Code: "DESP_OUTROS_INSUMOS", Type: entity.CCDespesa, Description: "Outros insumos agrícolas"},
	// 		{OrganizationID: organization.ID, Name: "Serviços Terceiros", Code: "DESP_SERV_TERCEIROS", Type: entity.CCDespesa, Description: "Serviços contratados de terceiros"},
	// 		{OrganizationID: organization.ID, Name: "Energia", Code: "DESP_ENERGIA", Type: entity.CCDespesa, Description: "Energia elétrica"},
	// 		{OrganizationID: organization.ID, Name: "Depreciação", Code: "DESP_DEPRECIACAO", Type: entity.CCDespesa, Description: "Depreciação de máquinas e benfeitorias"},
	// 		{OrganizationID: organization.ID, Name: "Administrativo", Code: "DESP_ADMINISTRATIVO", Type: entity.CCDespesa, Description: "Despesas administrativas"},
	// 		{OrganizationID: organization.ID, Name: "Outras Despesas", Code: "DESP_OUTRAS", Type: entity.CCDespesa, Description: "Outras despesas operacionais"},
	// 		{OrganizationID: organization.ID, Name: "Venda de Café", Code: "REC_CAFE", Type: entity.CCReceita, Description: "Receita com venda de café"},
	// 		{OrganizationID: organization.ID, Name: "Venda de Mudas", Code: "REC_MUDAS", Type: entity.CCReceita, Description: "Receita com venda de mudas"},
	// 		{OrganizationID: organization.ID, Name: "Outras Receitas", Code: "REC_OUTRAS", Type: entity.CCReceita, Description: "Outras receitas"},
	// 	}
	// 	for i := range ccDefs {
	// 		cc := entity.CostCenter{}
	// 		if err := db.Where("organization_id = ? AND code = ?", organization.ID, ccDefs[i].Code).First(&cc).Error; err != nil {
	// 			if err := db.Create(&ccDefs[i]).Error; err != nil {
	// 				return fmt.Errorf("create cost center %s: %w", ccDefs[i].Name, err)
	// 			}
	// 			cc = ccDefs[i]
	// 		}
	// 		ccIDByCode[cc.Code] = cc.ID
	// 	}
	// 	fmt.Printf("  ✓ Centros de Custo: %d\n", len(ccDefs))
	// }

	// ── Budgets ───────────────────────────────────────────────────
	// var budgetCount int64
	// db.Model(&entity.Budget{}).Where("organization_id = ?", organization.ID).Count(&budgetCount)
	// if budgetCount == 0 {
	// 	budgets := []entity.Budget{
	// 		{OrganizationID: organization.ID, HarvestID: harvests[0].ID, CostCenterID: ccIDByCode["DESP_ADUBOS"], PlannedAmount: 15000, Description: "Orçamento adubos safra 2024"},
	// 		{OrganizationID: organization.ID, HarvestID: harvests[0].ID, CostCenterID: ccIDByCode["DESP_MAO_OBRA"], PlannedAmount: 25000, Description: "Orçamento mão de obra safra 2024"},
	// 		{OrganizationID: organization.ID, HarvestID: harvests[1].ID, CostCenterID: ccIDByCode["DESP_ADUBOS"], PlannedAmount: 18000, Description: "Orçamento adubos safra 2025"},
	// 		{OrganizationID: organization.ID, HarvestID: harvests[1].ID, CostCenterID: ccIDByCode["DESP_COMBUSTIVEL"], PlannedAmount: 8000, Description: "Orçamento combustível safra 2025"},
	// 	}
	// 	for i := range budgets {
	// 		if err := db.Create(&budgets[i]).Error; err != nil {
	// 			return fmt.Errorf("create budget: %w", err)
	// 		}
	// 	}
	// 	fmt.Printf("  ✓ Orçamentos: %d\n", len(budgets))
	// } else {
	// 	fmt.Printf("  ✓ Orçamentos: %d (loaded)\n", budgetCount)
	// }

	// ── Cost Allocations ──────────────────────────────────────────
	// var allocCount int64
	// db.Model(&entity.CostAllocation{}).Where("organization_id = ?", organization.ID).Count(&allocCount)
	// if allocCount == 0 {
	// 	alloc := entity.CostAllocation{
	// 		OrganizationID: organization.ID,
	// 		HarvestID:      harvests[0].ID,
	// 		CostCenterID:   ccIDByCode["DESP_ADUBOS"],
	// 		Description:    "Rateio adubos safra 2024 — área proporcional",
	// 		TotalAmount:    10000,
	// 		Method:         entity.AllocAreaProportional,
	// 		Date:           now.AddDate(0, -4, 0),
	// 	}
	// 	if err := db.Create(&alloc).Error; err != nil {
	// 		return fmt.Errorf("create cost allocation: %w", err)
	// 	}
	// 	items := []entity.CostAllocationItem{
	// 		{AllocationID: alloc.ID, PlotID: plots[0].ID, Amount: 3158, Percentage: 31.58},
	// 		{AllocationID: alloc.ID, PlotID: plots[1].ID, Amount: 3684, Percentage: 36.84},
	// 		{AllocationID: alloc.ID, PlotID: plots[2].ID, Amount: 3158, Percentage: 31.58},
	// 	}
	// 	for i := range items {
	// 		if err := db.Create(&items[i]).Error; err != nil {
	// 			return fmt.Errorf("create cost allocation item: %w", err)
	// 		}
	// 	}
	// 	fmt.Printf("  ✓ Rateios: 1 (3 itens)\n")
	// } else {
	// 	fmt.Printf("  ✓ Rateios: %d (loaded)\n", allocCount)
	// }

	// ── Indicators ────────────────────────────────────────────────
	// var indCount int64
	// db.Model(&entity.Indicator{}).Where("organization_id = ?", organization.ID).Count(&indCount)
	// if indCount == 0 {
	// 	indicators := []entity.Indicator{
	// 		{OrganizationID: organization.ID, HarvestID: harvests[0].ID, Type: entity.IndProducaoTotal, Value: 3600, CalculatedAt: now.AddDate(0, -4, 0)},
	// 		{OrganizationID: organization.ID, HarvestID: harvests[0].ID, Type: entity.IndCustoTotal, Value: 45000, CalculatedAt: now.AddDate(0, -4, 0)},
	// 		{OrganizationID: organization.ID, HarvestID: harvests[0].ID, Type: entity.IndSacasHA, Value: 37.89, CalculatedAt: now.AddDate(0, -4, 0)},
	// 		{OrganizationID: organization.ID, HarvestID: harvests[0].ID, Type: entity.IndCustoSaca, Value: 12.50, CalculatedAt: now.AddDate(0, -4, 0)},
	// 		{OrganizationID: organization.ID, HarvestID: harvests[0].ID, Type: entity.IndCOE, Value: 30000, CalculatedAt: now.AddDate(0, -4, 0)},
	// 		{OrganizationID: organization.ID, HarvestID: harvests[0].ID, Type: entity.IndCOT, Value: 40000, CalculatedAt: now.AddDate(0, -4, 0)},
	// 		{OrganizationID: organization.ID, HarvestID: harvests[0].ID, Type: entity.IndCTProducao, Value: 45000, CalculatedAt: now.AddDate(0, -4, 0)},
	// 	}
	// 	for i := range indicators {
	// 		if err := db.Create(&indicators[i]).Error; err != nil {
	// 			return fmt.Errorf("create indicator: %w", err)
	// 		}
	// 	}
	// 	fmt.Printf("  ✓ Indicadores: %d\n", len(indicators))
	// } else {
	// 	fmt.Printf("  ✓ Indicadores: %d (loaded)\n", indCount)
	// }

	// ── Financial ─────────────────────────────────────────────────
	// var finCount int64
	// db.Model(&entity.FinancialTransaction{}).Where("organization_id = ?", organization.ID).Count(&finCount)
	// if finCount == 0 {
	// 	financials := []entity.FinancialTransaction{
	// 		{OrganizationID: organization.ID, FarmID: &farms[0].ID, Type: entity.TranDespesa, Description: "Adubação Recanto Verde", Amount: 4800, Date: now.AddDate(0, -3, 0), DueDate: now.AddDate(0, -3, 5), Status: "paid"},
	// 		{OrganizationID: organization.ID, FarmID: &farms[2].ID, Type: entity.TranReceita, Description: "Venda de café - Monte Alegre", Amount: 52000, Date: now.AddDate(0, -1, 0), DueDate: now.AddDate(0, -1, 0), Status: "paid"},
	// 		{OrganizationID: organization.ID, Type: entity.TranDespesa, Description: "Contabilidade (organização)", Amount: 1200, Date: now, DueDate: now.AddDate(0, 0, 10), Status: "pending"},
	// 	}
	// 	for i := range financials {
	// 		if err := db.Create(&financials[i]).Error; err != nil {
	// 			return fmt.Errorf("create financial transaction: %w", err)
	// 		}
	// 	}
	// 	fmt.Printf("  ✓ Transações Financeiras: %d\n", len(financials))
	// } else {
	// 	fmt.Printf("  ✓ Transações Financeiras: %d (loaded)\n", finCount)
	// }

	// ── Stock ─────────────────────────────────────────────────────
	// stockItems := make([]entity.StockItem, 0)
	// db.Where("organization_id = ?", organization.ID).Find(&stockItems)
	// if len(stockItems) == 0 {
	// 	stockDefs := []entity.StockItem{
	// 		{OrganizationID: organization.ID, FarmID: &farms[0].ID, ProductID: products[0].ID, Quantity: 500, Unit: "kg", Location: "Galpão Recanto Verde"},
	// 		{OrganizationID: organization.ID, FarmID: &farms[2].ID, ProductID: products[1].ID, Quantity: 300, Unit: "kg", Location: "Galpão Monte Alegre"},
	// 		{OrganizationID: organization.ID, ProductID: products[4].ID, Quantity: 1000, Unit: "l", Location: "Depósito Central"},
	// 	}
	// 	for i := range stockDefs {
	// 		if err := db.Create(&stockDefs[i]).Error; err != nil {
	// 			return fmt.Errorf("create stock item: %w", err)
	// 		}
	// 		stockItems = append(stockItems, stockDefs[i])
	// 	}
	// 	fmt.Printf("  ✓ Itens de Estoque: %d\n", len(stockItems))
	// } else {
	// 	fmt.Printf("  ✓ Itens de Estoque: %d (loaded)\n", len(stockItems))
	// }

	// ── Stock Movements ───────────────────────────────────────────
	// var smCount int64
	// db.Model(&entity.StockMovement{}).Where("organization_id = ?", organization.ID).Count(&smCount)
	// if smCount == 0 && len(stockItems) >= 3 {
	// 	movements := []entity.StockMovement{
	// 		{OrganizationID: organization.ID, ItemID: stockItems[0].ID, Type: "in", Quantity: 200, Date: now.AddDate(0, -1, 0), Reference: "NF-12345", Notes: "Reposição NPK"},
	// 		{OrganizationID: organization.ID, ItemID: stockItems[2].ID, Type: "out", Quantity: 50, Date: now.AddDate(0, 0, -3), Reference: "Abastecimento trator", Notes: "Óleo diesel para operações"},
	// 	}
	// 	for i := range movements {
	// 		if err := db.Create(&movements[i]).Error; err != nil {
	// 			return fmt.Errorf("create stock movement: %w", err)
	// 		}
	// 	}
	// 	fmt.Printf("  ✓ Movimentações Estoque: %d\n", len(movements))
	// }

	// ── Fleet ─────────────────────────────────────────────────────
	// vehicles := make([]entity.Vehicle, 0)
	// db.Where("organization_id = ?", organization.ID).Find(&vehicles)
	// if len(vehicles) == 0 {
	// 	vehDefs := []entity.Vehicle{
	// 		{OrganizationID: organization.ID, FarmID: &farms[0].ID, Name: "Trator Recanto Verde", Type: entity.VeicTractor, Plate: "ABC1D23", Brand: "Massey Ferguson", Model: "4275", Year: 2019, Status: "active"},
	// 		{OrganizationID: organization.ID, FarmID: &farms[2].ID, Name: "Trator Monte Alegre", Type: entity.VeicTractor, Plate: "XYZ9E87", Brand: "New Holland", Model: "TL75E", Year: 2021, Status: "active"},
	// 		{OrganizationID: organization.ID, Name: "Caminhão da Cooperativa", Type: entity.VeicCaminhao, Plate: "QAZ2W34", Brand: "Volkswagen", Model: "Delivery", Year: 2018, Status: "active"},
	// 	}
	// 	for i := range vehDefs {
	// 		if err := db.Create(&vehDefs[i]).Error; err != nil {
	// 			return fmt.Errorf("create vehicle %s: %w", vehDefs[i].Name, err)
	// 		}
	// 		vehicles = append(vehicles, vehDefs[i])
	// 	}
	// 	fmt.Printf("  ✓ Veículos: %d\n", len(vehicles))
	// } else {
	// 	fmt.Printf("  ✓ Veículos: %d (loaded)\n", len(vehicles))
	// }

	// ── Maintenance ───────────────────────────────────────────────
	// var mtCount int64
	// db.Model(&entity.Maintenance{}).Where("organization_id = ?", organization.ID).Count(&mtCount)
	// if mtCount == 0 && len(vehicles) >= 3 {
	// 	maintenances := []entity.Maintenance{
	// 		{OrganizationID: organization.ID, VehicleID: vehicles[0].ID, CostCenterID: strPtr(ccIDByCode["DESP_MANUTENCAO"]), Date: now.AddDate(0, -2, 0), Type: "preventive", Description: "Revisão 5000h", Cost: 1200, Odometer: 5000, Notes: ""},
	// 		{OrganizationID: organization.ID, VehicleID: vehicles[2].ID, CostCenterID: strPtr(ccIDByCode["DESP_MANUTENCAO"]), Date: now.AddDate(0, 0, -15), Type: "corrective", Description: "Troca de embreagem", Cost: 3500, Odometer: 15000, Notes: "Caminhão apresentou ruído"},
	// 	}
	// 	for i := range maintenances {
	// 		if err := db.Create(&maintenances[i]).Error; err != nil {
	// 			return fmt.Errorf("create maintenance: %w", err)
	// 		}
	// 	}
	// 	fmt.Printf("  ✓ Manutenções: %d\n", len(maintenances))
	// }

	// ── Labor: Teams + Workers + Shifts ──────────────────────────
	// var teamCount int64
	// db.Model(&entity.Team{}).Where("organization_id = ?", organization.ID).Count(&teamCount)
	// if teamCount == 0 {
	// 	team := entity.Team{
	// 		OrganizationID: organization.ID,
	// 		Name:           "Equipe Recanto Verde",
	// 		Leader:         "Ana Costa",
	// 		Description:    "Equipe de campo responsável pelas operações na Fazenda Recanto Verde",
	// 	}
	// 	if err := db.Create(&team).Error; err != nil {
	// 		return fmt.Errorf("create team: %w", err)
	// 	}
	// 	workerAna := entity.Worker{
	// 		OrganizationID: organization.ID,
	// 		TeamID:         team.ID,
	// 		Name:           "Ana Costa",
	// 		Role:           "operador_campo",
	// 		Phone:          "(35) 99999-0004",
	// 		HourlyRate:     15,
	// 		IsActive:       true,
	// 	}
	// 	if err := db.Create(&workerAna).Error; err != nil {
	// 		return fmt.Errorf("create worker Ana: %w", err)
	// 	}
	// 	workerCarlos := entity.Worker{
	// 		OrganizationID: organization.ID,
	// 		TeamID:         team.ID,
	// 		Name:           "Carlos Santos",
	// 		Role:           "engenheiro_agronomo",
	// 		Phone:          "(35) 99999-0003",
	// 		HourlyRate:     25,
	// 		IsActive:       true,
	// 	}
	// 	if err := db.Create(&workerCarlos).Error; err != nil {
	// 		return fmt.Errorf("create worker Carlos: %w", err)
	// 	}
	// 	shifts := []entity.WorkShift{
	// 		{OrganizationID: organization.ID, WorkerID: workerAna.ID, CostCenterID: strPtr(ccIDByCode["DESP_MAO_OBRA"]), Date: now.AddDate(0, 0, -5), Hours: 8, Cost: 120, Notes: "Adubação talhão A-1"},
	// 		{OrganizationID: organization.ID, WorkerID: workerCarlos.ID, CostCenterID: strPtr(ccIDByCode["DESP_MAO_OBRA"]), Date: now.AddDate(0, 0, -5), Hours: 6, Cost: 150, Notes: "Vistoria técnica talhão Sul"},
	// 	}
	// 	for i := range shifts {
	// 		if err := db.Create(&shifts[i]).Error; err != nil {
	// 			return fmt.Errorf("create work shift: %w", err)
	// 		}
	// 	}
	// 	fmt.Println("  ✓ Equipes: 1, Trabalhadores: 2, Turnos: 2")
	// } else {
	// 	fmt.Printf("  ✓ Equipes: %d (loaded)\n", teamCount)
	// }

	// ── Alerts ────────────────────────────────────────────────────
	// var alertCount int64
	// db.Model(&entity.Alert{}).Where("organization_id = ?", organization.ID).Count(&alertCount)
	// if alertCount == 0 {
	// 	alerts := []entity.Alert{
	// 		{OrganizationID: organization.ID, HarvestID: harvests[0].ID, RuleID: "custo_total", Message: "Custo total da safra 2024 está 15% acima do orçado (R$ 45.000 vs R$ 40.000 orçados)", Severity: "warning", Status: "aberto"},
	// 		{OrganizationID: organization.ID, HarvestID: harvests[1].ID, RuleID: "produtividade", Message: "Produtividade estimada da safra 2025 está abaixo da safra anterior (35 sc/ha vs 37,9 sc/ha)", Severity: "info", Status: "aberto"},
	// 	}
	// 	for i := range alerts {
	// 		if err := db.Create(&alerts[i]).Error; err != nil {
	// 			return fmt.Errorf("create alert: %w", err)
	// 		}
	// 	}
	// 	fmt.Printf("  ✓ Alertas: %d\n", len(alerts))
	// }

	return nil
}

func strPtr(s string) *string { return &s }
