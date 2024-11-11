package distributed_cache

// Cache is a simple cache interface in this file only you can find
// the Cache struct and the methods Set, Get, Delete and Clean
// that are used to interact with the cache

// Cache is a simple cache interface
type Cache struct {
	storage map[string]interface{}
	Name    string
	Address string
}

// Set sets a value in the cache and sends it to the other nodes
func (c *Cache) Set(key string, value interface{}) {
	if value == nil {
		c.Delete(key)
		return
	}
	c.sendSet(key, value)
	c.storage[key] = value
}

// Get gets a value from the cache
func (c *Cache) Get(key string) interface{} {
	return c.storage[key]
}

// Delete deletes a value from the cache
// and sends the delete message to the other nodes
func (c *Cache) Delete(key string) {
	c.sendDelete(key)
	delete(c.storage, key)
}

// Clean deletes all values from the cache
// and sends the clean message to the other nodes
func (c *Cache) Clean() {
	c.clean()
}

func NewCache(name, address string) *Cache {
	c := &Cache{
		Name:    name,
		Address: address,
		storage: make(map[string]interface{}),
	}
	go c.startListener()
	return c
}
