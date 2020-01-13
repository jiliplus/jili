package clock

import (
	"container/heap"
	"time"
)

// TODO: 重命名这个结构体
type mockTimer struct {
	deadline time.Time
	fire     func() time.Duration
	mock     *Mock
	index    int
}

const removed = -1

func newTask(m *Mock, d time.Time) *mockTimer {
	return &mockTimer{
		deadline: d,
		mock:     m,
		index:    removed,
	}
}

func (t mockTimer) hasStopped() bool {
	return t.index == removed
}

// taskHeap implements mockTimers with a heap.
type taskHeap []*mockTimer

func (h taskHeap) Len() int { return len(h) }

func (h taskHeap) Less(i, j int) bool {
	return h[i].deadline.Before(h[j].deadline)
}

func (h taskHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}

func (h *taskHeap) Push(x interface{}) {
	n := len(*h)
	t := x.(*mockTimer)
	t.index = n
	*h = append(*h, t)
}

func (h *taskHeap) Pop() interface{} {
	old := *h
	n := len(old)
	t := old[n-1]
	t.index = removed
	*h = old[0 : n-1]
	return t
}

func (h *taskHeap) start(t *mockTimer) {
	heap.Push(h, t)
}

func (h *taskHeap) stop(t *mockTimer) {
	if !t.hasStopped() {
		heap.Remove(h, t.index)
	}
}

func (h *taskHeap) reset(t *mockTimer) {
	if !t.hasStopped() {
		heap.Fix(h, t.index)
	} else {
		heap.Push(h, t)
	}
}

func (h taskHeap) next() *mockTimer {
	if len(h) == 0 {
		return nil
	}
	return h[0]
}
