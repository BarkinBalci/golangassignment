package memory

import (
	"testing"
)

func TestCreate(t *testing.T) {
	item := NewItem()
	key := "testKey"
	value := "testValue"

	item.Create(key, value)
	fetchedValue, ok := item.Fetch(key)

	if !ok {
		t.Fatalf("Expected to find the key %v in the item after create", key)
	}
	if fetchedValue != value {
		t.Fatalf("Expected value %v, but found %v", value, fetchedValue)
	}
}

func TestFetch(t *testing.T) {
	item := NewItem()
	key := "testKey"
	value := "testValue"
	item.Create(key, value)

	fetchedValue, ok := item.Fetch(key)
	if !ok {
		t.Fatalf("Expected to find the key %v in the item", key)
	}
	if fetchedValue != value {
		t.Fatalf("Expected value %v, but found %v", value, fetchedValue)
	}
}

func TestFetchKeyNotFound(t *testing.T) {
	item := NewItem()
	key := "nonExistentKey"

	_, ok := item.Fetch(key)
	if ok {
		t.Fatalf("Expected key %s to not be found, but it was", key)
	}
}
