package main

import (
	distributed_cache "github.com/diogenes-moreira/distributed-cache"
	"time"
)

func main() {
	// Create a new cache with the name "cache" and the address UDP Port ":12345".
	// Name is used to identify the cache and address is used to send messages to the cache.
	// start the listener in a goroutine. is used to listen for incoming messages.
	cache := distributed_cache.NewCache("cache", ":12345")
	cache.Set("key", "value")
	value := cache.Get("key")
	if value != nil {
		println(value.(string))
	}

	// LRU Cache Extend the Cache struct and limit the number of entries to 10.
	lruCache := distributed_cache.NewLRUCache("lru", ":12345", 10)
	lruCache.Set("key", "value")
	value = lruCache.Get("key")

	// LRU Cache with TTL Extend the lruCache struct and add a TTL of 10 seconds.
	lruCacheWithTTL := distributed_cache.NewLRUCacheWithTTL("lru", ":12345", 10, time.Second*10)
	lruCacheWithTTL.Set("key", "value")
	value = lruCache.Get("key")
}
