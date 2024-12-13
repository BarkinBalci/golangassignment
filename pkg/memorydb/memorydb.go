package memorydb

import (
	"sync"
)

type InMemoryDB struct {
	data map[string]interface{}
	mu   sync.RWMutex
}

func NewInMemoryDB() *InMemoryDB {
	return &InMemoryDB{
		data: make(map[string]interface{}),
	}
}

func (db *InMemoryDB) Create(key string, value interface{}) {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.data[key] = value
}

func (db *InMemoryDB) FetchAll() map[string]interface{} {
	db.mu.RLock()
	defer db.mu.RUnlock()

	dataCopy := make(map[string]interface{})
	for key, value := range db.data {
		dataCopy[key] = value
	}
	return dataCopy
}
