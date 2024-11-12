package distributed_cache

import (
	"context"
	"sync"
)

// Cache is a simple cache interface in this file only you can find
// the Cache struct and the methods Set, Get, Delete and Clean
// that are used to interact with the cache

// Cache is a simple cache interface, to create a cache you must Use NewCache Method
type Cache struct {
	mutex        sync.Mutex
	storage      map[string]interface{}
	Name         string             // Name of the cache
	Address      string             // Port over which the cache will communicate
	StopListener context.CancelFunc // Cancel function to stop the listener
	context      context.Context
}

// Set sets a value in the cache and sends it to the other nodes
func (c *Cache) Set(key string, value interface{}) {
	if value == nil {
		c.Delete(key)
		return
	}
	c.sendSet(key, value)
	c.mutex.Lock()
	c.storage[key] = value
	c.mutex.Unlock()
}

// Get gets a value from the cache
func (c *Cache) Get(key string) interface{} {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.storage[key]
}

// Delete deletes a value from the cache
// and sends the delete message to the other nodes
func (c *Cache) Delete(key string) {
	c.sendDelete(key)
	c.mutex.Lock()
	delete(c.storage, key)
	c.mutex.Unlock()
}

// Clean deletes all values from the cache
// and sends the sendClean message to the other nodes
func (c *Cache) Clean() {
	c.mutex.Lock()
	c.storage = make(map[string]interface{})
	c.mutex.Unlock()
	c.sendClean()
}

// NewCache creates a new Cache with the given name and address
// It also starts a listener to receive messages from other nodes
func NewCache(name, address string) *Cache {
	ctx, cancel := context.WithCancel(context.Background())
	c := Cache{
		mutex:        sync.Mutex{},
		Name:         name,
		Address:      address,
		storage:      make(map[string]interface{}),
		StopListener: cancel,
		context:      ctx,
	}
	go c.startListener(ctx)
	return &c
}
