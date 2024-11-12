package distributed_cache

import (
	"testing"
	"time"
)

func TestCache_Set(t *testing.T) {
	cache.Set("key1", "value1")

	if val := cache.Get("key1"); val != "value1" {
		t.Errorf("Expected value1, got %v", val)
	}

	cache.Set("key2", nil)
	if val := cache.Get("key2"); val != nil {
		t.Errorf("Expected nil, got %v", val)
	}
}

func TestCache_Get(t *testing.T) {
	cache.Set("key1", "value1")

	if val := cache.Get("key1"); val != "value1" {
		t.Errorf("Expected value1, got %v", val)
	}

	if val := cache.Get("key2"); val != nil {
		t.Errorf("Expected nil, got %v", val)
	}
}

func TestCache_Delete(t *testing.T) {
	cache.Set("key1", "value1")
	cache.Delete("key1")

	if val := cache.Get("key1"); val != nil {
		t.Errorf("Expected nil, got %v", val)
	}
}

func TestCache_Clean(t *testing.T) {
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Clean()

	if val := cache.Get("key1"); val != nil {
		t.Errorf("Expected nil, got %v", val)
	}

	if val := cache.Get("key2"); val != nil {
		t.Errorf("Expected nil, got %v", val)
	}
}

func TestCache_Remove(t *testing.T) {
	var removedKeys []string
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.RemoveHook = func(key string, value interface{}) {
		removedKeys = append(removedKeys, key)
	}
	cache.Clean()
	time.Sleep(5 * time.Second)
	if len(removedKeys) != 2 {
		t.Errorf("Expected 2 removed keys, got %v", len(removedKeys))
	}
	if removedKeys[0] != "key1" {
		t.Errorf("Expected key1, got %v", removedKeys[0])
	}
	if removedKeys[1] != "key2" {
		t.Errorf("Expected key2, got %v", removedKeys[1])
	}
}
