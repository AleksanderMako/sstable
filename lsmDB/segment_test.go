package lsmdbsegment

import (
	"bufio"
	"fmt"
	"lsm-tree/index"
	"lsm-tree/types"
	"os"
	"strings"
	"testing"
)

func Test_Commit_Should_DumpDataSortedInFile(t *testing.T) {

	// Arrange
	logsDir := "/home/alex/Desktop/go/src/lsm-tree/"
	os.Setenv("LOGS_DIR", logsDir+"/logs")
	mapIndex := index.NewMapIndex()
	segment := NewSegment(mapIndex)

	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("%v", i)
		segment.Put(types.Node{
			Key: key,
			Val: key,
		})
	}

	// Act
	segment.Commit()

	// Assert

	file, err := segment.OpenFileinReadOnly()
	if err != nil {
		t.Error(err.Error())
	}

	scanner := bufio.NewScanner(file)
	fileData := []string{}
	for scanner.Scan() {
		data := scanner.Text()
		colonIndex := strings.Index(data, ":")
		currentKey := data[0:colonIndex]
		fileData = append(fileData, currentKey)
	}

	for i := 0; i < len(fileData); i++ {
		if i == 0 {
			continue
		}
		if fileData[i] < fileData[i-1] {
			t.Errorf("slice is not sorted val at i is %v  , val at i-1 is %v", fileData[i], fileData[i-1])
		}
	}
}

func Test_SearchPriorityQueue_Should_Find_Existing_Key(t *testing.T) {

	// Arrange
	logsDir := "/home/alex/Desktop/go/src/lsm-tree/"
	os.Setenv("LOGS_DIR", logsDir+"/logs")
	mapIndex := index.NewMapIndex()
	segment := NewSegment(mapIndex)
	for i := 0; i < 800; i++ {
		key := fmt.Sprintf("%v", i)
		segment.Put(types.Node{
			Key: key,
			Val: key,
		})
	}

	// Act

	node := segment.SearchPriorityQueue("564")

	//Assert

	if node.Key != "564" || node.Val != "564" {
		t.Errorf("expecetd key %v , val %v but got key %v , val %v", "564", "564", node.Key, node.Val)
	}
}

func Test_SearchPriorityQueue_Should_Return_EmptyNode_WhenKey_DoesNotExist(t *testing.T) {

	// Arrange
	logsDir := "/home/alex/Desktop/go/src/lsm-tree/"
	os.Setenv("LOGS_DIR", logsDir+"/logs")
	mapIndex := index.NewMapIndex()
	segment := NewSegment(mapIndex)
	for i := 0; i < 800; i++ {
		key := fmt.Sprintf("%v", i)
		segment.Put(types.Node{
			Key: key,
			Val: key,
		})
	}

	// Act

	node := segment.SearchPriorityQueue("1000")

	//Assert

	if node.Key != "" || node.Val != "" {
		t.Errorf("expecetd key %v , val %v but got key %v , val %v", "empty", "empty", node.Key, node.Val)
	}
}
