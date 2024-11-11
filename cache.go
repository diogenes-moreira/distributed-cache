package distributed_cache

type Cache struct {
	storage map[string]interface{}
	Name    string
	Address string
}

func (c *Cache) Set(key string, value interface{}) {
	if value == nil {
		c.Delete(key)
		return
	}
	c.sendSet(key, value)
	c.storage[key] = value
}

func (c *Cache) Get(key string) interface{} {
	return c.storage[key]
}

func (c *Cache) Delete(key string) {
	c.sendDelete(key)
	delete(c.storage, key)
}

func (c *Cache) Clean() {
	c.clean()
}

func NewCache(name, address string) *Cache {
	c := &Cache{
		Name:    name,
		Address: address,
		storage: make(map[string]interface{}),
	}
	go c.StartListener()
	return c
}
