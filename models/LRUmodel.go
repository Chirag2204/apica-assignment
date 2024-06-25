package models

import (
	"log"
	"sync"
	"time"
)

type cacheItem struct {
	value      interface{}
	expiration int64
}

type LRUCache struct {
	capacity int
	items    map[string]*cacheItem
	order    []string
	mutex    sync.Mutex
}

func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		items:    make(map[string]*cacheItem),
		order:    make([]string, 0, capacity),
	}
}

func (c *LRUCache) Get(key string) (interface{}, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	item, exists := c.items[key]
	if !exists {
		log.Println("Get:", key, "does not exist")
		return nil, false
	}
	if time.Now().Unix() > item.expiration {
		log.Println("Get:", key, "has expired")
		delete(c.items, key)
		c.removeFromOrder(key)
		return nil, false
	}

	log.Println("Get:", key, "is found with value", item.value)
	c.updateOrder(key)
	return item.value, true
}

func (c *LRUCache) Set(key string, value interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
    
	if _, exists := c.items[key]; !exists && len(c.items) >= c.capacity {
		delete(c.items, c.order[0])
		c.order = c.order[1:]
	}

	expiration := time.Now().Add(50 * time.Second).Unix()
	c.items[key] = &cacheItem{value: value, expiration: expiration}
	c.updateOrder(key)
	log.Printf("%s added successfully", key)
}

func (c *LRUCache) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, exists := c.items[key]; exists {
		delete(c.items, key)
		c.removeFromOrder(key)
	}
}

func (c *LRUCache) updateOrder(key string) {
	c.removeFromOrder(key)
	c.order = append(c.order, key)
}

func (c *LRUCache) removeFromOrder(key string) {
	for i, v := range c.order {
		if v == key {
			c.order = append(c.order[:i], c.order[i+1:]...)
			break
		}
	}
}
