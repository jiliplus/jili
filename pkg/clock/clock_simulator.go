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
	now  time.Time
	heap *taskHeap
}

// NewSimulator 返回一个以 now 为当前时间的虚拟时钟。
func NewSimulator(now time.Time) *Simulator {
	return &Simulator{
		now:  now,
		heap: newTaskHeap(),
	}
}

// Add advances the current time by duration d and fires all expired timers if d >= 0,
// else DO NOTHING
//
// Returns the current time.
// 推荐使用 AddOrPanic 替换此方法
func (s *Simulator) Add(d time.Duration) time.Time {
	s.Lock()
	defer s.Unlock()
	if d < 0 {
		return s.now
	}
	now, _ := s.set(s.now.Add(d))
	return now
}

const (
	timeReversal = "继续执行此操作会导致 Simulator 的时间逆转"
)

// AddOrPanic advances the current time by duration d and fires all expired timers if d >= 0
// else panic
// Returns the new current time.
func (s *Simulator) AddOrPanic(d time.Duration) time.Time {
	s.Lock()
	defer s.Unlock()
	if d < 0 {
		panic(timeReversal)
	}
	now, _ := s.set(s.now.Add(d))
	return now
}

// Move advances the current time to the next available timer deadline
// Returns the new current time and the advanced duration.
func (s *Simulator) Move() (time.Time, time.Duration) {
	s.Lock()
	defer s.Unlock()
	last := s.now
	if s.heap.hasTask() {
		s.accomplishNextTask()
	}
	return s.now, s.now.Sub(last)
}

// Set advances the current time to t and fires all expired timers if s.now <= t
// else DO NOTHING
// Returns the advanced duration.
// NOTICE: 返回 0 还有可能是 t < s.now，不仅仅时 t = s.now
// 推荐使用 SetOrPanic 替代此方法
func (s *Simulator) Set(t time.Time) time.Duration {
	s.Lock()
	defer s.Unlock()
	if t.Before(s.now) {
		return 0
	}
	_, d := s.set(t)
	return d
}

// SetOrPanic advances the current time to t and fires all expired timers if s.now <= t
// else panic with time reversal
// Returns the advanced duration.
func (s *Simulator) SetOrPanic(t time.Time) time.Duration {
	s.Lock()
	defer s.Unlock()
	if t.Before(s.now) {
		panic(timeReversal)
	}
	_, d := s.set(t)
	return d
}

// Since returns the time elapsed since t.
func (s *Simulator) Since(t time.Time) time.Duration {
	s.Lock()
	defer s.Unlock()
	return s.now.Sub(t)
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
	for s.heap.hasTaskToRun(now) {
		s.accomplishNextTask()
		s.gosched()
	}
	// 如果有多个 goroutine 在并行运行 set 的话。
	// 由于 s.gosched() 的存在
	// s.now 有可能已经大于 now 了
	// 所以此处不能直接使用 s.now = now
	s.setNowTo(now)
	return s.now, s.now.Sub(last)
}

// NOTICE: 务必在临界区内运行此方法，否则会 panic。
// 目的是，其他 Goroutine 的 simulator.Now() 的操作。
// 可以在每次更新 simulator.now 后，有机会得到执行。
// 而不是得等到整个 set 函数执行完毕后，才能执行。
// 当输入参数 now 的值特别大的时候，
// simulator 的运行情况更接近 real clock
// 所以，才必须在临界区内执行
func (s *Simulator) gosched() {
	s.Unlock()
	runtime.Gosched()
	s.Lock()
}

func (s *Simulator) accomplishNextTask() {
	t := s.heap.pop()
	s.setNowTo(t.deadline)
	t = t.run()
	s.start(t)
}

func (s *Simulator) start(t *task) {
	if t == nil {
		return
	}
	if !t.deadline.After(s.now) {
		s.start(t.run())
		return
	}
	s.heap.push(t)
}

// setNowTo make m.now equal to t if m.now < t
// else do nothing
func (s *Simulator) setNowTo(t time.Time) {
	if s.now.Before(t) {
		// Simulator 的所有方法中，
		// 应该只有这一处存在 .now =
		// 需要改变 s.now 的话，就调用此方法。
		s.now = t
	}
}
