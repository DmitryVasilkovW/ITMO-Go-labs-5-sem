//go:build !solution

package genericsum

import (
	"cmp"
	"math/cmplx"
	"slices"

	"golang.org/x/exp/constraints"
)

func Min[T cmp.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func SortSlice[S ~[]E, E cmp.Ordered](a S) {
	slices.Sort(a)
}

func MapsEqual[M ~map[K]V, K comparable, V comparable](a, b M) bool {
	if len(a) != len(b) {
		return false
	}

	for k, v := range a {
		if bValue, ok := b[k]; !ok || bValue != v {
			return false
		}
	}

	return true
}

func SliceContains[T comparable](s []T, v T) bool {
	for _, elem := range s {
		if elem == v {
			return true
		}
	}
	return false
}

func MergeChans[C <-chan E, E any](chs ...C) C {
	result := make(chan E)

	go func() {
		openedChans := len(chs)

		for openedChans > 0 {
			openedChans = processChannels(chs, result, openedChans)
		}

		close(result)
	}()

	return result
}

func processChannels[C <-chan E, E any](chs []C, result chan<- E, openedChans int) int {
	for _, ch := range chs {
		if !handleChannel(ch, result) {
			openedChans--
		}
	}
	return openedChans
}

func handleChannel[C <-chan E, E any](ch C, result chan<- E) bool {
	select {
	case value, isOpened := <-ch:
		if isOpened {
			result <- value
			return true
		}
		return false
	default:
		return true
	}
}

func IsHermitianMatrix[T constraints.Integer | constraints.Complex | constraints.Float](m [][]T) bool {
	if !isSquareMatrix(m) {
		return false
	}

	for y := 0; y < len(m); y++ {
		for x := 0; x < len(m[0]); x++ {
			if !isHermitianPair(m, x, y) {
				return false
			}
		}
	}

	return true
}

func isSquareMatrix[T any](m [][]T) bool {
	height := len(m)
	if height < 1 {
		return true
	}
	width := len(m[0])
	if width < 1 {
		return true
	}
	return height == width
}

func isHermitianPair[T constraints.Integer | constraints.Complex | constraints.Float](m [][]T, x, y int) bool {
	el := m[x][y]
	switch any(el).(type) {
	case complex64:
		return checkComplex64Pair(m, x, y)
	case complex128:
		return checkComplex128Pair(m, x, y)
	default:
		return m[x][y] == m[y][x]
	}
}

func checkComplex64Pair[T any](m [][]T, x, y int) bool {
	opposite, ok1 := any(m[y][x]).(complex64)
	cur, ok2 := any(m[x][y]).(complex64)
	if !ok1 || !ok2 {
		return false
	}
	return real(opposite) == real(cur) && imag(opposite) == -imag(cur)
}

func checkComplex128Pair[T any](m [][]T, x, y int) bool {
	opposite, ok1 := any(m[y][x]).(complex128)
	cur, ok2 := any(m[x][y]).(complex128)
	if !ok1 || !ok2 {
		return false
	}
	return cmplx.Conj(cur) == opposite
}
