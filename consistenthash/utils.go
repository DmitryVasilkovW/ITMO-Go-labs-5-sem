package consistenthash

import "sort"

func removeFromSlice(slice []int, value int) []int {
	i := sort.Search(len(slice), func(i int) bool { return slice[i] >= value })
	if i < len(slice) && slice[i] == value {
		return append(slice[:i], slice[i+1:]...)
	}
	return slice
}
