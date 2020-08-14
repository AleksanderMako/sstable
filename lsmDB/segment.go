package lsmdbsegment

// TODO: figure out where to close the file
import (
	"bufio"
	"container/heap"
	"fmt"
	"lsm-tree/index"
	"lsm-tree/types"
	"os"
	"strings"
	"time"
)

// StringHeap is a min heap of string keys
type StringHeap []types.Node

func (h StringHeap) Len() int           { return len(h) }
func (h StringHeap) Less(i, j int) bool { return (strings.Compare(h[i].Key, h[j].Key) == -1) }
func (h StringHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

// Push adds a value to min heap
func (h *StringHeap) Push(x interface{}) {

	*h = append(*h, x.(types.Node))
}

// Pop pops the min from string heap
func (h *StringHeap) Pop() interface{} {

	old := *h
	heapLen := len(old)
	min := old[heapLen-1]
	*h = old[0 : heapLen-1]
	return min
}

// Segment repreents a db segment
type Segment struct {
	file            *os.File
	index           index.Index
	priorityQueue   *StringHeap
	queueSize       int
	logsDirectory   string
	segmentFileName string
}

// NewSegment constructs a new segment and returns a ptr to it
func NewSegment(indexImplementation index.Index) *Segment {

	pq := &StringHeap{}
	heap.Init(pq)
	logsDirectory := os.Getenv("LOGS_DIR")
	return &Segment{
		index:         indexImplementation,
		priorityQueue: pq,
		queueSize:     1000,
		logsDirectory: logsDirectory,
	}
}

// Put adds a key to the priority node
func (s *Segment) Put(node types.Node) bool {

	if len(*s.priorityQueue)+1 < s.queueSize {

		heap.Push(s.priorityQueue, node)
		return true
	}
	return false
}

// Commit commts the contents of the priority queue to disk
func (s *Segment) Commit() error {

	file, err := s.OpenFileinAppendOrCreate()
	if err != nil {
		return err
	}
	s.file = file

	sparseIndexOffset := 0

	for len(*s.priorityQueue) > 0 {

		value := heap.Pop(s.priorityQueue)
		node := value.(types.Node)
		fileStats, err := s.file.Stat()
		if err != nil {
			return err
		}

		if sparseIndexOffset == 0 || sparseIndexOffset%10 == 0 {
			s.index.Set(node.Key, fileStats.Size())
		}
		_, err = s.file.Write(append([]byte(node.Key+":"+node.Val), '\n'))
		if err != nil {
			return err
		}
	}
	s.index.SortIndex()
	s.file.Close()
	return nil
}

// SearchPriorityQueue searches for the key in the in memory priority queue
func (s *Segment) SearchPriorityQueue(key string) types.Node {

	if len(*s.priorityQueue) < 1 {
		return types.Node{}
	}

	l := 0
	r := len(*s.priorityQueue)

	for l <= r {
		mid := l + (r-l)/2
		midValue := (*s.priorityQueue)[mid].Key
		if key == midValue {
			return (*s.priorityQueue)[mid]
		}
		if key < midValue {
			r = mid - 1
		} else {
			l = mid + 1
		}
	}
	return types.Node{}
}

// SearchIndex serches the in momory map for a closer offest to the given key
// and then searches the segment file from the offest found in the map
func (s *Segment) SearchIndex(key string) (val string, found bool) {

	// TODO: implement an order map version

	file, err := s.OpenFileinReadOnly()
	if err != nil {
		return "", false
	}
	s.file = file
	defer s.file.Close()
	if offset := s.index.Get(key); offset != -1 {

		s.file.Seek(offset, 0)
		scanner := bufio.NewScanner(s.file)
		scanner.Scan()
		data := scanner.Text()
		colonIndex := strings.Index(data, ":")
		result := data[colonIndex+1:]
		if err := scanner.Err(); err != nil {
			return "", false
		}
		return result, true
	}

	lowerBound := s.index.LowerBound(key)

	lowerBoundOffset := s.index.Get(lowerBound)
	s.file.Seek(lowerBoundOffset, 0)
	scanner := bufio.NewScanner(s.file)
	result := ""
	for scanner.Scan() {
		data := scanner.Text()
		colonIndex := strings.Index(data, ":")
		currentKey := data[0:colonIndex]
		if currentKey == key {
			result = data[colonIndex+1:]
			return result, true
		}
	}
	return result, false
}

// OpenFileinAppendOrCreate opens or creates the current segment file in append mode
func (s *Segment) OpenFileinAppendOrCreate() (file *os.File, err error) {

	loc, _ := time.LoadLocation("UTC")
	fileName := fmt.Sprintf("segment %v", time.Now().In(loc))
	s.segmentFileName = fileName
	file, err = os.OpenFile(s.logsDirectory+"/"+fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777)

	if err != nil {
		// TODO: add some more info for the errs
		return nil, err
	}
	return file, nil
}

// OpenFileinReadOnly opens the current segment file in readonly mode
func (s *Segment) OpenFileinReadOnly() (file *os.File, err error) {

	file, err = os.OpenFile(s.logsDirectory+"/"+s.segmentFileName, os.O_RDONLY, 0444)

	if err != nil {
		// TODO: add some more info for the errs
		return nil, err
	}
	return file, nil
}
