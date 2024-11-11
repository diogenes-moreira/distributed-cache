package distributed_cache

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
}

// pushFront adds a key to the front of the queue
func (c *LRUCache) pushFront(key string) {
	c.queue = append([]string{key}, c.queue...)
}

// deleteFromQueue deletes a key from the queue
func (c *LRUCache) deleteFromQueue(key string) {
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
	key := c.queue[len(c.queue)-1]
	c.queue = c.queue[:len(c.queue)-1]
	c.Delete(key)
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
	value := c.storage[key]
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

// NewLRUCache creates a new LRUCache with the given name,
// address, and maxEntries.
// It also starts a listener to receive messages from other nodes
// name is the name of the cache
// address is the address of the cache
// maxEntries is the maximum number of entries that the cache can have
func NewLRUCache(name, address string, maxEntries int) *LRUCache {
	c := &LRUCache{
		Cache: Cache{
			Name:    name,
			Address: address,
			storage: make(map[string]interface{}),
		},
		MaxEntries: maxEntries,
	}
	go c.startListener()
	return c
}
