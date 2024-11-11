package distributed_cache

import "time"

type LRUCacheWithTTL struct {
	LRUCache
	TTL    time.Duration
	TTLMap map[string]time.Time
}

func (c *LRUCacheWithTTL) pushFront(key string) {
	c.LRUCache.pushFront(key)
	c.TTLMap[key] = time.Now().Add(c.TTL)
}

func (c *LRUCacheWithTTL) deleteLast() {
	if len(c.queue) == 0 {
		return
	}
	key := c.queue[len(c.queue)-1]
	delete(c.TTLMap, key)
	c.LRUCache.deleteLast()
}

func (c *LRUCacheWithTTL) deleteFromQueue(key string) {
	delete(c.TTLMap, key)
	c.LRUCache.deleteFromQueue(key)
}

func (c *LRUCacheWithTTL) evict() {
	for key, ttl := range c.TTLMap {
		if ttl.Before(time.Now()) {
			c.Delete(key)
		}
	}
}

func NewLRUCacheWithTTL(name, address string, maxEntries int, ttl time.Duration) *LRUCacheWithTTL {
	return &LRUCacheWithTTL{
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
}
