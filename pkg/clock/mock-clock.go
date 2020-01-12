package clock

import (
	"context"
	"sync"
	"time"
)

// *mockClock 实现了 Clock 接口和 Updater 接口
// Clock 接口提供了以 mock clock 当前时间为基准的时间运算。
// Updater 接口会修改 mock clock 的当前时间。
// Clock 和 Updater 方法都是并发安全的。
type mockClock struct {
	sync.RWMutex
	now time.Time
	*pieceHeap
}

// NewMockClock returns a new Mock with current time set to now.
func NewMockClock(now time.Time) UpdatableClock {
	return &mockClock{
		now:       now,
		pieceHeap: newPieceHeap(),
	}
}

// Add advances the current time by duration d and fires all expired timers.
//
// Returns the new current time.
// To increase predictability and speed, Tickers are ticked only once per call.
// TODO: ? To increase predictability and speed, Tickers are ticked only once per call.
func (m *mockClock) Add(d time.Duration) time.Time {
	m.Lock()
	defer m.Unlock()
	now, _ := m.set(m.now.Add(d))
	return now
}

// Move advances the current time to the next available timer deadline
// and fires all expired timers.
//
// Returns the new current time and the advanced duration.
func (m *mockClock) Move() (time.Time, time.Duration) {
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
func (m *mockClock) Set(t time.Time) time.Duration {
	m.Lock()
	defer m.Unlock()
	if t.Sub(m.now) < 0 {
		panic("mockClock.Set: t < m.now")
	}
	_, d := m.set(t)
	return d
}

func (m *mockClock) set(now time.Time) (time.Time, time.Duration) {
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

// Now returns the current mocked time.
func (m *mockClock) Now() time.Time {
	m.Lock()
	defer m.Unlock()
	return m.now
}

// Since returns the time elapsed since t.
func (m *mockClock) Since(t time.Time) time.Duration {
	m.Lock()
	defer m.Unlock()
	return m.now.Sub(t)
}

// Until returns the duration until t.
func (m *mockClock) Until(t time.Time) time.Duration {
	m.Lock()
	defer m.Unlock()
	return t.Sub(m.now)
}

// ContextWithDeadline implements Clock.
func (m *mockClock) ContextWithDeadline(parent context.Context, d time.Time) (context.Context, context.CancelFunc) {
	m.Lock()
	defer m.Unlock()
	return m.contextWithDeadline(parent, d)
}

// ContextWithTimeout implements Clock.
func (m *mockClock) ContextWithTimeout(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	m.Lock()
	defer m.Unlock()
	return m.contextWithDeadline(parent, m.now.Add(timeout))
}

func (m *mockClock) contextWithDeadline(parent context.Context, deadline time.Time) (context.Context, context.CancelFunc) {
	cancelCtx, cancel := context.WithCancel(SetClock(parent, m))
	if pd, ok := parent.Deadline(); ok && !pd.After(deadline) {
		return cancelCtx, cancel
	}
	ctx := &mockContext{
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
