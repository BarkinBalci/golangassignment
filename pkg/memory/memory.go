package memory

import (
	"sync"
)

type Item struct {
	data map[string]string
	mu   sync.RWMutex
}

func NewItem() *Item {
	return &Item{
		data: make(map[string]string),
	}
}

func (item *Item) Create(key string, value string) {
	item.mu.Lock()
	defer item.mu.Unlock()

	item.data[key] = value
}

func (item *Item) Fetch(key string) (string, bool) {
	item.mu.RLock()
	defer item.mu.RUnlock()
	value, ok := item.data[key]
	return value, ok
}
