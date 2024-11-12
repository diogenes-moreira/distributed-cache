package distributed_cache

import (
	"testing"
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
