package memory

import (
	"sync"
)

type Item struct {
	data map[string]string
	mu   sync.RWMutex
}

// NewItem Initializes a new Item.
func NewItem() *Item {
	return &Item{
		data: make(map[string]string),
	}
}

// Create adds or updates a key-value pair in the Item.
func (item *Item) Create(key string, value string) {
	item.mu.Lock()
	defer item.mu.Unlock()

	item.data[key] = value
}

// Fetch retrieves the value associated with a given key.
func (item *Item) Fetch(key string) (string, bool) {
	item.mu.RLock()
	defer item.mu.RUnlock()
	value, ok := item.data[key]
	return value, ok
}
