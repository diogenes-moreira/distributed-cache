package distributed_cache

import (
	"testing"
)

func TestLRUCache_Set(t *testing.T) {
	lruCache.Set("key1", "value1")

	if val := lruCache.Get("key1"); val != "value1" {
		t.Errorf("Expected value1, got %v", val)
	}

	lruCache.Set("key2", nil)
	if val := lruCache.Get("key2"); val != nil {
		t.Errorf("Expected nil, got %v", val)
	}
}

func TestLRUCache_Get(t *testing.T) {
	lruCache.Set("key1", "value1")

	if val := lruCache.Get("key1"); val != "value1" {
		t.Errorf("Expected value1, got %v", val)
	}

	if val := lruCache.Get("key2"); val != nil {
		t.Errorf("Expected nil, got %v", val)
	}
}

func TestLRUCache_Delete(t *testing.T) {
	lruCache.Set("key1", "value1")
	lruCache.Delete("key1")

	if val := lruCache.Get("key1"); val != nil {
		t.Errorf("Expected nil, got %v", val)
	}
}

func TestLRUCache_Clean(t *testing.T) {
	lruCache.Set("key1", "value1")
	lruCache.Set("key2", "value2")
	lruCache.Clean()

	if val := lruCache.Get("key1"); val != nil {
		t.Errorf("Expected nil, got %v", val)
	}

	if val := lruCache.Get("key2"); val != nil {
		t.Errorf("Expected nil, got %v", val)
	}
}

func TestLRUCache_Eviction(t *testing.T) {
	lruCache.Set("key1", "value1")
	lruCache.Set("key2", "value2")
	lruCache.Set("key3", "value3") // This should evict "key1"

	if val := lruCache.Get("key1"); val != nil {
		t.Errorf("Expected nil, got %v", val)
	}

	if val := lruCache.Get("key2"); val != "value2" {
		t.Errorf("Expected value2, got %v", val)
	}

	if val := lruCache.Get("key3"); val != "value3" {
		t.Errorf("Expected value3, got %v", val)
	}
}
