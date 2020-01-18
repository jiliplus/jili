package clock

import (
	"container/heap"
	"time"
)

type taskManager interface {
	hasTaskToRun(now time.Time) bool
	pop() *task
	push(t *task)
	remove(t *task)
}

type task struct {
	deadline time.Time
	// fire 是用来区分 task 的上级时 timer 还是 tick
	// 当是 tick 的时候，返回 tick 周期，好计算下一个 tick 的时间。
	// 但是，现在的实现方式，有一个大问题，
	// 当 mock clock 的调整时间，大于 tick 的周期的时候，
	// 只会触发 tick 一次，这明显与实际不符
	// TODO: 删除此处内容
	fire func() time.Duration
	// TODO: 删除此处内容
	// task 不必包含 mock 的内容。
	// 由 mock 驱动 heap 的时候，丢弃 !isActive 的 task 就好了。
	mock *Mock
	// tick 或 timer
	// 在 Stop 之前，isStopped = true
	// 在 Stop 之后，isStopped = false
	// isStopped bool
	// 用于替代 fire，
	runTask func(t *task) *task
	index   int
}

const removed = -1

func newTask2(deadline time.Time, runTask func(t *task) *task) *task {
	return &task{
		deadline: deadline,
		runTask:  runTask,
		index:    removed,
	}
}

func newTask(m *Mock, d time.Time) *task {
	return &task{
		deadline: d,
		mock:     m,
		index:    removed,
	}
}

func (t *task) run() *task {
	return t.runTask(t)
}

func (t task) hasStopped() bool {
	return t.index == removed
}

// TODO: 删除此处内容
func (h *taskHeap) start(t *task) {
	heap.Push(h, t)
}

// TODO: 删除此处内容
func (h *taskHeap) stop(t *task) {
	if !t.hasStopped() {
		heap.Remove(h, t.index)
	}
}

// TODO: 删除此处内容
func (h *taskHeap) reset(t *task) {
	if !t.hasStopped() {
		heap.Fix(h, t.index)
	} else {
		heap.Push(h, t)
	}
}

// TODO: 删除此处内容
func (h taskHeap) next() *task {
	if len(h) == 0 {
		return nil
	}
	return h[0]
}

type taskHeap []*task

func newTaskHeap() *taskHeap {
	t := make(taskHeap, 0, 1024)
	return &t
}

// *taskHeap 实现了 taskOrder 接口
func (h *taskHeap) push(t *task) {
	heap.Push(h, t)
}

func (h *taskHeap) pop() (t *task) {
	t, _ = heap.Pop(h).(*task)
	return
}

// TODO: 删除此处内容
func (h taskHeap) hasTaskToRun(now time.Time) bool {
	return len(h) != 0 && !now.Before(h[0].deadline)
}

func (h taskHeap) hasExpiredTask(now time.Time) bool {
	return len(h) != 0 && !now.Before(h[0].deadline)
}

func (h taskHeap) hasTask() bool {
	return len(h) != 0
}

func (h *taskHeap) remove(t *task) {
	if !t.hasStopped() {
		heap.Remove(h, t.index)
	}
}

// TODO: 删除此处内容
func (h taskHeap) hasNext() bool {
	return len(h) != 0
}

// *taskHeap 实现了 heap.Interface
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
	t := x.(*task)
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
