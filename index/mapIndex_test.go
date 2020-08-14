package index

import (
	"testing"
)

func Test_LowerBound_ShoudReturnCorrectKey(t *testing.T) {

	// Arrange
	mapIndex := NewMapIndex()

	mapIndex.Set("aaa", int64(9))
	mapIndex.Set("bbb", int64(9))
	mapIndex.Set("hhhh", int64(9))
	mapIndex.Set("kk", int64(9))
	mapIndex.Set("ww", int64(9))
	mapIndex.Set("yy", int64(9))

	mapIndex.SortIndex()

	// Act
	lowerBound := mapIndex.LowerBound("cccccccccc")
	lowerBound1 := mapIndex.LowerBound("d")
	lowerBound2 := mapIndex.LowerBound("eeee")
	lowerBound3 := mapIndex.LowerBound("1")
	lowerBound4 := mapIndex.LowerBound("zzz")

	// Assert

	if lowerBound != "bbb" {
		t.Error("expected bbb but got " + lowerBound)
	}
	if lowerBound1 != "bbb" {
		t.Error("expected bbb but got " + lowerBound1)
	}
	if lowerBound2 != "bbb" {
		t.Error("expected bbb but got " + lowerBound2)
	}

	if lowerBound3 != mapIndex.Begin() {
		t.Error("expected mapIndex.Begin() but got " + lowerBound2)
	}
	if lowerBound4 != mapIndex.End() {
		t.Error("expected mapIndex.Begin() but got " + lowerBound2)
	}
}

func Test_LowerBound_SortedIndex_HasCorrectOffesets(t *testing.T) {

	// Arrange
	mapIndex := NewMapIndex()

	mapIndex.Set("aaa", int64(9))
	mapIndex.Set("bbb", int64(9))
	mapIndex.Set("hhhh", int64(9))
	mapIndex.Set("kk", int64(9))
	mapIndex.Set("ww", int64(9))
	mapIndex.Set("yy", int64(1000))

	// Act
	mapIndex.SortIndex()

	//Assert

	if end := mapIndex.Get("end"); end == -1 {
		t.Error("end was not found t should have existed")
	} else if end != 1000 {
		t.Error("end should have the value 1000 same as the last key yy")
	}

	if begin := mapIndex.Get("begin"); begin == -1 {
		t.Error("begin was not found t should have existed")
	} else if begin != 0 {
		t.Error("begin should always have the offest 0 ")
	}

}
