package clock

import (
	"runtime"
	"sync"
	"time"
)

// Simulator 实现了 Clock 接口，并提供了 .Add*，.Set* 和 .Move 方法驱动时钟运行
//
// 为了尽可能真实地模拟时间的流逝，Simulator.now 只会不断变大，不会出现逆转情况。
//
// RWMutex 锁住 Simulator 时，其他 goroutine 的 Simulator 方法会被阻塞。
// Simulator 的运行也不适均匀的，有可能下一个时刻就是很久以后。
// 这是与 time 标准库的主要差异，使用 Simulator 时，请特别注意。
//
type Simulator struct {
	sync.RWMutex
	now         time.Time
	taskManager taskManager
}

// NewSimulator 返回一个以 now 为当前时间的虚拟时钟。
func NewSimulator(now time.Time) *Mock {
	return newSimulator(now, newTaskHeap())
}

// newSimulator 返回一个以 now 为当前时间的虚拟时钟。
func newSimulator(now time.Time, tm taskManager) *Mock {
	return &Mock{
		now:         now,
		taskManager: tm,
	}
}

// set 是 Simulator 的核心逻辑，
// 把 now 时间点之前需要完成的任务，由早到晚依次触发。
// 一边触发，一边把 t.deadline 设置为 Simulator.now
// 每个 for 循环的结尾，利用 s.gosched 离开临界区，
// 给其他阻塞的操作执行的机会。
// 全部触发完毕后，让 Simulator.now = now
//
// 当多个 set 同时调用时，会交替运行，并发安全。
// 只是较小的输入参数 now，可能无法被赋值到 Simulator.now
func (s *Simulator) set(now time.Time) (time.Time, time.Duration) {
	last := s.now
	for s.taskManager.hasTaskToRun(now) {
		s.accomplishNextTask()
		s.gosched()
	}
	s.setNowTo(now)
	return s.now, s.now.Sub(last)
}

// NOTICE: 务必在临界区内运行此方法，否则会 panic。
// 因为在修改了 mock.now 后，解锁又上锁的操作。
// 目的是，其他 Goroutine 的 mock.Now() 的操作。
// 可以在每次更新 mock.now 后，有机会得到执行。
// 而不是得等到整个 set 函数执行完毕后，才能执行。
// 当输入参数 now 的值特别大的时候，
// mock clock 的运行情况更接近 real clock
// 所以，才必须在临界区内执行
func (s *Simulator) gosched() {
	s.Unlock()
	runtime.Gosched()
	// TODO: 可是，如果这样的话，等 Lock 成功的时候，都不知道时什么时候了。
	// 比如说，有两个 set2 在运行，情况会时怎么样的呢？
	s.Lock()
}

func (s *Simulator) accomplishNextTask() {
	t := s.taskManager.pop()
	s.setNowTo(t.deadline)
	t = t.run()
	s.start(t)
}

func (s *Simulator) start(t *task) {
	if t == nil {
		return
	}
	if !t.deadline.After(s.now) {
		t.run()
	}
	s.taskManager.push(t)
}

// setNowTo make m.now equal to t if m.now < t
// else do nothing
func (s *Simulator) setNowTo(t time.Time) {
	if s.now.Before(t) {
		// Simulator 的所有方法中，
		// 应该只有这一处存在 .now =
		s.now = t
	}
}
