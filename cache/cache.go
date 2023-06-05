package cache

import (
	"container/list"
	"image/gif"
	"sync"
)

type Cache struct {
	size     int
	capacity int
	ll       *list.List
	cache    map[string]*list.Element
	lock     sync.RWMutex
}

type entry struct {
	key   string
	value gif.GIF
}

// NewCache creates a new Cache with the given capacity.
func NewCache(capacity int) *Cache {
	return &Cache{
		capacity: capacity,
		ll:       list.New(),
		cache:    make(map[string]*list.Element),
	}
}

// Fetch gets the value for a key. If the key does not exist, it returns nil.
func (c *Cache) Fetch(key string) (gif.GIF, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	if element, ok := c.cache[key]; ok {
		c.ll.MoveToFront(element)
		return element.Value.(*entry).value, true
	}

	return gif.GIF{}, false
}

// Store sets the value for a key.
func (c *Cache) Store(key string, value gif.GIF) {
	c.lock.Lock()
	defer c.lock.Unlock()

	// If the key already exists, just update the value and move the element to the front.
	if element, ok := c.cache[key]; ok {
		c.ll.MoveToFront(element)
		element.Value.(*entry).value = value
		return
	}

	// Add new item
	element := c.ll.PushFront(&entry{key, value})
	c.cache[key] = element
	c.size++

	// Verify if cache exceeded its capacity
	if c.size > c.capacity {
		c.removeOldest()
	}
}

// removeOldest removes the oldest (least recently used) item from the cache.
func (c *Cache) removeOldest() {
	element := c.ll.Back()
	if element != nil {
		c.ll.Remove(element)
		key := element.Value.(*entry).key
		delete(c.cache, key)
		c.size--
	}
}
