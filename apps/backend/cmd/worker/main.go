package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	"github.com/aralvesandrade/cafeos/internal/event"
	"github.com/aralvesandrade/cafeos/internal/infra/config"
	"github.com/aralvesandrade/cafeos/internal/infra/db/postgres"
	infraRepo "github.com/aralvesandrade/cafeos/internal/infra/db/repository"
	infraLogger "github.com/aralvesandrade/cafeos/internal/infra/logger"
	"github.com/aralvesandrade/cafeos/internal/infra/messaging"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	cfg := config.Load()

	log := infraLogger.New(infraLogger.Config{
		Level:  cfg.LogLevel,
		Format: cfg.LogFormat,
	})
	slog.SetDefault(log)

	var err error
	db, err = postgres.NewConnection(cfg.DatabaseURL, log, slog.LevelInfo)
	if err != nil {
		log.Error("failed to connect database", "error", err)
		os.Exit(1)
	}

	conn, err := messaging.NewConnection(cfg.RabbitMQURL)
	if err != nil {
		log.Error("failed to connect rabbitmq", "error", err)
		os.Exit(1)
	}
	defer conn.Close()

	ch := conn.Channel()

	queues := []string{"sync.operations", "sync.stock", "sync.harvest", "sync.financial", "sync.labor"}
	pub := messaging.NewPublisher(ch)
	for _, q := range queues {
		if err := pub.DeclareQueue(q); err != nil {
			log.Error("failed to declare queue", "queue", q, "error", err)
			os.Exit(1)
		}
	}

	eventBus := event.NewInMemoryBus()
	event.SetupHandlers(eventBus)

	for _, q := range queues {
		consumer := messaging.NewConsumer(ch, q, processMessage)
		if err := consumer.Start(); err != nil {
			log.Error("failed to start consumer", "queue", q, "error", err)
			os.Exit(1)
		}
		log.Info("listening on queue", "queue", q)
	}

	log.Info("worker started, waiting for messages")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	log.Info("worker shutting down")
}

func processMessage(msg messaging.SyncMessage) error {
	slog.Info("processing message", "event", msg.EventType, "client", msg.ClientID, "organization", msg.OrganizationID)

	payload, err := json.Marshal(msg.Payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	switch msg.EventType {
	case "operation.created":
		return processOperationCreated(msg.OrganizationID, payload)
	case "stock.moved":
		return processStockMovement(msg.OrganizationID, payload)
	case "financial.created":
		return processFinancialTransaction(msg.OrganizationID, payload)
	case "harvest.production":
		return processHarvestProduction(msg.OrganizationID, payload)
	case "labor.shift":
		return processLaborShift(msg.OrganizationID, payload)
	default:
		slog.Warn("unknown event type", "event", msg.EventType)
		return nil
	}
}

type operationPayload struct {
	PlotID      string  `json:"plot_id"`
	Type        string  `json:"type"`
	Date        string  `json:"date"`
	Responsible string  `json:"responsible"`
	ProductUsed string  `json:"product_used"`
	Quantity    float64 `json:"quantity"`
	Cost        float64 `json:"cost"`
	Notes       string  `json:"notes"`
}

func processOperationCreated(organizationID string, payload []byte) error {
	var p operationPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		return err
	}
	otRepo := infraRepo.NewOperationTypeRepository(db)
	ot, err := otRepo.GetByOrganizationAndCode(organizationID, p.Type)
	if err != nil {
		return fmt.Errorf("resolve operation type %q: %w", p.Type, err)
	}
	repo := infraRepo.NewOperationRepository(db)
	op := &entity.Operation{
		ID:             uuid.New().String(),
		OrganizationID: organizationID,
		PlotID:         p.PlotID,
		TypeID:         ot.ID,
		Date:           parseTime(p.Date),
		Responsible:    p.Responsible,
		ProductUsed:    p.ProductUsed,
		Quantity:       p.Quantity,
		Cost:           p.Cost,
		Notes:          p.Notes,
	}
	return repo.Create(op)
}

type stockPayload struct {
	ProductID string  `json:"product_id"`
	Quantity  float64 `json:"quantity"`
	Type      string  `json:"type"`
	Reference string  `json:"reference"`
	Notes     string  `json:"notes"`
}

func processStockMovement(organizationID string, payload []byte) error {
	var p stockPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		return err
	}
	repo := infraRepo.NewStockMovementRepository(db)
	mov := &entity.StockMovement{
		ID:             uuid.New().String(),
		OrganizationID: organizationID,
		ItemID:         p.ProductID,
		Type:           p.Type,
		Quantity:       p.Quantity,
		Reference:      p.Reference,
		Notes:          p.Notes,
	}
	return repo.Create(mov)
}

type financialPayload struct {
	Type         string  `json:"type"`
	CostCenterID *string `json:"cost_center_id"`
	Description  string  `json:"description"`
	Amount       float64 `json:"amount"`
	Date         string  `json:"date"`
	Notes        string  `json:"notes"`
}

func processFinancialTransaction(organizationID string, payload []byte) error {
	var p financialPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		return err
	}
	repo := infraRepo.NewFinancialRepository(db)
	tx := &entity.FinancialTransaction{
		ID:             uuid.New().String(),
		OrganizationID: organizationID,
		Type:           entity.TransactionType(p.Type),
		CostCenterID:   p.CostCenterID,
		Description:    p.Description,
		Amount:         p.Amount,
		Date:           parseTime(p.Date),
		Status:         "pending",
	}
	return repo.Create(tx)
}

type harvestPayload struct {
	HarvestID string  `json:"harvest_id"`
	PlotID    string  `json:"plot_id"`
	Quantity  float64 `json:"quantity"`
	Notes     string  `json:"notes"`
}

func processHarvestProduction(organizationID string, payload []byte) error {
	var p harvestPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		return err
	}
	repo := infraRepo.NewHarvestProductionRepository(db)
	prod := &entity.HarvestProduction{
		ID:             uuid.New().String(),
		OrganizationID: organizationID,
		HarvestID:      p.HarvestID,
		PlotID:         p.PlotID,
		Quantity:       p.Quantity,
		Notes:          p.Notes,
	}
	return repo.Create(prod)
}

type laborPayload struct {
	WorkerID    string  `json:"worker_id"`
	OperationID string  `json:"operation_id"`
	Hours       float64 `json:"hours"`
	Cost        float64 `json:"cost"`
	Date        string  `json:"date"`
	Notes       string  `json:"notes"`
}

func processLaborShift(organizationID string, payload []byte) error {
	var p laborPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		return err
	}
	repo := infraRepo.NewWorkShiftRepository(db)
	var opRef *string
	if p.OperationID != "" {
		opRef = &p.OperationID
	}
	shift := &entity.WorkShift{
		ID:             uuid.New().String(),
		OrganizationID: organizationID,
		WorkerID:       p.WorkerID,
		OperationID:    opRef,
		Hours:          p.Hours,
		Cost:           p.Cost,
		Date:           parseTime(p.Date),
		Notes:          p.Notes,
	}
	return repo.Create(shift)
}

func parseTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		t, err = time.Parse("2006-01-02", s)
		if err != nil {
			return time.Now()
		}
	}
	return t
}
