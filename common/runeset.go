package common

import (
	"slices"
	"sort"
)

// jesus christ go is shit
// func CompareArrays[T comparable](a, b []T) int {
func CompareRunes(a, b []rune) int {
	c := len(a) - len(b)
	if c != 0 {
		return c
	}
	for i, aa := range a {
		switch {
		case aa < b[i]:
			return -1
		case aa > b[i]:
			return 1
		}
	}
	return 0
}

func ToRuneSet(v []int) map[rune]bool {
	set := make(map[rune]bool)
	for _, x := range v {
		set[rune(x)] = true
	}
	return set
}

func ToRuneArray(v []int) []rune {
	runes := make([]rune, len(v))
	for i, x := range v {
		runes[i] = rune(x)
	}
	return runes
}

type RuneSet struct {
	sorted []rune
}

func RuneSetFromInts(v []int) RuneSet {
	sorted := make([]rune, len(v))
	for i, x := range v {
		sorted[i] = rune(x)
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i] < sorted[j]
	})
	return RuneSet{sorted}
}

func RuneSetFromKeys[T any](m map[rune]T) RuneSet {
	sorted := make([]rune, 0, len(m))
	for x := range m {
		sorted = append(sorted, x)
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i] < sorted[j]
	})
	return RuneSet{sorted}
}

func (set RuneSet) Contains(cp rune) bool {
	_, exists := slices.BinarySearch(set.sorted, cp)
	return exists
}
func (set RuneSet) Size() int {
	return len(set.sorted)
}
func (set RuneSet) Filter(fn func(cp rune) bool) RuneSet {
	v := make([]rune, 0, len(set.sorted))
	for _, x := range set.sorted {
		if fn(x) {
			v = append(v, x)
		}
	}
	return RuneSet{v}
}
func (set RuneSet) ToArray() []rune {
	v := make([]rune, len(set.sorted))
	copy(v, set.sorted)
	return v
}

type ReadOnlySet[K comparable, V any] struct {
	m map[K]V
}

func MakeReadOnlySet[K comparable, V any](m map[K]V) ReadOnlySet[K, V] {
	return ReadOnlySet[K, V]{m}
}

func (set ReadOnlySet[K, V]) Size() int {
	return len(set.m)
}

func (set ReadOnlySet[K, V]) Contains(k K) bool {
	_, ok := set.m[k]
	return ok
}

func (set ReadOnlySet[K, V]) ToArray() []K {
	v := make([]K, 0, len(set.m))
	for k := range set.m {
		v = append(v, k)
	}
	return v
}
