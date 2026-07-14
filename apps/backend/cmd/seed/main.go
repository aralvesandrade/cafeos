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
	// Organization
	organization := entity.Organization{
		Name:         "CafeOS Padrão",
		Slug:         "cafeos",
		BrandName:    "CafeOS",
		Plan:         "pro",
		PrimaryColor: "#2E7D32",
	}
	if err := db.Create(&organization).Error; err != nil {
		return fmt.Errorf("create organization: %w", err)
	}
	fmt.Println("  ✓ Organization: CafeOS Padrão (cafeos)")

	// Users
	users := []entity.User{
		{OrganizationID: organization.ID, Name: "Administrador", Email: adminEmail, PasswordHash: hash(adminPass), Role: entity.RolePlatformOwner, IsActive: true},
		{OrganizationID: organization.ID, Name: "João Silva", Email: "joao@cafeos.com.br", PasswordHash: hash("123456"), Role: entity.RoleProprietario, IsActive: true},
		{OrganizationID: organization.ID, Name: "Maria Oliveira", Email: "maria@cafeos.com.br", PasswordHash: hash("123456"), Role: entity.RoleProprietario, IsActive: true},
		{OrganizationID: organization.ID, Name: "Carlos Santos", Email: "carlos@cafeos.com.br", PasswordHash: hash("123456"), Role: entity.RoleEngenheiro, IsActive: true},
		{OrganizationID: organization.ID, Name: "Ana Costa", Email: "ana@cafeos.com.br", PasswordHash: hash("123456"), Role: entity.RoleOperador, IsActive: true},
		{OrganizationID: organization.ID, Name: "Fernanda Lima", Email: "fernanda@cafeos.com.br", PasswordHash: hash("123456"), Role: entity.RoleOrganizationAdmin, IsActive: true},
		{OrganizationID: organization.ID, Name: "Rodrigo Alves", Email: "rodrigo@cafeos.com.br", PasswordHash: hash("123456"), Role: entity.RoleConsultor, IsActive: true},
	}
	for i := range users {
		if err := db.Create(&users[i]).Error; err != nil {
			return fmt.Errorf("create user %s: %w", users[i].Email, err)
		}
	}
	fmt.Printf("  ✓ Usuários: %d criados\n", len(users))

	// Agricultural Products
	products := []entity.AgriculturalProduct{
		{OrganizationID: organization.ID, Name: "NPK 20-05-20", Type: entity.ProdFertilizante, Unit: "kg"},
		{OrganizationID: organization.ID, Name: "Calcário Dolomítico", Type: entity.ProdFertilizante, Unit: "kg"},
		{OrganizationID: organization.ID, Name: "Glyphosate", Type: entity.ProdDefensivo, Unit: "l"},
		{OrganizationID: organization.ID, Name: "Óleo Mineral", Type: entity.ProdDefensivo, Unit: "l"},
		{OrganizationID: organization.ID, Name: "Óleo Diesel", Type: entity.ProdCombustivel, Unit: "l"},
	}
	for i := range products {
		if err := db.Create(&products[i]).Error; err != nil {
			return fmt.Errorf("create product %s: %w", products[i].Name, err)
		}
	}
	fmt.Printf("  ✓ Produtos Agrícolas: %d criados\n", len(products))

	// Farms
	farms := []entity.Farm{
		{OrganizationID: organization.ID, Name: "Fazenda Recanto Verde", Owner: "João Silva", Location: "Alfenas - MG", TotalAreaHA: 120, PlantedAreaHA: 95},
		{OrganizationID: organization.ID, Name: "Sítio Boa Esperança", Owner: "João Silva", Location: "Machado - MG", TotalAreaHA: 45, PlantedAreaHA: 40},
		{OrganizationID: organization.ID, Name: "Fazenda Monte Alegre", Owner: "Maria Oliveira", Location: "Poços de Caldas - MG", TotalAreaHA: 200, PlantedAreaHA: 160},
	}
	for i := range farms {
		if err := db.Create(&farms[i]).Error; err != nil {
			return fmt.Errorf("create farm %s: %w", farms[i].Name, err)
		}
	}
	fmt.Printf("  ✓ Fazendas: %d criadas\n", len(farms))

	// Producers
	joaoUserID := users[1].ID
	mariaUserID := users[2].ID
	producers := []entity.Producer{
		{OrganizationID: organization.ID, FarmID: farms[0].ID, UserID: &joaoUserID, CPF: "123.456.789-00", Name: "João Silva", Phone: "(35) 99999-0001", Email: "joao@cafeos.com.br"},
		{OrganizationID: organization.ID, FarmID: farms[1].ID, UserID: &joaoUserID, CPF: "123.456.789-00", Name: "João Silva", Phone: "(35) 99999-0001", Email: "joao@cafeos.com.br"},
		{OrganizationID: organization.ID, FarmID: farms[2].ID, UserID: &mariaUserID, CPF: "987.654.321-00", Name: "Maria Oliveira", Phone: "(35) 99999-0002", Email: "maria@cafeos.com.br"},
	}
	for i := range producers {
		if err := db.Create(&producers[i]).Error; err != nil {
			return fmt.Errorf("create producer %s: %w", producers[i].Name, err)
		}
	}
	fmt.Printf("  ✓ Produtores: %d criados\n", len(producers))

	// Plots
	plots := []entity.Plot{
		{OrganizationID: organization.ID, FarmID: farms[0].ID, Name: "Talhão A-1", AreaHA: 30, Cultivar: "Catuaí Vermelho", PlantingYear: 2018, Altitude: 950, SoilType: "argiloso"},
		{OrganizationID: organization.ID, FarmID: farms[0].ID, Name: "Talhão A-2", AreaHA: 35, Cultivar: "Mundo Novo", PlantingYear: 2019, Altitude: 920, SoilType: "argiloso"},
		{OrganizationID: organization.ID, FarmID: farms[0].ID, Name: "Talhão B-1", AreaHA: 30, Cultivar: "Catuaí Amarelo", PlantingYear: 2020, Altitude: 980, SoilType: "arenoso"},
		{OrganizationID: organization.ID, FarmID: farms[1].ID, Name: "Talhão Único", AreaHA: 40, Cultivar: "Bourbon Amarelo", PlantingYear: 2017, Altitude: 1050, SoilType: "organico"},
		{OrganizationID: organization.ID, FarmID: farms[2].ID, Name: "Talhão Sul", AreaHA: 80, Cultivar: "Catuaí Vermelho", PlantingYear: 2016, Altitude: 1100, SoilType: "argiloso"},
		{OrganizationID: organization.ID, FarmID: farms[2].ID, Name: "Talhão Norte", AreaHA: 80, Cultivar: "Acauã", PlantingYear: 2018, Altitude: 1080, SoilType: "siltoso"},
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
		{OrganizationID: organization.ID, PlotID: plots[0].ID, Type: entity.OpAdubacao, Date: now.AddDate(0, -3, 0), Responsible: "Ana Costa", ProductUsed: "NPK 20-05-20", Quantity: 600, Cost: 4800, Notes: "Adubação de cobertura"},
		{OrganizationID: organization.ID, PlotID: plots[1].ID, Type: entity.OpAdubacao, Date: now.AddDate(0, -3, -2), Responsible: "Ana Costa", ProductUsed: "NPK 20-05-20", Quantity: 700, Cost: 5600, Notes: ""},
		{OrganizationID: organization.ID, PlotID: plots[0].ID, Type: entity.OpPulverizacao, Date: now.AddDate(0, -2, 0), Responsible: "Carlos Santos", ProductUsed: "Glyphosate", Quantity: 30, Cost: 1200, Notes: "Controle de plantas daninhas"},
		{OrganizationID: organization.ID, PlotID: plots[3].ID, Type: entity.OpPoda, Date: now.AddDate(0, -1, -15), Responsible: "Maria Oliveira", ProductUsed: "", Quantity: 0, Cost: 2500, Notes: "Poda de formação"},
		{OrganizationID: organization.ID, PlotID: plots[2].ID, Type: entity.OpIrrigacao, Date: now.AddDate(0, -1, -5), Responsible: "Ana Costa", ProductUsed: "", Quantity: 0, Cost: 1800, Notes: "Irrigação de salvamento"},
		{OrganizationID: organization.ID, PlotID: plots[4].ID, Type: entity.OpPulverizacao, Date: now.AddDate(0, 0, -10), Responsible: "Carlos Santos", ProductUsed: "Óleo Mineral", Quantity: 50, Cost: 2000, Notes: "Controle de ácaros"},
		{OrganizationID: organization.ID, PlotID: plots[0].ID, Type: entity.OpAdubacao, Date: now.AddDate(0, 0, -5), Responsible: "Ana Costa", ProductUsed: "Calcário Dolomítico", Quantity: 1500, Cost: 3000, Notes: "Calagem"},
	}
	for i := range operations {
		if err := db.Create(&operations[i]).Error; err != nil {
			return fmt.Errorf("create operation: %w", err)
		}
	}
	fmt.Printf("  ✓ Operações: %d criadas\n", len(operations))

	// Harvests
	harvests := []entity.Harvest{
		{OrganizationID: organization.ID, Year: 2024, Description: "Safra 2024", EstimatedProduction: 3000, Status: entity.HarvestFinalizada},
		{OrganizationID: organization.ID, Year: 2025, Description: "Safra 2025", EstimatedProduction: 3500, Status: entity.HarvestEmAndamento},
	}
	for i := range harvests {
		if err := db.Create(&harvests[i]).Error; err != nil {
			return fmt.Errorf("create harvest: %w", err)
		}
	}
	fmt.Printf("  ✓ Safras: %d criadas\n", len(harvests))

	// Harvest Productions
	productions := []entity.HarvestProduction{
		{OrganizationID: organization.ID, HarvestID: harvests[0].ID, PlotID: plots[0].ID, Quantity: 800, RecordedAt: now.AddDate(0, -6, 0), Notes: "Lote 1"},
		{OrganizationID: organization.ID, HarvestID: harvests[0].ID, PlotID: plots[1].ID, Quantity: 950, RecordedAt: now.AddDate(0, -6, -5), Notes: "Lote 2"},
		{OrganizationID: organization.ID, HarvestID: harvests[0].ID, PlotID: plots[2].ID, Quantity: 750, RecordedAt: now.AddDate(0, -5, -10), Notes: "Lote 3"},
		{OrganizationID: organization.ID, HarvestID: harvests[0].ID, PlotID: plots[3].ID, Quantity: 1100, RecordedAt: now.AddDate(0, -5, -15), Notes: ""},
		{OrganizationID: organization.ID, HarvestID: harvests[1].ID, PlotID: plots[0].ID, Quantity: 400, RecordedAt: now.AddDate(0, -1, 0), Notes: "Parcial 1"},
	}
	for i := range productions {
		if err := db.Create(&productions[i]).Error; err != nil {
			return fmt.Errorf("create production: %w", err)
		}
	}
	fmt.Printf("  ✓ Produções: %d criadas\n", len(productions))

	// Cost Centers (default seed)
	costCenters := []entity.CostCenter{
		{OrganizationID: organization.ID, Name: "Adubos", Code: "DESP_ADUBOS", Type: entity.CCDespesa, Description: "Fertilizantes e corretivos"},
		{OrganizationID: organization.ID, Name: "Defensivos", Code: "DESP_DEFENSIVOS", Type: entity.CCDespesa, Description: "Agrotóxicos e defensivos agrícolas"},
		{OrganizationID: organization.ID, Name: "Combustíveis", Code: "DESP_COMBUSTIVEL", Type: entity.CCDespesa, Description: "Diesel, gasolina e lubrificantes"},
		{OrganizationID: organization.ID, Name: "Mão de Obra", Code: "DESP_MAO_OBRA", Type: entity.CCDespesa, Description: "Salários e encargos trabalhistas"},
		{OrganizationID: organization.ID, Name: "Frete", Code: "DESP_FRETE", Type: entity.CCDespesa, Description: "Transporte de insumos e produção"},
		{OrganizationID: organization.ID, Name: "Manutenção", Code: "DESP_MANUTENCAO", Type: entity.CCDespesa, Description: "Manutenção de máquinas e equipamentos"},
		{OrganizationID: organization.ID, Name: "Irrigação", Code: "DESP_IRRIGACAO", Type: entity.CCDespesa, Description: "Custo de irrigação e água"},
		{OrganizationID: organization.ID, Name: "Análise de Solo", Code: "DESP_ANALISE_SOLO", Type: entity.CCDespesa, Description: "Análises laboratoriais de solo e folha"},
		{OrganizationID: organization.ID, Name: "Outros Insumos", Code: "DESP_OUTROS_INSUMOS", Type: entity.CCDespesa, Description: "Outros insumos agrícolas"},
		{OrganizationID: organization.ID, Name: "Serviços Terceiros", Code: "DESP_SERV_TERCEIROS", Type: entity.CCDespesa, Description: "Serviços contratados de terceiros"},
		{OrganizationID: organization.ID, Name: "Energia", Code: "DESP_ENERGIA", Type: entity.CCDespesa, Description: "Energia elétrica"},
		{OrganizationID: organization.ID, Name: "Depreciação", Code: "DESP_DEPRECIACAO", Type: entity.CCDespesa, Description: "Depreciação de máquinas e benfeitorias"},
		{OrganizationID: organization.ID, Name: "Administrativo", Code: "DESP_ADMINISTRATIVO", Type: entity.CCDespesa, Description: "Despesas administrativas"},
		{OrganizationID: organization.ID, Name: "Outras Despesas", Code: "DESP_OUTRAS", Type: entity.CCDespesa, Description: "Outras despesas operacionais"},
		{OrganizationID: organization.ID, Name: "Venda de Café", Code: "REC_CAFE", Type: entity.CCReceita, Description: "Receita com venda de café"},
		{OrganizationID: organization.ID, Name: "Venda de Mudas", Code: "REC_MUDAS", Type: entity.CCReceita, Description: "Receita com venda de mudas"},
		{OrganizationID: organization.ID, Name: "Outras Receitas", Code: "REC_OUTRAS", Type: entity.CCReceita, Description: "Outras receitas"},
	}
	for _, cc := range costCenters {
		if err := db.Create(&cc).Error; err != nil {
			return fmt.Errorf("create cost center %s: %w", cc.Name, err)
		}
	}
	fmt.Printf("  ✓ Centros de Custo: %d criados\n", len(costCenters))

	// Indicators for finalized harvest
	indicators := []entity.Indicator{
		{OrganizationID: organization.ID, HarvestID: harvests[0].ID, Type: entity.IndProducaoTotal, Value: 3600, CalculatedAt: now.AddDate(0, -4, 0)},
		{OrganizationID: organization.ID, HarvestID: harvests[0].ID, Type: entity.IndCustoTotal, Value: 45000, CalculatedAt: now.AddDate(0, -4, 0)},
		{OrganizationID: organization.ID, HarvestID: harvests[0].ID, Type: entity.IndSacasHA, Value: 37.89, CalculatedAt: now.AddDate(0, -4, 0)},
		{OrganizationID: organization.ID, HarvestID: harvests[0].ID, Type: entity.IndCustoSaca, Value: 12.50, CalculatedAt: now.AddDate(0, -4, 0)},
		{OrganizationID: organization.ID, HarvestID: harvests[0].ID, Type: entity.IndRentabilidade, Value: 85000, CalculatedAt: now.AddDate(0, -4, 0)},
	}
	for i := range indicators {
		if err := db.Create(&indicators[i]).Error; err != nil {
			return fmt.Errorf("create indicator: %w", err)
		}
	}
	fmt.Printf("  ✓ Indicadores: %d criados\n", len(indicators))

	// Financial (one per farm, plus one org-wide with no farm link)
	financials := []entity.FinancialTransaction{
		{OrganizationID: organization.ID, FarmID: &farms[0].ID, Type: entity.TranDespesa, Description: "Adubação Recanto Verde", Amount: 4800, Date: now.AddDate(0, -3, 0), DueDate: now.AddDate(0, -3, 5), Status: "paid"},
		{OrganizationID: organization.ID, FarmID: &farms[2].ID, Type: entity.TranReceita, Description: "Venda de café - Monte Alegre", Amount: 52000, Date: now.AddDate(0, -1, 0), DueDate: now.AddDate(0, -1, 0), Status: "paid"},
		{OrganizationID: organization.ID, Type: entity.TranDespesa, Description: "Contabilidade (organização)", Amount: 1200, Date: now, DueDate: now.AddDate(0, 0, 10), Status: "pending"},
	}
	for i := range financials {
		if err := db.Create(&financials[i]).Error; err != nil {
			return fmt.Errorf("create financial transaction: %w", err)
		}
	}
	fmt.Printf("  ✓ Transações Financeiras: %d criadas\n", len(financials))

	// Stock (one item per farm, plus one org-wide with no farm link)
	stockItems := []entity.StockItem{
		{OrganizationID: organization.ID, FarmID: &farms[0].ID, ProductID: products[0].ID, Quantity: 500, Unit: "kg", Location: "Galpão Recanto Verde"},
		{OrganizationID: organization.ID, FarmID: &farms[2].ID, ProductID: products[1].ID, Quantity: 300, Unit: "kg", Location: "Galpão Monte Alegre"},
		{OrganizationID: organization.ID, ProductID: products[4].ID, Quantity: 1000, Unit: "l", Location: "Depósito Central"},
	}
	for i := range stockItems {
		if err := db.Create(&stockItems[i]).Error; err != nil {
			return fmt.Errorf("create stock item: %w", err)
		}
	}
	fmt.Printf("  ✓ Itens de Estoque: %d criados\n", len(stockItems))

	// Fleet (one vehicle per farm, plus one shared/org-wide with no farm link)
	vehicles := []entity.Vehicle{
		{OrganizationID: organization.ID, FarmID: &farms[0].ID, Name: "Trator Recanto Verde", Type: entity.VeicTractor, Plate: "ABC1D23", Brand: "Massey Ferguson", Model: "4275", Year: 2019, Status: "active"},
		{OrganizationID: organization.ID, FarmID: &farms[2].ID, Name: "Trator Monte Alegre", Type: entity.VeicTractor, Plate: "XYZ9E87", Brand: "New Holland", Model: "TL75E", Year: 2021, Status: "active"},
		{OrganizationID: organization.ID, Name: "Caminhão da Cooperativa", Type: entity.VeicCaminhao, Plate: "QAZ2W34", Brand: "Volkswagen", Model: "Delivery", Year: 2018, Status: "active"},
	}
	for i := range vehicles {
		if err := db.Create(&vehicles[i]).Error; err != nil {
			return fmt.Errorf("create vehicle %s: %w", vehicles[i].Name, err)
		}
	}
	fmt.Printf("  ✓ Veículos: %d criados\n", len(vehicles))

	return nil
}
