package clock

import (
	"container/heap"
	"time"
)

type executer interface {
	execute() *task
}

type task struct {
	deadline  time.Time
	owner     executer
	heapIndex int
}

const removed = -1

func newTask(d time.Time, e executer) *task {
	return &task{
		deadline:  d,
		owner:     e,
		heapIndex: removed,
	}
}

func (t task) hasStopped() bool {
	return t.heapIndex == removed
}

// taskHeap implements mockTimers with a heap.
type taskHeap []*task

func newTaskHeap() *taskHeap {
	res := make(taskHeap, 0, 1024)
	return &res
}

func (h taskHeap) Len() int { return len(h) }

func (h taskHeap) Less(i, j int) bool {
	return h[i].deadline.Before(h[j].deadline)
}

func (h taskHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].heapIndex = i
	h[j].heapIndex = j
}

func (h *taskHeap) Push(x interface{}) {
	n := len(*h)
	t := x.(*task)
	t.heapIndex = n
	*h = append(*h, t)
}

func (h *taskHeap) Pop() interface{} {
	old := *h
	n := len(old)
	t := old[n-1]
	t.heapIndex = removed
	*h = old[0 : n-1]
	return t
}

func (h *taskHeap) start(t *task) {
	heap.Push(h, t)
}

func (h *taskHeap) stop(t *task) {
	if !t.hasStopped() {
		heap.Remove(h, t.heapIndex)
	}
}

func (h *taskHeap) reset(t *task) {
	if !t.hasStopped() {
		heap.Fix(h, t.heapIndex)
	} else {
		heap.Push(h, t)
	}
}

func (h taskHeap) next() *task {
	if len(h) == 0 {
		return nil
	}
	return h[0]
}
