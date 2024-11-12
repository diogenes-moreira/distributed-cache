package distributed_cache

import (
	"context"
	"github.com/google/uuid"
	"log"
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
	StopListener context.CancelFunc // Cancel function
	// to stop the listener
	Filler func(string) (interface{}, error) // Function to fill the cache
	// when the key is not found
	RemoveHook func(string, interface{}) // Function to remove the key from the cache
	context    context.Context
	node       uuid.UUID
}

func (c *Cache) getNode() uuid.UUID {
	return c.node
}

func (c *Cache) getAddress() string {
	return c.Address
}

func (c *Cache) getName() string {
	return c.Name
}

func (c *Cache) clean() {
	c.mutex.Lock()
	if c.RemoveHook != nil {
		for key, value := range c.storage {
			go c.RemoveHook(key, value)
		}
	}
	c.storage = make(map[string]interface{})
	c.mutex.Unlock()
}

func (c *Cache) set(key string, value interface{}) {
	c.mutex.Lock()
	c.storage[key] = value
	c.mutex.Unlock()
}

func (c *Cache) delete(key string) {
	c.mutex.Lock()
	if c.RemoveHook != nil {
		go c.RemoveHook(key, c.storage[key])
	}
	delete(c.storage, key)
	c.mutex.Unlock()
}

// Set sets a value in the cache and sends it to the other nodes
func (c *Cache) Set(key string, value interface{}) {
	if value == nil {
		c.Delete(key)
		return
	}
	c.sendSet(key, value)
	c.set(key, value)
}

// Get gets a value from the cache
func (c *Cache) Get(key string) interface{} {
	var err error
	c.mutex.Lock()
	out, exists := c.storage[key]
	c.mutex.Unlock()
	if !exists && c.Filler != nil {
		out, err = c.Filler(key)
		if err != nil {
			log.Println(err)
		}
		c.Set(key, out)
	}
	return out
}

// Delete deletes a value from the cache
// and sends the delete message to the other nodes
func (c *Cache) Delete(key string) {
	c.sendDelete(key)
	c.delete(key)
}

// Clean deletes all values from the cache
// and sends the sendClean message to the other nodes
func (c *Cache) Clean() {
	c.clean()
	c.sendClean()
}

// NewCache creates a new Cache with the given name and address
// It also starts a listener to receive messages from other nodes
func NewCache(name, address string) *Cache {
	ctx, cancel := context.WithCancel(context.Background())
	c := &Cache{
		mutex:        sync.Mutex{},
		Name:         name,
		Address:      address,
		storage:      make(map[string]interface{}),
		StopListener: cancel,
		context:      ctx,
		node:         uuid.New(),
	}
	go startListener(c, ctx)
	return c
}
