package distributed_cache

// In this file, you can find the LRUCacheWithTTL struct that is used to create a cache
// All methods are implemented in this file to create a cache that uses the Least
//Recently Used (LRU) algorithm to evict entries when the cache is full.
//Only NewLRUCacheWithTTL is exported, the rest of the methods are private

import (
	"context"
	"github.com/google/uuid"
	"log"
	"slices"
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
	TTL    time.Duration
	ttlMap map[string]time.Time
}

func (c *LRUCacheWithTTL) clean() {
	c.mutex.Lock()
	if c.RemoveHook != nil {
		for key, value := range c.storage {
			go c.RemoveHook(key, value)
		}
	}
	c.queue = make([]string, 0)
	c.storage = make(map[string]interface{})
	c.ttlMap = make(map[string]time.Time)
	c.mutex.Unlock()
}

func (c *LRUCacheWithTTL) set(key string, value interface{}) {
	c.evict()
	c.mutex.Lock()
	_, exists := c.storage[key]
	if !exists {
		c.queue = append(c.queue, key)
		if len(c.queue) > c.MaxEntries {
			if c.RemoveHook != nil {
				go c.RemoveHook(c.queue[0], c.storage[c.queue[0]])
			}
			delete(c.storage, c.queue[0])
			c.queue = c.queue[1:]
		}
	} else {
		index := slices.Index(c.queue, key)
		c.queue = append(c.queue[:index], c.queue[index+1:]...)
		c.queue = append(c.queue, key)
	}
	c.ttlMap[key] = time.Now().Add(c.TTL)
	c.storage[key] = value
	c.mutex.Unlock()
}

func (c *LRUCacheWithTTL) delete(key string) {
	c.mutex.Lock()
	index := slices.Index(c.queue, key)
	if index != -1 {
		if c.RemoveHook != nil {
			go c.RemoveHook(key, c.storage[key])
		}
		c.queue = append(c.queue[:index], c.queue[index+1:]...)
		delete(c.storage, key)
		delete(c.ttlMap, key)
	}
	c.mutex.Unlock()
}

// evict deletes entries that have expired
func (c *LRUCacheWithTTL) evict() {
	c.mutex.Lock()
	for key, ttl := range c.ttlMap {
		if time.Now().After(ttl) {
			index := slices.Index(c.queue, key)
			if c.RemoveHook != nil {
				go c.RemoveHook(key, c.storage[key])
			}
			c.queue = append(c.queue[:index], c.queue[index+1:]...)
			delete(c.storage, key)
			delete(c.ttlMap, key)
		}
	}
	c.mutex.Unlock()
}

func (c *LRUCacheWithTTL) Set(key string, value interface{}) {
	if value == nil {
		c.Delete(key)
		return
	}
	c.sendSet(key, value)
	c.set(key, value)
}

// Get gets a value from the cache
func (c *LRUCacheWithTTL) Get(key string) interface{} {
	c.evict()
	c.mutex.Lock()
	out, exists := c.storage[key]
	if exists {
		c.ttlMap[key] = time.Now().Add(c.TTL)
		c.mutex.Unlock()
	} else {
		c.mutex.Unlock()
		var err error
		if c.Filler != nil {
			out, err = c.Filler(key)
			if err != nil {
				log.Println(err)
			}
			c.Set(key, out)
		}
	}

	return out
}

// Delete deletes a value from the cache
// and sends the delete message to the other nodes
func (c *LRUCacheWithTTL) Delete(key string) {
	c.sendDelete(key)
	c.delete(key)
}

// Clean deletes all values from the cache
// and sends the sendClean message to the other nodes
func (c *LRUCacheWithTTL) Clean() {
	c.clean()
	c.sendClean()
}

// NewLRUCacheWithTTL creates a new LRUCacheWithTTL with the given name,
// address, maxEntries and TTL.
// It also starts a listener to receive messages from other nodes and starts a
// for cache eviction.
// name is the name of the cache
// address is the address of the cache
// maxEntries is the maximum number of entries that the cache can have
// ttl is the time-to-live for each entry in the cache
func NewLRUCacheWithTTL(name, broadcast, address string, maxEntries int, ttl time.Duration) *LRUCacheWithTTL {
	ctx, cancel := context.WithCancel(context.Background())
	c := &LRUCacheWithTTL{
		LRUCache: LRUCache{
			Cache: Cache{
				mutex:        sync.Mutex{},
				storage:      make(map[string]interface{}),
				Name:         name,
				Address:      address,
				Broadcast:    broadcast,
				context:      ctx,
				StopListener: cancel,
				node:         uuid.New(),
			},
			MaxEntries: maxEntries,
			queue:      make([]string, 0),
		},
		TTL:    ttl,
		ttlMap: make(map[string]time.Time),
	}

	go startListener(c, ctx)
	return c
}
