package event

import "sync"

type Handler func(event interface{})

type Bus interface {
	Publish(event interface{})
	Subscribe(handler Handler) func()
}

type InMemoryBus struct {
	mu       sync.RWMutex
	handlers []Handler
}

func NewInMemoryBus() *InMemoryBus {
	return &InMemoryBus{}
}

func (b *InMemoryBus) Publish(event interface{}) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, h := range b.handlers {
		h(event)
	}
}

func (b *InMemoryBus) Subscribe(handler Handler) func() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.handlers = append(b.handlers, handler)

	idx := len(b.handlers) - 1
	return func() {
		b.mu.Lock()
		defer b.mu.Unlock()
		b.handlers = append(b.handlers[:idx], b.handlers[idx+1:]...)
	}
}
