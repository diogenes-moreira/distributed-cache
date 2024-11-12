package distributed_cache

import (
	"testing"
	"time"
)

func TestLRUCacheWithTTL_Set(t *testing.T) {
	lruCacheWithTTL.Set("key1", "value1")

	if val := lruCacheWithTTL.Get("key1"); val != "value1" {
		t.Errorf("Expected value1, got %v", val)
	}

	lruCacheWithTTL.Set("key2", nil)
	if val := lruCacheWithTTL.Get("key2"); val != nil {
		t.Errorf("Expected nil, got %v", val)
	}
}

func TestLRUCacheWithTTL_Get(t *testing.T) {
	lruCacheWithTTL.Set("key1", "value1")

	if val := lruCacheWithTTL.Get("key1"); val != "value1" {
		t.Errorf("Expected value1, got %v", val)
	}

	if val := lruCacheWithTTL.Get("key2"); val != nil {
		t.Errorf("Expected nil, got %v", val)
	}
}

func TestLRUCacheWithTTL_Delete(t *testing.T) {
	lruCacheWithTTL.Set("key1", "value1")
	lruCacheWithTTL.Delete("key1")

	if val := lruCacheWithTTL.Get("key1"); val != nil {
		t.Errorf("Expected nil, got %v", val)
	}
}

func TestLRUCacheWithTTL_Clean(t *testing.T) {
	lruCacheWithTTL.Set("key1", "value1")
	lruCacheWithTTL.Set("key2", "value2")
	lruCacheWithTTL.Clean()

	if val := lruCacheWithTTL.Get("key1"); val != nil {
		t.Errorf("Expected nil, got %v", val)
	}

	if val := lruCacheWithTTL.Get("key2"); val != nil {
		t.Errorf("Expected nil, got %v", val)
	}
}

func TestLRUCacheWithTTL_TTLExpiration(t *testing.T) {
	lruCacheWithTTL.Clean()
	lruCacheWithTTL.Set("key1", "value1")
	time.Sleep(3 * time.Second)

	if val := lruCacheWithTTL.Get("key1"); val != nil {
		t.Errorf("Expected nil, got %v", val)
	}
}
