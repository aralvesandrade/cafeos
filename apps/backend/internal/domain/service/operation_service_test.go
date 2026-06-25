package service

import (
	"testing"
	"time"

	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	mock "github.com/aralvesandrade/cafeos/internal/domain/service/testing"
	"github.com/aralvesandrade/cafeos/internal/event"
)

func TestOperationService_Create(t *testing.T) {
	repo := mock.NewInMemoryOperationRepo()
	bus := event.NewInMemoryBus()
	svc := NewOperationService(repo, bus)

	op, err := svc.Create("t1", "plot-1", entity.OpAdubacao, time.Now(), "João", "NPK 20-05-20", 500, 1500.00, "")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if op.Type != entity.OpAdubacao {
		t.Errorf("expected adubacao, got %s", op.Type)
	}
	if op.Cost != 1500.00 {
		t.Errorf("expected 1500.00, got %f", op.Cost)
	}
}

func TestOperationService_Validation(t *testing.T) {
	repo := mock.NewInMemoryOperationRepo()
	bus := event.NewInMemoryBus()
	svc := NewOperationService(repo, bus)

	tests := []struct {
		name  string
		plotID string
		opType entity.OperationType
		cost  float64
	}{
		{"empty plot", "", entity.OpAdubacao, 100},
		{"empty type", "p1", "", 100},
		{"negative cost", "p1", entity.OpColheita, -10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.Create("t1", tt.plotID, tt.opType, time.Now(), "", "", 0, tt.cost, "")
			if err == nil {
				t.Error("expected error, got nil")
			}
		})
	}
}

func TestOperationService_EventPublished(t *testing.T) {
	repo := mock.NewInMemoryOperationRepo()
	bus := event.NewInMemoryBus()
	svc := NewOperationService(repo, bus)

	received := false
	bus.Subscribe(func(e interface{}) {
		if _, ok := e.(event.OperationRegistered); ok {
			received = true
		}
	})

	svc.Create("t1", "p1", entity.OpPulverizacao, time.Now(), "", "", 0, 500, "")

	if !received {
		t.Error("expected OperationRegistered event to be published")
	}
}

func TestOperationService_ListRecent(t *testing.T) {
	repo := mock.NewInMemoryOperationRepo()
	bus := event.NewInMemoryBus()
	svc := NewOperationService(repo, bus)

	for i := 0; i < 15; i++ {
		svc.Create("t1", "p1", entity.OpIrrigacao, time.Now(), "", "", 0, 100, "")
	}

	recent, _ := svc.ListRecent("t1", 5)
	if len(recent) > 5 {
		t.Errorf("expected at most 5 recent operations, got %d", len(recent))
	}
}

func TestOperationService_Delete(t *testing.T) {
	repo := mock.NewInMemoryOperationRepo()
	bus := event.NewInMemoryBus()
	svc := NewOperationService(repo, bus)

	op, _ := svc.Create("t1", "p1", entity.OpPoda, time.Now(), "", "", 0, 100, "")
	err := svc.Delete(op.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
