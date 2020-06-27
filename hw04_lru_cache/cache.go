package hw04_lru_cache //nolint:golint,stylecheck

import "sync"

// Cache ...
type Cache interface {
	Set(key string, value interface{}) bool
	Get(key string) (interface{}, bool)
	Clear()
}

type lruCache struct {
	cap      int
	queue    List
	cacheMap map[string]cacheItem
	mux      sync.Mutex
}

func (c *lruCache) Set(key string, value interface{}) bool {
	c.mux.Lock()
	defer c.mux.Unlock()
	if item, ok := c.cacheMap[key]; ok {
		c.queue.MoveToFront(item.qPointer)
		item.val = value
		c.cacheMap[key] = item

		return ok
	}
	cItem := cacheItem{
		val:      value,
		qPointer: c.queue.PushFront(key),
	}
	c.cacheMap[key] = cItem
	c.rotate()

	return false
}

func (c *lruCache) Get(key string) (interface{}, bool) {
	c.mux.Lock()
	defer c.mux.Unlock()
	v, ok := c.cacheMap[key]
	if ok {
		c.queue.MoveToFront(v.qPointer)
	}
	return v.val, ok
}

func (c *lruCache) rotate() {
	if c.queue.Len() <= c.cap {
		return
	}
	elem := c.queue.Back()
	if elem != nil {
		delete(c.cacheMap, elem.Value.(string))
		c.queue.Remove(elem)
	}
}

func (c *lruCache) Clear() {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.queue = NewList()
	c.cacheMap = make(map[string]cacheItem)
}

type cacheItem struct {
	val      interface{}
	qPointer *listItem
}

// NewCache ...
func NewCache(capacity int) Cache {
	return &lruCache{
		cap:      capacity,
		queue:    NewList(),
		cacheMap: make(map[string]cacheItem),
	}
}
