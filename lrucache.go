package distributed_cache

type LRUCache struct {
	Cache
	MaxEntries int
	queue      []string
}

func (c *LRUCache) pushFront(key string) {
	c.queue = append([]string{key}, c.queue...)
}

func (c *LRUCache) deleteFromQueue(key string) {
	for i, k := range c.queue {
		if k == key {
			c.queue = append(c.queue[:i], c.queue[i+1:]...)
			break
		}
	}
}

func (c *LRUCache) deleteLast() {
	if len(c.queue) == 0 {
		return
	}
	key := c.queue[len(c.queue)-1]
	c.queue = c.queue[:len(c.queue)-1]
	c.Delete(key)
}

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

func (c *LRUCache) evict() {
	c.deleteLast()
}

func (c *LRUCache) Get(key string) interface{} {
	value := c.storage[key]
	if value != nil {
		c.pushFront(key)
	}
	return value
}

func (c *LRUCache) Delete(key string) {
	c.deleteFromQueue(key)
	c.Cache.Delete(key)
}

func NewLRUCache(name, address string, maxEntries int) *LRUCache {
	return &LRUCache{
		Cache: Cache{
			Name:    name,
			Address: address,
			storage: make(map[string]interface{}),
		},
		MaxEntries: maxEntries,
	}
}
