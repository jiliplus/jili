package clock

import (
	"context"
	"runtime"
	"sync"
	"time"
)

// Mock 实现了 Clock 接口，并提供了 .Add*，.Set* 和 .Move 方法驱动时钟运行
//
// 为了尽可能真实地模拟时间的流逝，Mock.now 只会不断变大，不会出现逆转情况。
//
// RWMutex 锁住 Mock 时，其他 goroutine 的 Mock 方法会被阻塞。
// Mock 的运行也不适均匀的，有可能下一个时刻就是很久以后。
// 这是与 time 标准库的主要差异，使用 Mock 时，请特别注意。
//
// TODO: 删除此处内容
type Mock struct {
	sync.RWMutex
	now time.Time
	// TODO: 删除此处内容
	mockTimers
	taskManager taskManager
}

// NewMockClock 返回一个以 now 为当前时间的虚拟时钟。
// TODO: 删除此处内容
func NewMockClock(now time.Time) *Mock {
	return &Mock{
		now:         now,
		mockTimers:  &taskHeap{},
		taskManager: newTaskHeap(),
	}
}

// TODO: 删除此处内容
type mockTimers interface {
	start(t *task)
	stop(t *task)
	reset(t *task)
	next() *task
}

// Add advances the current time by duration d and fires all expired timers.
//
// Returns the new current time.
// To increase predictability and speed, Tickers are ticked only once per call.
// TODO: 删除此处内容
func (m *Mock) Add(d time.Duration) time.Time {
	m.Lock()
	defer m.Unlock()
	now, _ := m.set(m.now.Add(d))
	return now
}

// AddNext advances the current time to the next available timer deadline
// and fires all expired timers.
//
// Returns the new current time and the advanced duration.
// TODO: 删除此处内容
func (m *Mock) AddNext() (time.Time, time.Duration) {
	m.Lock()
	defer m.Unlock()
	t := m.next()
	if t == nil {
		return m.now, 0
	}
	return m.set(t.deadline)
}

// Set advances the current time to t and fires all expired timers.
//
// Returns the advanced duration.
// To increase predictability and speed, Tickers are ticked only once per call.
// TODO: 删除此处内容
func (m *Mock) Set(t time.Time) time.Duration {
	m.Lock()
	defer m.Unlock()
	_, d := m.set(t)
	return d
}

func (m *Mock) set(now time.Time) (time.Time, time.Duration) {
	cur := m.now
	for {
		t := m.next()
		if t == nil || t.deadline.After(now) {
			m.now = now
			return m.now, m.now.Sub(cur)
		}
		m.now = t.deadline
		if d := t.fire(); d == 0 {
			// Timers are always stopped.
			m.stop(t)
		} else {
			// Ticker's next deadline is set to the first tick after the new now.
			dd := (now.Sub(m.now)/d + 1) * d
			t.deadline = m.now.Add(dd)
			m.reset(t)
		}
	}
}

// NOTICE: 务必在临界区内运行此方法，否则会 panic。
// 因为在修改了 mock.now 后，解锁又上锁的操作。
// 目的是，其他 Goroutine 的 mock.Now() 的操作。
// 可以在每次更新 mock.now 后，有机会得到执行。
// 而不是得等到整个 set 函数执行完毕后，才能执行。
// 当输入参数 now 的值特别大的时候，
// mock clock 的运行情况更接近 real clock
// 所以，才必须在临界区内执行
func (m *Mock) gosched() {
	m.Unlock()
	runtime.Gosched()
	// TODO: 可是，如果这样的话，等 Lock 成功的时候，都不知道时什么时候了。
	// 比如说，有两个 set2 在运行，情况会时怎么样的呢？
	m.Lock()
}

// Now returns the current mocked time.
func (m *Mock) Now() time.Time {
	m.Lock()
	defer m.Unlock()
	return m.now
}

// Since returns the time elapsed since t.
// TODO: 删除此处内容
func (m *Mock) Since(t time.Time) time.Duration {
	m.Lock()
	defer m.Unlock()
	return m.now.Sub(t)
}

// Until returns the duration until t.
// TODO: 删除此处内容
func (m *Mock) Until(t time.Time) time.Duration {
	m.Lock()
	defer m.Unlock()
	return t.Sub(m.now)
}

// ContextWithDeadline implements Clock.
func (m *Mock) ContextWithDeadline(parent context.Context, d time.Time) (context.Context, context.CancelFunc) {
	m.Lock()
	defer m.Unlock()
	return m.contextWithDeadline(parent, d)
}

// ContextWithTimeout implements Clock.
func (m *Mock) ContextWithTimeout(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	m.Lock()
	defer m.Unlock()
	return m.contextWithDeadline(parent, m.now.Add(timeout))
}

func (m *Mock) contextWithDeadline(parent context.Context, deadline time.Time) (context.Context, context.CancelFunc) {
	cancelCtx, cancel := context.WithCancel(Set(parent, m))
	if pd, ok := parent.Deadline(); ok && !pd.After(deadline) {
		return cancelCtx, cancel
	}
	// TODO: 把以下代码放入 newMockContext
	ctx := &mockCtx{
		Context:  cancelCtx,
		done:     make(chan struct{}),
		deadline: deadline,
	}
	t := m.newTimerFunc(deadline, nil)
	go func() {
		select {
		case <-t.C:
			ctx.err = context.DeadlineExceeded
		case <-cancelCtx.Done():
			ctx.err = cancelCtx.Err()
			defer t.Stop()
		}
		close(ctx.done)
	}()
	return ctx, cancel
}

// TODO: 删除此处内容
func (m *Mock) start(t *task) {
	if t == nil {
		return
	}
	if !t.deadline.After(m.now) {
		t.run()
	}
	m.taskManager.push(t)
}
