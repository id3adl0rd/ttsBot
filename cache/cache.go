package cache

import (
	"container/list"
	"sync"
)

type Item struct {
	Key   string
	Value interface{}
}

type LRU struct {
	queue    *list.List
	mutex    *sync.RWMutex
	items    map[string]*list.Element
	capacity int
}

func NewLru(capacity int16) *LRU {
	return &LRU{
		queue:    list.New(),
		mutex:    &sync.RWMutex{},
		items:    make(map[string]*list.Element),
		capacity: int(capacity),
	}
}

func (c *LRU) Set(key string, value interface{}) bool {
	c.mutex.Lock()
	if element, exists := c.items[key]; exists == true {
		c.queue.MoveToFront(element)
		element.Value.(*Item).Value = value
		return true
	}
	c.mutex.Unlock()

	if c.queue.Len() == c.capacity {
		c.purge()
	}

	item := &Item{
		Key:   key,
		Value: value,
	}

	c.mutex.Lock()
	element := c.queue.PushFront(item)
	c.items[item.Key] = element
	c.mutex.Unlock()

	return true
}

func (c *LRU) purge() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if element := c.queue.Back(); element != nil {
		item := c.queue.Remove(element).(*Item)
		delete(c.items, item.Key)
	}
}

func (c *LRU) Get(key string) interface{} {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	element, exists := c.items[key]
	if exists == false {
		return nil
	}
	c.queue.MoveToFront(element)

	return element.Value.(*Item).Value
}
