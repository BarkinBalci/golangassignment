package memorydb

import (
	"sync"
)

type InMemoryDB struct {
	data map[string]interface{}
	mu   sync.RWMutex
}
