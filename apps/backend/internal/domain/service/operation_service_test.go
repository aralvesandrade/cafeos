package service

import (
	"testing"
	"time"

	mock "github.com/aralvesandrade/cafeos/internal/domain/service/testing"
	"github.com/aralvesandrade/cafeos/internal/event"
)

func TestOperationService_Create(t *testing.T) {
	repo := mock.NewInMemoryOperationRepo()
	bus := event.NewInMemoryBus()
	svc := NewOperationService(repo, bus)

	op, err := svc.Create("t1", "plot-1", "type-adubacao", time.Now(), "João", "NPK 20-05-20", 500, 1500.00, "", nil, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if op.TypeID != "type-adubacao" {
		t.Errorf("expected type-adubacao, got %s", op.TypeID)
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
		name   string
		plotID string
		typeID string
		cost   float64
	}{
		{"empty plot", "", "type-adubacao", 100},
		{"empty type", "p1", "", 100},
		{"negative cost", "p1", "type-colheita", -10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.Create("t1", tt.plotID, tt.typeID, time.Now(), "", "", 0, tt.cost, "", nil, nil)
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

	svc.Create("t1", "p1", "type-pulverizacao", time.Now(), "", "", 0, 500, "", nil, nil)

	if !received {
		t.Error("expected OperationRegistered event to be published")
	}
}

func TestOperationService_ListRecent(t *testing.T) {
	repo := mock.NewInMemoryOperationRepo()
	bus := event.NewInMemoryBus()
	svc := NewOperationService(repo, bus)

	for i := 0; i < 15; i++ {
		svc.Create("t1", "p1", "type-irrigacao", time.Now(), "", "", 0, 100, "", nil, nil)
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

	op, _ := svc.Create("t1", "p1", "type-poda", time.Now(), "", "", 0, 100, "", nil, nil)
	err := svc.Delete(op.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
