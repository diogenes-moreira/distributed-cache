package distributed_cache

import (
	"os"
	"testing"
	"time"
)

var cache *Cache
var lruCache *LRUCache
var lruCacheWithTTL *LRUCacheWithTTL

func TestMain(m *testing.M) {
	// Setup code
	cache = NewCache("testCache", "255.255.255.255", ":12345")
	lruCache = NewLRUCache("testCache", "255.255.255.255", ":12346", 2)
	lruCacheWithTTL = NewLRUCacheWithTTL("testCacheWithTTL", "255.255.255.255", ":12347", 2, 2*time.Second)

	// Run tests
	code := m.Run()

	// Teardown code
	cache.StopListener()
	lruCache.StopListener()
	lruCacheWithTTL.StopListener()

	// Exit with the code from m.Run()
	os.Exit(code)
}
