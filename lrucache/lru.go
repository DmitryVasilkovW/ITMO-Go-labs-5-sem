//go:build !solution

package lrucache

import "container/list"

func New(cap int) LRUCache {
	return &LRU{
		capacity: cap,
		cache:    make(map[int]*list.Element, cap),
		usages:   list.New(),
	}
}
