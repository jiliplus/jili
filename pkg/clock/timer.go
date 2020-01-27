package clock

import "time"

// Timer represents a time.Timer.
type Timer struct {
	C <-chan time.Time
	// TODO: 修改 Stop2 为 Stop
	// https://golang.org/pkg/time/#Timer.Reset
	// TODO: 好好阅读标准库中,关于 Reset 和 Stop 的处理逻辑
	// Stop prevents the Timer from firing.
	// It returns true if the call stops the timer, false if the timer has already
	// expired or been stopped.
	Stop2 func() bool
	// Reset changes the timer to expire after duration d.
	// It returns true if the timer had been active, false if the timer had
	// expired or been stopped.
	//
	// A negative or zero duration fires the timer immediately.
	Reset2 func(d time.Duration) bool
	// 当 timer != nil 的时候, Timer 代表了 real clock
	timer *time.Timer
	// TODO: 删除此处内容,使用 Stop2 以后,可以不用保留 task 属性了
	*task
}

// After waits for the duration to elapse and then sends the current time on
// the returned channel.
//
// A negative or zero duration fires the underlying timer immediately.
func (m *Mock) After(d time.Duration) <-chan time.Time {
	return m.NewTimer(d).C
}

// AfterFunc waits for the duration to elapse and then calls f in its own goroutine.
// It returns a Timer that can be used to cancel the call using its Stop method.
//
// A negative or zero duration fires the timer immediately.
func (m *Mock) AfterFunc(d time.Duration, f func()) *Timer {
	m.Lock()
	defer m.Unlock()
	return m.newTimerFunc(m.now.Add(d), f)
}

// NewTimer creates a new Timer that will send the current time on its channel
// after at least duration d.
//
// A negative or zero duration fires the timer immediately.
func (m *Mock) NewTimer(d time.Duration) *Timer {
	m.Lock()
	defer m.Unlock()
	return m.newTimerFunc(m.now.Add(d), nil)
}

// Sleep pauses the current goroutine for at least the duration d.
//
// A negative or zero duration causes Sleep to return immediately.
func (m *Mock) Sleep(d time.Duration) {
	<-m.After(d)
}

func (m *Mock) newTimerFunc(deadline time.Time, afterFunc func()) *Timer {
	t := &Timer{
		task: newTask(m, deadline),
	}
	if afterFunc != nil {
		t.fire = func() time.Duration {
			go afterFunc()
			return 0
		}
	} else {
		c := make(chan time.Time, 1)
		t.C = c
		t.fire = func() time.Duration {
			select {
			case c <- m.now:
			default:
			}
			return 0
		}
	}
	if !t.deadline.After(m.now) {
		t.fire()
	} else {
		m.start(t.task)
	}
	return t
}

// Stop prevents the Timer from firing.
// It returns true if the call stops the timer, false if the timer has already
// expired or been stopped.
// TODO: 删除此处内容
func (t *Timer) Stop() bool {
	if t.timer != nil {
		return t.timer.Stop()
	}
	t.mock.Lock()
	defer t.mock.Unlock()
	wasActive := !t.task.hasStopped()
	t.mock.stop(t.task)
	return wasActive
}

// Reset changes the timer to expire after duration d.
// It returns true if the timer had been active, false if the timer had
// expired or been stopped.
//
// A negative or zero duration fires the timer immediately.
// TODO: 删除此处内容
func (t *Timer) Reset(d time.Duration) bool {
	if t.timer != nil {
		return t.timer.Reset(d)
	}
	t.mock.Lock()
	defer t.mock.Unlock()
	wasActive := !t.task.hasStopped()
	t.deadline = t.mock.now.Add(d)
	if !t.deadline.After(t.mock.now) {
		t.fire()
		t.mock.stop(t.task)
	} else {
		t.mock.reset(t.task)
	}
	return wasActive
}
