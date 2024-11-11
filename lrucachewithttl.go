package distributed_cache

// In this file, you can find the LRUCacheWithTTL struct that is used to create a cache
// All methods are implemented in this file to create a cache that uses the Least
//Recently Used (LRU) algorithm to evict entries when the cache is full.
//Only NewLRUCacheWithTTL is exported, the rest of the methods are private

import "time"

// LRUCacheWithTTL is a cache that uses the Least Recently Used (LRU) algorithm
// to evict entries when the cache is full.
// It also has a time-to-live (TTL) for each entry if an entry is not accessed,
// it will be deleted after the TTL expires.
// The Evict method is called periodically (TTL Period)
// to delete entries that have expired.
type LRUCacheWithTTL struct {
	LRUCache
	TTL    time.Duration
	TTLMap map[string]time.Time
}

// pushFront adds a key to the front of the queue and updates the TTLMap.
func (c *LRUCacheWithTTL) pushFront(key string) {
	c.LRUCache.pushFront(key)
	c.TTLMap[key] = time.Now().Add(c.TTL)
}

// deleteLast deletes the last key from the queue and updates the TTLMap.
func (c *LRUCacheWithTTL) deleteLast() {
	if len(c.queue) == 0 {
		return
	}
	key := c.queue[len(c.queue)-1]
	delete(c.TTLMap, key)
	c.LRUCache.deleteLast()
}

// deleteFromQueue deletes a key from the queue and updates the TTLMap.
func (c *LRUCacheWithTTL) deleteFromQueue(key string) {
	delete(c.TTLMap, key)
	c.LRUCache.deleteFromQueue(key)
}

// evict deletes the last key from the queue and updates the TTLMap.
func (c *LRUCacheWithTTL) evict() {
	for key, ttl := range c.TTLMap {
		if ttl.Before(time.Now()) {
			c.Delete(key)
		}
	}
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
		LRUCache: LRUCache{
			Cache: Cache{
				Name:    name,
				Address: address,
				storage: make(map[string]interface{}),
			},
			MaxEntries: maxEntries,
		},
		TTL:    ttl,
		TTLMap: make(map[string]time.Time),
	}

	go c.startListener()
	go func() {
		for {
			time.Sleep(ttl)
			c.evict()
		}
	}()
	return c
}
