package cache

import (
	"sync"
	"time"
)

type inMemoryCache struct {
	c     map[string][]byte
	mutex sync.RWMutex
	Stat  // TODO: Why not pointer here but uses pointer in cacheHandler in http/cache.go??
	ch    chan *pair
}

type pair struct {
	k string
	v []byte
}

func (c *inMemoryCache) Set(k string, v []byte) error {
	c.ch <- &pair{k, v}
	return nil
}

// BatchSize : the size a batch of k-v pairs to write to cache at once
const BatchSize = 100

func (c *inMemoryCache) flushBatch(keys [BatchSize]string, vals [BatchSize][]byte, count int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for i := 0; i < count; i++ {
		k := keys[i]
		v := vals[i]
		tmp, exist := c.c[k]
		if exist {
			c.del(k, tmp)
		}
		c.c[k] = v
		c.add(k, v) // add stat
	}
}

func (c *inMemoryCache) writeFunc(ch chan *pair) {
	count := 0
	t := time.NewTimer(time.Second)
	var keys [BatchSize]string
	var vals [BatchSize][]byte
	for {
		select {
		case p := <-ch:
			keys[count] = p.k
			vals[count] = p.v
			count++
			if count == BatchSize {
				c.flushBatch(keys, vals, count)
				count = 0
			}
			if !t.Stop() { // t.Stop() returns false when the timer has been triggered
				<-t.C // remove an item for timer's channel to avoid going the the second case at once
			}
			t.Reset(time.Second)
		case <-t.C:
			if count != 0 {
				c.flushBatch(keys, vals, count)
				count = 0
			}
			t.Reset(time.Second)
		}
	}
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
	ch := make(chan *pair, 5000)
	c := &inMemoryCache{
		make(map[string][]byte), sync.RWMutex{}, Stat{}, ch}
	go c.writeFunc(ch)
	return c
}
