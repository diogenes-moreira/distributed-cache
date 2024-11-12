# distributed-cache
GO Distributed Cache

## Description
This is a distributed cache system that  to distribute the data across multiple nodes. The system is designed to be fault-tolerant and can handle node failures. The system is implemented in Go and uses UPD Communication 

## Features

LRU Cache


LRU Cache With TTL

Example
```go
package main

import (
	distributed_cache "github.com/diogenes-moreira/distributed-cache"
	"time"
)

	


func main() {
	// Create a new cache with the name "cache" "255.255.255.255" is the mask to broadcast the cache sync messages
	// and the address UDP Port ":12345".
	// Name is used to identify the cache and address is used to send messages to the cache.
	// start the listener in a goroutine. is used to listen for incoming messages.
	cache := distributed_cache.NewCache("cache","255.255.255.255", ":12345")
	cache.Set("key", "value")
	value := cache.Get("key")
	if value != nil {
		println(value.(string))
	}
	
	//In addition hoy can set a function to be called when the key is not found in the cache
	cache.Filler = func(key string) (interface{}, error) {
        return "value",nil
    }
	
	// LRU Cache Extend the Cache struct and limit the number of entries to 10.
	lruCache := distributed_cache.NewLRUCache("lru","255.255.255.255", ":12345", 10)
	lruCache.Set("key", "value")
	value = lruCache.Get("key")
	
	// LRU Cache with TTL Extend the lruCache struct and add a TTL of 10 seconds.
	lruCacheWithTTL := distributed_cache.NewLRUCacheWithTTL("lru","255.255.255.255", ":12345", 10, time.Second*10)
	lruCache.Set("key", "value")
	value = lruCache.Get("key")
	
	//If you need to stop the cache listener you can call the Stop method
	cache.StopListener()
	
	//If you need add a Hook to be called when a key is removed from the cache
	//the hook run in a goroutine
	cache.RemoveHook=func(key string, value interface{}) {
        println("Key removed: ", key)
    }
	
}