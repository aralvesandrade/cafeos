package main

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	"github.com/aralvesandrade/cafeos/internal/infra/config"
	"github.com/aralvesandrade/cafeos/internal/infra/db/postgres"
	infraLogger "github.com/aralvesandrade/cafeos/internal/infra/logger"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	adminEmail = "admin@cafeos.com.br"
	adminPass  = "admin123"
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

func seed(db *gorm.DB) error {
	// Tenant
	tenant := entity.Tenant{
		Name:         "CafeOS Padrão",
		Slug:         "cafeos",
		BrandName:    "CafeOS",
		Plan:         "pro",
		PrimaryColor: "#2E7D32",
	}
	if err := db.Create(&tenant).Error; err != nil {
		return fmt.Errorf("create tenant: %w", err)
	}
	fmt.Println("  ✓ Tenant: CafeOS Padrão (cafeos)")

	// Users
	users := []entity.User{
		{TenantID: tenant.ID, Name: "Administrador", Email: adminEmail, PasswordHash: hash(adminPass), Role: entity.RolePlatformOwner, IsActive: true},
		{TenantID: tenant.ID, Name: "João Silva", Email: "joao@cafeos.com.br", PasswordHash: hash("123456"), Role: entity.RoleProprietario, IsActive: true},
		{TenantID: tenant.ID, Name: "Maria Oliveira", Email: "maria@cafeos.com.br", PasswordHash: hash("123456"), Role: entity.RoleGerente, IsActive: true},
		{TenantID: tenant.ID, Name: "Carlos Santos", Email: "carlos@cafeos.com.br", PasswordHash: hash("123456"), Role: entity.RoleEngenheiro, IsActive: true},
		{TenantID: tenant.ID, Name: "Ana Costa", Email: "ana@cafeos.com.br", PasswordHash: hash("123456"), Role: entity.RoleOperador, IsActive: true},
	}
	for _, u := range users {
		if err := db.Create(&u).Error; err != nil {
			return fmt.Errorf("create user %s: %w", u.Email, err)
		}
	}
	fmt.Printf("  ✓ Usuários: %d criados\n", len(users))

	// Agricultural Products
	products := []entity.AgriculturalProduct{
		{TenantID: tenant.ID, Name: "NPK 20-05-20", Type: entity.ProdFertilizante, Unit: "kg"},
		{TenantID: tenant.ID, Name: "Calcário Dolomítico", Type: entity.ProdFertilizante, Unit: "kg"},
		{TenantID: tenant.ID, Name: "Glyphosate", Type: entity.ProdDefensivo, Unit: "l"},
		{TenantID: tenant.ID, Name: "Óleo Mineral", Type: entity.ProdDefensivo, Unit: "l"},
		{TenantID: tenant.ID, Name: "Óleo Diesel", Type: entity.ProdCombustivel, Unit: "l"},
	}
	for _, p := range products {
		if err := db.Create(&p).Error; err != nil {
			return fmt.Errorf("create product %s: %w", p.Name, err)
		}
	}
	fmt.Printf("  ✓ Produtos Agrícolas: %d criados\n", len(products))

	// Farms
	farms := []entity.Farm{
		{TenantID: tenant.ID, Name: "Fazenda Recanto Verde", Owner: "João Silva", Location: "Alfenas - MG", TotalAreaHA: 120, PlantedAreaHA: 95},
		{TenantID: tenant.ID, Name: "Sítio Boa Esperança", Owner: "João Silva", Location: "Machado - MG", TotalAreaHA: 45, PlantedAreaHA: 40},
		{TenantID: tenant.ID, Name: "Fazenda Monte Alegre", Owner: "Maria Oliveira", Location: "Poços de Caldas - MG", TotalAreaHA: 200, PlantedAreaHA: 160},
	}
	for i := range farms {
		if err := db.Create(&farms[i]).Error; err != nil {
			return fmt.Errorf("create farm %s: %w", farms[i].Name, err)
		}
	}
	fmt.Printf("  ✓ Fazendas: %d criadas\n", len(farms))

	// Plots
	plots := []entity.Plot{
		{TenantID: tenant.ID, FarmID: farms[0].ID, Name: "Talhão A-1", AreaHA: 30, Cultivar: "Catuaí Vermelho", PlantingYear: 2018, Altitude: 950, SoilType: "argiloso"},
		{TenantID: tenant.ID, FarmID: farms[0].ID, Name: "Talhão A-2", AreaHA: 35, Cultivar: "Mundo Novo", PlantingYear: 2019, Altitude: 920, SoilType: "argiloso"},
		{TenantID: tenant.ID, FarmID: farms[0].ID, Name: "Talhão B-1", AreaHA: 30, Cultivar: "Catuaí Amarelo", PlantingYear: 2020, Altitude: 980, SoilType: "arenoso"},
		{TenantID: tenant.ID, FarmID: farms[1].ID, Name: "Talhão Único", AreaHA: 40, Cultivar: "Bourbon Amarelo", PlantingYear: 2017, Altitude: 1050, SoilType: "organico"},
		{TenantID: tenant.ID, FarmID: farms[2].ID, Name: "Talhão Sul", AreaHA: 80, Cultivar: "Catuaí Vermelho", PlantingYear: 2016, Altitude: 1100, SoilType: "argiloso"},
		{TenantID: tenant.ID, FarmID: farms[2].ID, Name: "Talhão Norte", AreaHA: 80, Cultivar: "Acauã", PlantingYear: 2018, Altitude: 1080, SoilType: "siltoso"},
	}
	for i := range plots {
		if err := db.Create(&plots[i]).Error; err != nil {
			return fmt.Errorf("create plot %s: %w", plots[i].Name, err)
		}
	}
	fmt.Printf("  ✓ Talhões: %d criados\n", len(plots))

	// Operations
	now := time.Now()
	operations := []entity.Operation{
		{TenantID: tenant.ID, PlotID: plots[0].ID, Type: entity.OpAdubacao, Date: now.AddDate(0, -3, 0), Responsible: "Ana Costa", ProductUsed: "NPK 20-05-20", Quantity: 600, Cost: 4800, Notes: "Adubação de cobertura"},
		{TenantID: tenant.ID, PlotID: plots[1].ID, Type: entity.OpAdubacao, Date: now.AddDate(0, -3, -2), Responsible: "Ana Costa", ProductUsed: "NPK 20-05-20", Quantity: 700, Cost: 5600, Notes: ""},
		{TenantID: tenant.ID, PlotID: plots[0].ID, Type: entity.OpPulverizacao, Date: now.AddDate(0, -2, 0), Responsible: "Carlos Santos", ProductUsed: "Glyphosate", Quantity: 30, Cost: 1200, Notes: "Controle de plantas daninhas"},
		{TenantID: tenant.ID, PlotID: plots[3].ID, Type: entity.OpPoda, Date: now.AddDate(0, -1, -15), Responsible: "Maria Oliveira", ProductUsed: "", Quantity: 0, Cost: 2500, Notes: "Poda de formação"},
		{TenantID: tenant.ID, PlotID: plots[2].ID, Type: entity.OpIrrigacao, Date: now.AddDate(0, -1, -5), Responsible: "Ana Costa", ProductUsed: "", Quantity: 0, Cost: 1800, Notes: "Irrigação de salvamento"},
		{TenantID: tenant.ID, PlotID: plots[4].ID, Type: entity.OpPulverizacao, Date: now.AddDate(0, 0, -10), Responsible: "Carlos Santos", ProductUsed: "Óleo Mineral", Quantity: 50, Cost: 2000, Notes: "Controle de ácaros"},
		{TenantID: tenant.ID, PlotID: plots[0].ID, Type: entity.OpAdubacao, Date: now.AddDate(0, 0, -5), Responsible: "Ana Costa", ProductUsed: "Calcário Dolomítico", Quantity: 1500, Cost: 3000, Notes: "Calagem"},
	}
	for i := range operations {
		if err := db.Create(&operations[i]).Error; err != nil {
			return fmt.Errorf("create operation: %w", err)
		}
	}
	fmt.Printf("  ✓ Operações: %d criadas\n", len(operations))

	// Harvests
	harvests := []entity.Harvest{
		{TenantID: tenant.ID, Year: 2024, Description: "Safra 2024", EstimatedProduction: 3000, Status: entity.HarvestFinalizada},
		{TenantID: tenant.ID, Year: 2025, Description: "Safra 2025", EstimatedProduction: 3500, Status: entity.HarvestEmAndamento},
	}
	for i := range harvests {
		if err := db.Create(&harvests[i]).Error; err != nil {
			return fmt.Errorf("create harvest: %w", err)
		}
	}
	fmt.Printf("  ✓ Safras: %d criadas\n", len(harvests))

	// Harvest Productions
	productions := []entity.HarvestProduction{
		{TenantID: tenant.ID, HarvestID: harvests[0].ID, PlotID: plots[0].ID, Quantity: 800, RecordedAt: now.AddDate(0, -6, 0), Notes: "Lote 1"},
		{TenantID: tenant.ID, HarvestID: harvests[0].ID, PlotID: plots[1].ID, Quantity: 950, RecordedAt: now.AddDate(0, -6, -5), Notes: "Lote 2"},
		{TenantID: tenant.ID, HarvestID: harvests[0].ID, PlotID: plots[2].ID, Quantity: 750, RecordedAt: now.AddDate(0, -5, -10), Notes: "Lote 3"},
		{TenantID: tenant.ID, HarvestID: harvests[0].ID, PlotID: plots[3].ID, Quantity: 1100, RecordedAt: now.AddDate(0, -5, -15), Notes: ""},
		{TenantID: tenant.ID, HarvestID: harvests[1].ID, PlotID: plots[0].ID, Quantity: 400, RecordedAt: now.AddDate(0, -1, 0), Notes: "Parcial 1"},
	}
	for i := range productions {
		if err := db.Create(&productions[i]).Error; err != nil {
			return fmt.Errorf("create production: %w", err)
		}
	}
	fmt.Printf("  ✓ Produções: %d criadas\n", len(productions))

	// Indicators for finalized harvest
	indicators := []entity.Indicator{
		{TenantID: tenant.ID, HarvestID: harvests[0].ID, Type: entity.IndProducaoTotal, Value: 3600, CalculatedAt: now.AddDate(0, -4, 0)},
		{TenantID: tenant.ID, HarvestID: harvests[0].ID, Type: entity.IndCustoTotal, Value: 45000, CalculatedAt: now.AddDate(0, -4, 0)},
		{TenantID: tenant.ID, HarvestID: harvests[0].ID, Type: entity.IndSacasHA, Value: 37.89, CalculatedAt: now.AddDate(0, -4, 0)},
		{TenantID: tenant.ID, HarvestID: harvests[0].ID, Type: entity.IndCustoSaca, Value: 12.50, CalculatedAt: now.AddDate(0, -4, 0)},
		{TenantID: tenant.ID, HarvestID: harvests[0].ID, Type: entity.IndRentabilidade, Value: 85000, CalculatedAt: now.AddDate(0, -4, 0)},
	}
	for i := range indicators {
		if err := db.Create(&indicators[i]).Error; err != nil {
			return fmt.Errorf("create indicator: %w", err)
		}
	}
	fmt.Printf("  ✓ Indicadores: %d criados\n", len(indicators))

	return nil
}
