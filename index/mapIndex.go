package index

import (
	"sort"
)

type Index interface {
	Get(key string) int64
	Set(key string, val int64)
	LowerBound(key string) string
	SortIndex()
	Begin() string
	End() string
}

// MapIndex is index structure based on simple hashmap
type MapIndex struct {
	index       map[string]int64
	sortedIndex []string
}

// NewMapIndex returns an implementation of index data structure
func NewMapIndex() Index {

	index := make(map[string]int64)
	sortedIndex := []string{}
	return &MapIndex{
		index:       index,
		sortedIndex: sortedIndex,
	}
}
func (mi *MapIndex) Get(key string) int64 {

	if val, ok := mi.index[key]; ok {
		return val
	}

	return -1
}

func (mi *MapIndex) Set(key string, val int64) {
	mi.index[key] = val
	mi.sortedIndex = append(mi.sortedIndex, key)
}

// LowerBound returns the largest key considered to go before the given key
func (mi *MapIndex) LowerBound(key string) string {

	l := 1
	r := len(mi.sortedIndex) - 2

	lowerBoundIndex := len(mi.sortedIndex) - 1
	for l <= r {

		mid := l + ((r - l) / 2)

		if key < mi.sortedIndex[mid] {
			lowerBoundIndex = mid
			r = mid - 1
		} else if key > mi.sortedIndex[mid] {
			l = mid + 1
		}
	}

	if lowerBoundIndex == len(mi.sortedIndex)-1 {
		return mi.End()
	}

	return mi.sortedIndex[lowerBoundIndex-1]

}

func (mi *MapIndex) SortIndex() {
	boundedSortedIndex := []string{"begin"}
	sort.Strings(mi.sortedIndex)
	boundedSortedIndex = append(boundedSortedIndex, mi.sortedIndex...)
	boundedSortedIndex = append(boundedSortedIndex, "end")
	mi.sortedIndex = boundedSortedIndex
	mi.index["begin"] = 0
	mi.index["end"] = mi.index[mi.sortedIndex[len(mi.sortedIndex)-2]]
}

func (mi *MapIndex) Begin() string {

	return mi.sortedIndex[0]
}

func (mi *MapIndex) End() string {

	return mi.sortedIndex[len(mi.sortedIndex)-1]
}
