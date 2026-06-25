package event

import (
	"testing"
)

func TestInMemoryBus_Publish(t *testing.T) {
	bus := NewInMemoryBus()

	received := false
	unsub := bus.Subscribe(func(e interface{}) {
		received = true
	})
	defer unsub()

	bus.Publish(OperationRegistered{OperationID: "op-1"})

	if !received {
		t.Error("expected handler to be called")
	}
}

func TestInMemoryBus_MultipleHandlers(t *testing.T) {
	bus := NewInMemoryBus()

	count := 0
	h1 := func(e interface{}) { count++ }
	h2 := func(e interface{}) { count++ }

	unsub1 := bus.Subscribe(h1)
	defer unsub1()
	unsub2 := bus.Subscribe(h2)
	defer unsub2()

	bus.Publish(OperationRegistered{})

	if count != 2 {
		t.Errorf("expected count to be 2, got %d", count)
	}
}

func TestInMemoryBus_Unsubscribe(t *testing.T) {
	bus := NewInMemoryBus()

	count := 0
	h := func(e interface{}) { count++ }

	unsub := bus.Subscribe(h)
	bus.Publish(OperationRegistered{})

	unsub()
	bus.Publish(OperationRegistered{})

	if count != 1 {
		t.Errorf("expected count to be 1 (after unsubscribe), got %d", count)
	}
}

func TestInMemoryBus_NoHandlers(t *testing.T) {
	bus := NewInMemoryBus()

	bus.Publish(OperationRegistered{})
}

func TestInMemoryBus_EventTypes(t *testing.T) {
	bus := NewInMemoryBus()

	var receivedEvents []string
	bus.Subscribe(func(e interface{}) {
		switch e.(type) {
		case OperationRegistered:
			receivedEvents = append(receivedEvents, "op")
		case HarvestFinalized:
			receivedEvents = append(receivedEvents, "harvest")
		case IndicatorUpdated:
			receivedEvents = append(receivedEvents, "indicator")
		case AlertGenerated:
			receivedEvents = append(receivedEvents, "alert")
		}
	})

	bus.Publish(OperationRegistered{})
	bus.Publish(HarvestFinalized{})
	bus.Publish(IndicatorUpdated{})
	bus.Publish(AlertGenerated{})

	if len(receivedEvents) != 4 {
		t.Errorf("expected 4 events, got %d", len(receivedEvents))
	}
}
