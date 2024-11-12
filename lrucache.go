package distributed_cache

import "sync"

//In this file, you can find the LRUCache struct.
//That implements the Cache interface.
//All methods are implemented in this file to create a cache that uses the Least
//Recently Used (LRU) algorithm to evict entries when the cache is full.
//Only NewLRUCache is exported, the rest of the methods are private

// LRUCache is a cache
// that uses the Least Recently Used algorithm to evict entries
type LRUCache struct {
	Cache
	MaxEntries int
	queue      []string
	queueMutex sync.Mutex
}

// pushFront adds a key to the front of the queue
func (c *LRUCache) pushFront(key string) {
	c.queueMutex.Lock()
	c.queue = append([]string{key}, c.queue...)
	c.queueMutex.Unlock()
}

// deleteFromQueue deletes a key from the queue
func (c *LRUCache) deleteFromQueue(key string) {
	c.queueMutex.Lock()
	defer c.queueMutex.Unlock()
	for i, k := range c.queue {
		if k == key {
			c.queue = append(c.queue[:i], c.queue[i+1:]...)
			break
		}
	}
}

// deleteLast deletes the last key from the queue
func (c *LRUCache) deleteLast() {
	if len(c.queue) == 0 {
		return
	}
	c.queueMutex.Lock()
	key := c.queue[len(c.queue)-1]
	c.queue = c.queue[:len(c.queue)-1]
	c.queueMutex.Unlock()
	c.Cache.Delete(key)
}

// Set sets a value in the cache and sends it to the other nodes
func (c *LRUCache) Set(key string, value interface{}) {
	if value == nil {
		c.Delete(key)
		return
	}
	c.sendSet(key, value)
	if len(c.storage) >= c.MaxEntries {
		c.evict()
	}
	c.storage[key] = value
	c.pushFront(key)
}

// evict deletes the last key from the queue
func (c *LRUCache) evict() {
	c.deleteLast()
}

// Get gets a value from the cache
func (c *LRUCache) Get(key string) interface{} {
	value := c.Cache.Get(key)
	if value != nil {
		c.pushFront(key)
	}
	return value
}

// Delete deletes a value from the cache
func (c *LRUCache) Delete(key string) {
	c.deleteFromQueue(key)
	c.Cache.Delete(key)
}

func (c *LRUCache) Clean() {
	c.Cache.Clean()
	c.queueMutex.Lock()
	c.queue = make([]string, 0)
	c.queueMutex.Unlock()
}

// NewLRUCache creates a new LRUCache with the given name,
// address, and maxEntries.
// It also starts a listener to receive messages from other nodes
// name is the name of the cache
// address is the address of the cache
// maxEntries is the maximum number of entries that the cache can have
func NewLRUCache(name, address string, maxEntries int) *LRUCache {
	c := LRUCache{
		Cache:      *NewCache(name, address),
		MaxEntries: maxEntries,
		queue:      make([]string, 0),
		queueMutex: sync.Mutex{},
	}
	return &c
}
