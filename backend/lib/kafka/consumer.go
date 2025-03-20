package kafka

import (
	"context"
	"sync"
)

type ConsumerTracker struct {
	mu        sync.Mutex
	consumers map[string]context.CancelFunc
}

func NewConsumerTracker() *ConsumerTracker {
	return &ConsumerTracker{
		consumers: make(map[string]context.CancelFunc),
	}
}

func (t *ConsumerTracker) AddConsumer(id string, cancel context.CancelFunc) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.consumers[id] = cancel
}

func (t *ConsumerTracker) RemoveConsumer(id string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	if cancel, exists := t.consumers[id]; exists {
		cancel()
		delete(t.consumers, id)
		return true
	}
	return false
}

func (t *ConsumerTracker) StopAllConsumers() {
	t.mu.Lock()
	defer t.mu.Unlock()
	for id, cancel := range t.consumers {
		cancel()
		delete(t.consumers, id)
	}
}

func (t *ConsumerTracker) ConsumerExists(id string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	_, exists := t.consumers[id]
	return exists
}

