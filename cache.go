package main

import (
	"container/list"
	"sync"
)

type Item struct {
	Key   string
	Value interface{}
}

type LRU struct {
	queue *list.List
	items map[string]*list.Element
	*sync.Mutex
	capacity int
}

func NewLru(capacity int) *LRU {
	return &LRU{
		capacity: capacity,
		items:    make(map[string]*list.Element),
		queue:    list.New(),
	}
}

func (c *LRU) Set(key string, value interface{}) bool {
	var mutex sync.Mutex

	mutex.Lock()
	if element, exists := c.items[key]; exists == true {
		c.queue.MoveToFront(element)
		element.Value.(*Item).Value = value
		return true
	}

	if c.queue.Len() == c.capacity {
		c.purge()
	}

	item := &Item{
		Key:   key,
		Value: value,
	}

	element := c.queue.PushFront(item)
	c.items[item.Key] = element
	mutex.Unlock()

	return true
}

func (c *LRU) purge() {
	var mutex sync.Mutex
	mutex.Lock()
	if element := c.queue.Back(); element != nil {
		item := c.queue.Remove(element).(*Item)
		delete(c.items, item.Key)
	}
	mutex.Unlock()
}

func (c *LRU) Get(key string) interface{} {
	var mutex sync.Mutex
	mutex.Lock()
	element, exists := c.items[key]
	if exists == false {
		return nil
	}
	c.queue.MoveToFront(element)
	mutex.Unlock()

	return element.Value.(*Item).Value
}
