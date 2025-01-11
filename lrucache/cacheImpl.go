package lrucache

import "container/list"

type LRU struct {
	capacity int
	cache    map[int]*list.Element
	usages   *list.List
}

type entry struct {
	key   int
	value int
}

func (c *LRU) Get(key int) (int, bool) {
	if value, found := c.cache[key]; found {
		updateUsage(c, value)
		return getValue(value), true
	}

	return -1, false
}

func (c *LRU) Set(key, value int) {
	if c.capacity == 0 {
		return
	}

	if val, found := c.cache[key]; found {
		setValue(val, value)
		updateUsage(c, val)
		return
	}

	c.updateUsage()
	c.updateCache(key, value)
}

func (c *LRU) updateCache(key, value int) {
	c.cache[key] = c.usages.PushFront(
		&entry{
			key:   key,
			value: value,
		})
}

func (c *LRU) updateUsage() {
	if c.usages.Len() >= c.capacity {
		if elem := c.usages.Back(); elem != nil {
			c.usages.Remove(elem)
			delete(c.cache, elem.Value.(*entry).key)
		}

		return
	}
}

func updateUsage(c *LRU, listValue *list.Element) {
	c.usages.MoveToFront(listValue)
}

func getValue(listValue *list.Element) int {
	return listValue.Value.(*entry).value
}

func setValue(listValue *list.Element, value int) {
	listValue.Value.(*entry).value = value
}

func (c *LRU) Range(f func(key, value int) bool) {
	for e := c.usages.Back(); e != nil; e = e.Prev() {
		val := e.Value.(*entry)
		if !f(val.key, val.value) {
			break
		}
	}
}

func (c *LRU) Clear() {
	c.cache = make(map[int]*list.Element, c.capacity)
	c.usages.Init()
}
