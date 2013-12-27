// Package biscuit is used for simple linguistic computations.
package biscuit

import (
	"sort"
)

// SortedMap is a structure that contains the original map and associated
// sorted keys.
type SortedMap struct {
	m map[string]float64
	s []string
}

func (sm *SortedMap) Len() int {
	return len(sm.m)
}

func (sm *SortedMap) Less(i, j int) bool {
	return sm.m[sm.s[i]] > sm.m[sm.s[j]]
}

func (sm *SortedMap) Swap(i, j int) {
	sm.s[i], sm.s[j] = sm.s[j], sm.s[i]
}

// SortedKeys takes the specified map and returns a sorted array containing the
// original map's keys. The order in which the keys are returned are determined
// by their associated values in the original map.
func SortedKeys(m map[string]float64) []string {
	sm := new(SortedMap)
	sm.m = m
	sm.s = make([]string, len(m))
	i := 0
	for key := range m {
		sm.s[i] = key
		i++
	}
	sort.Sort(sm)
	return sm.s
}
