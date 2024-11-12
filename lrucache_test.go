package distributed_cache

import (
	"testing"
)

func TestLRUCache_Set(t *testing.T) {
	lruCache.Set("key1", "value1")
	lruCache.Set("key2", "value2")

	if len(lruCache.storage) != 2 {
		t.Errorf("Expected cache size 2, got %d", len(lruCache.storage))
	}

	lruCache.Set("key3", "value3")
	if len(lruCache.storage) != 2 {
		t.Errorf("Expected cache size 2 after eviction, got %d", len(lruCache.storage))
	}

	if _, exists := lruCache.storage["key1"]; exists {
		t.Errorf("Expected key1 to be evicted")
	}
}

func TestLRUCache_Get(t *testing.T) {
	lruCache.Set("key1", "value1")

	value := lruCache.Get("key1")
	if value != "value1" {
		t.Errorf("Expected value1, got %v", value)
	}

	value = lruCache.Get("key2")
	if value != nil {
		t.Errorf("Expected nil, got %v", value)
	}
}

func TestLRUCache_Delete(t *testing.T) {
	lruCache.Set("key1", "value1")
	lruCache.Delete("key1")

	if _, exists := lruCache.storage["key1"]; exists {
		t.Errorf("Expected key1 to be deleted")
	}
}

func TestLRUCache_pushFront(t *testing.T) {
	lruCache.Clean()
	lruCache.pushFront("key1")

	if len(lruCache.queue) != 1 || lruCache.queue[0] != "key1" {
		t.Errorf("Expected key1 at the front of the queue, got %v", lruCache.queue)
	}
}

func TestLRUCache_deleteFromQueue(t *testing.T) {
	lruCache.Clean()
	lruCache.pushFront("key1")
	lruCache.pushFront("key2")
	lruCache.deleteFromQueue("key1")

	if len(lruCache.queue) != 1 || lruCache.queue[0] != "key2" {
		t.Errorf("Expected key1 to be deleted from the queue, got %v", lruCache.queue)
	}
}

func TestLRUCache_deleteLast(t *testing.T) {
	lruCache.Clean()
	lruCache.pushFront("key1")
	lruCache.pushFront("key2")
	lruCache.deleteLast()

	if len(lruCache.queue) != 1 || lruCache.queue[0] != "key2" {
		t.Errorf("Expected key1 to be deleted from the queue, got %v", lruCache.queue)
	}
}
