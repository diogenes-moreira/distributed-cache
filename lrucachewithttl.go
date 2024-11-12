package distributed_cache

// In this file, you can find the LRUCacheWithTTL struct that is used to create a cache
// All methods are implemented in this file to create a cache that uses the Least
//Recently Used (LRU) algorithm to evict entries when the cache is full.
//Only NewLRUCacheWithTTL is exported, the rest of the methods are private

import (
	"context"
	"sync"
	"time"
)

// LRUCacheWithTTL is a cache that uses the Least Recently Used (LRU) algorithm
// to evict entries when the cache is full.
// It also has a time-to-live (TTL) for each entry if an entry is not accessed,
// it will be deleted after the TTL expires.
// The Evict method is called periodically (TTL Period)
// to delete entries that have expired.
type LRUCacheWithTTL struct {
	LRUCache
	TTL      time.Duration
	ttlMap   map[string]time.Time
	ttlMutex sync.Mutex
}

// pushFront adds a key to the front of the queue and updates the ttlMap.
func (c *LRUCacheWithTTL) pushFront(key string) {
	c.LRUCache.pushFront(key)
	c.ttlMutex.Lock()
	c.ttlMap[key] = time.Now().Add(c.TTL)
	c.ttlMutex.Unlock()
}

// deleteLast deletes the last key from the queue and updates the ttlMap.
func (c *LRUCacheWithTTL) deleteLast() {
	c.ttlMutex.Lock()
	defer c.ttlMutex.Unlock()
	if len(c.queue) == 0 {
		return
	}
	key := c.queue[len(c.queue)-1]
	delete(c.ttlMap, key)
	c.LRUCache.deleteLast()
}

// deleteFromQueue deletes a key from the queue and updates the ttlMap.
func (c *LRUCacheWithTTL) deleteFromQueue(key string) {
	delete(c.ttlMap, key)
	c.LRUCache.deleteFromQueue(key)
}

func (c *LRUCacheWithTTL) Set(key string, value interface{}) {
	c.evict()
	c.pushFront(key)
	c.LRUCache.Set(key, value)
}

func (c *LRUCacheWithTTL) Get(key string) interface{} {
	c.evict()
	out := c.LRUCache.Get(key)
	if out != nil {
		c.pushFront(key)
	}
	return out
}

func (c *LRUCacheWithTTL) Delete(key string) {
	c.ttlMutex.Lock()
	c.delete(key)
	c.ttlMutex.Unlock()
}

func (c *LRUCacheWithTTL) delete(key string) {
	c.deleteFromQueue(key)
	c.LRUCache.Delete(key)
	delete(c.ttlMap, key)
}

// evict deletes the last key from the queue and updates the ttlMap.
func (c *LRUCacheWithTTL) evict() {
	c.ttlMutex.Lock()
	defer c.ttlMutex.Unlock()
	for key, ttl := range c.ttlMap {
		if ttl.Before(time.Now()) {
			c.delete(key)
		}
	}
}

func (c *LRUCacheWithTTL) Clean() {
	c.LRUCache.Clean()
	c.ttlMutex.Lock()
	c.ttlMap = make(map[string]time.Time)
	c.ttlMutex.Unlock()
}

// NewLRUCacheWithTTL creates a new LRUCacheWithTTL with the given name,
// address, maxEntries and TTL.
// It also starts a listener to receive messages from other nodes and starts a
// for cache eviction.
// name is the name of the cache
// address is the address of the cache
// maxEntries is the maximum number of entries that the cache can have
// ttl is the time-to-live for each entry in the cache
func NewLRUCacheWithTTL(name, address string, maxEntries int, ttl time.Duration) *LRUCacheWithTTL {
	c := &LRUCacheWithTTL{
		LRUCache: *NewLRUCache(name, address, maxEntries),
		TTL:      ttl,
		ttlMap:   make(map[string]time.Time),
		ttlMutex: sync.Mutex{},
	}

	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				time.Sleep(ttl)
				c.evict()
			}
		}
	}(c.context)
	return c
}
