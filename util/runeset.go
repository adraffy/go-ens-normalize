package util

import (
	"slices"
	"sort"
)

type RuneSet struct {
	sorted []rune
}

func NewRuneSetFromInts(v []int) RuneSet {
	sorted := make([]rune, len(v))
	for i, x := range v {
		sorted[i] = rune(x)
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i] < sorted[j]
	})
	return RuneSet{sorted}
}

func NewRuneSetFromKeys[T any](m map[rune]T) RuneSet {
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
