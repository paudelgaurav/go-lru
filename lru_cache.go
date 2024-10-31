package lrucache

import (
	"container/list"
	"sync"
)

type ValueNode struct {
	value any
	node  *list.Element
}

type LRUCache struct {
	keys        *list.List
	data        map[any]ValueNode
	maxCapacity int
	mu          sync.RWMutex
}

func NewCache(maxCapacity int) LRUCache {
	emptyMap := make(map[any]ValueNode)
	l := list.New()
	return LRUCache{keys: l, data: emptyMap, maxCapacity: maxCapacity}
}

func (c *LRUCache) Get(key any) (any, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	valueNode, ok := c.data[key]
	if !ok {
		return nil, false
	}

	// move the node to the front
	c.keys.MoveToFront(valueNode.node)

	return valueNode.value, true

}

func (c *LRUCache) Add(key, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	valueNode, ok := c.data[key]
	// if key exists
	if ok {
		valueNode.value = value
		c.data[key] = valueNode
		c.keys.MoveToFront(valueNode.node)

		return
	}

	// if key does not exists we need to check for the length and remove the least used one

	if c.keys.Len() == c.maxCapacity {
		if lastElem := c.keys.Back(); lastElem != nil {
			delete(c.data, lastElem.Value)
			c.keys.Remove(lastElem)
		}
	}

	node := c.keys.PushFront(key)
	newValueNode := ValueNode{
		value: value,
		node:  node,
	}
	c.data[key] = newValueNode

}

func (c *LRUCache) Remove(key any) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	valueNode, ok := c.data[key]
	if !ok {
		return false
	}

	c.keys.Remove(valueNode.node)
	delete(c.data, key)

	return true
}

func (c *LRUCache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.keys.Len() // Wait for all goroutines to finish

}

func (c *LRUCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.keys.Init()
	for k := range c.data {
		delete(c.data, k)
	}
}
