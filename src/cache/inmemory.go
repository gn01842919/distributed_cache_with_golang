package cache

import (
	"sync"
)

type inMemoryCache struct {
	c     map[string][]byte
	mutex sync.RWMutex
	Stat  // TODO: Why not pointer here but uses pointer in cacheHandler in http/cache.go??
}

func (c *inMemoryCache) Set(k string, v []byte) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	tmp, exist := c.c[k]
	if exist {
		c.del(k, tmp) // Seems to make Stat correct
	}
	c.c[k] = v
	c.add(k, v)
	return nil
}

func (c *inMemoryCache) Del(k string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	v, exist := c.c[k]
	if exist {
		delete(c.c, k)
		c.del(k, v)
	}
	return nil
}

func (c *inMemoryCache) Get(k string) ([]byte, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.c[k], nil
}

func (c *inMemoryCache) GetStat() Stat {
	return c.Stat
}

func newInMemoryCache() *inMemoryCache {
	// Solved: So we can return a reference to a local variable?
	// ==> Yes. In go, this will be saved in heap. It is same as using new()
	return &inMemoryCache{
		make(map[string][]byte), sync.RWMutex{}, Stat{}}
}
