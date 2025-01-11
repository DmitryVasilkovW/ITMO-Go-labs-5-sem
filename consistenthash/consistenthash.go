//go:build !solution

package consistenthash

import (
	"sync"
)

type Node interface {
	ID() string
}

type ConsistentHash[N Node] struct {
	virtualReplicas int
	hashRing        []int
	hashToNode      map[int]*N
	mutex           sync.RWMutex
}

func New[N Node]() *ConsistentHash[N] {
	return &ConsistentHash[N]{
		hashRing:        []int{},
		virtualReplicas: 94,
		hashToNode:      make(map[int]*N),
	}
}

func (h *ConsistentHash[N]) AddNode(node *N) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	for i := 0; i < h.virtualReplicas; i++ {
		virtualNodeID := generateVirtualNodeID((*node).ID(), i)
		hash := hashKey(virtualNodeID)
		h.addHashToRing(hash, node)
	}
}

func (h *ConsistentHash[N]) GetNode(key string) *N {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	if len(h.hashRing) == 0 {
		return nil
	}

	hash := hashKey(key)
	index := findClosestHash(h.hashRing, hash)
	return h.hashToNode[h.hashRing[index]]
}

func (h *ConsistentHash[N]) RemoveNode(node *N) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	for i := 0; i < h.virtualReplicas; i++ {
		virtualNodeID := generateVirtualNodeID((*node).ID(), i)
		hash := hashKey(virtualNodeID)
		h.removeHashFromRing(hash)
	}
}
