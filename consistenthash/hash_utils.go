package consistenthash

import (
	"crypto/sha256"
	"fmt"
	"sort"
)

func generateVirtualNodeID(nodeID string, replica int) string {
	return fmt.Sprintf("%s#%d", nodeID, replica)
}

func hashKey(key string) int {
	hash := sha256.Sum256([]byte(key))
	return int(hash[0])<<24 | int(hash[1])<<16 | int(hash[2])<<8 | int(hash[3])
}

func (h *ConsistentHash[N]) addHashToRing(hash int, node *N) {
	h.hashRing = insertInOrder(h.hashRing, hash)
	h.hashToNode[hash] = node
}

func (h *ConsistentHash[N]) removeHashFromRing(hash int) {
	h.hashRing = removeFromOrder(h.hashRing, hash)
	delete(h.hashToNode, hash)
}

func findInsertionIndex(slice []int, value int) int {
	return sort.Search(len(slice), func(i int) bool { return slice[i] >= value })
}

func findClosestHash(ring []int, hash int) int {
	index := findInsertionIndex(ring, hash)
	if index == len(ring) {
		return 0
	}
	return index
}

func extendSlice(slice []int) []int {
	return append(slice, 0)
}

func shiftSlice(slice []int, fromIndex int) {
	copy(slice[fromIndex+1:], slice[fromIndex:])
}

func insertInOrder(slice []int, value int) []int {
	index := findInsertionIndex(slice, value)
	slice = extendSlice(slice)
	shiftSlice(slice, index)
	slice[index] = value
	return slice
}

func isValueAt(slice []int, index int, value int) bool {
	return index < len(slice) && slice[index] == value
}

func removeValueAt(slice []int, index int) []int {
	return append(slice[:index], slice[index+1:]...)
}

func removeFromOrder(slice []int, value int) []int {
	index := findInsertionIndex(slice, value)
	if isValueAt(slice, index, value) {
		return removeValueAt(slice, index)
	}
	return slice
}
