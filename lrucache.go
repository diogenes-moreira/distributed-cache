package distributed_cache

import (
	"context"
	"github.com/google/uuid"
	"slices"
)

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

func (c *LRUCache) clean() {
	c.mutex.Lock()
	c.queue = make([]string, 0)
	c.storage = make(map[string]interface{})
	c.mutex.Unlock()
}

func (c *LRUCache) set(key string, value interface{}) {
	c.mutex.Lock()
	_, exists := c.storage[key]
	if !exists {
		c.queue = append(c.queue, key)
		if len(c.queue) > c.MaxEntries {
			delete(c.storage, c.queue[0])
			c.queue = c.queue[1:]
		}
	} else {
		index := slices.Index(c.queue, key)
		c.queue = append(c.queue[:index], c.queue[index+1:]...)
		c.queue = append(c.queue, key)
	}
	c.storage[key] = value
	c.mutex.Unlock()
}

func (c *LRUCache) delete(key string) {
	c.mutex.Lock()
	index := slices.Index(c.queue, key)
	if index != -1 {
		c.queue = append(c.queue[:index], c.queue[index+1:]...)
		delete(c.storage, key)
	}
	c.mutex.Unlock()
}

func (c *LRUCache) Set(key string, value interface{}) {
	if value == nil {
		c.Delete(key)
		return
	}
	c.sendSet(key, value)
	c.set(key, value)
}

// Delete deletes a value from the cache
// and sends the delete message to the other nodes
func (c *LRUCache) Delete(key string) {
	c.sendDelete(key)
	c.delete(key)
}

// Clean deletes all values from the cache
// and sends the sendClean message to the other nodes
func (c *LRUCache) Clean() {
	c.clean()
	c.sendClean()
}

func NewLRUCache(name string, address string, maxEntries int) *LRUCache {
	ctx, cancel := context.WithCancel(context.Background())
	c := &LRUCache{
		Cache: Cache{
			Name:         name,
			Address:      address,
			storage:      make(map[string]interface{}),
			context:      ctx,
			StopListener: cancel,
			node:         uuid.New(),
		},
		MaxEntries: maxEntries,
		queue:      make([]string, 0),
	}
	go startListener(c, ctx)
	return c
}
