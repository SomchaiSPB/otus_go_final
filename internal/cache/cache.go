package imagecache

import "sync"

type Key string

func (k Key) String() string {
	return string(k)
}

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	mu       sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type cacheItem struct {
	key   string
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	newCacheItem := &cacheItem{
		key:   key.String(),
		value: value,
	}
	if _, ok := c.items[key]; ok {
		c.items[key].Value = newCacheItem

		return true
	}

	if c.capacity == c.queue.Len() {
		lastItem := c.queue.Back()
		key := lastItem.Value.(*cacheItem).key
		delete(c.items, Key(key))
		c.queue.Remove(lastItem)
	}

	listNode := c.queue.PushFront(newCacheItem)
	c.items[key] = listNode

	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	listNode, ok := c.items[key]

	if !ok {
		return nil, false
	}

	c.queue.MoveToFront(listNode)

	return listNode.Value.(*cacheItem).value, ok
}

func (c *lruCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[Key]*ListItem, c.capacity)
	c.queue = NewList()
}
