package clock

import "time"

// Timer2 替代 time.Timer.
// TODO: 修改名称为 Timer
type Timer2 struct {
	C <-chan time.Time
	// 当 timer != nil 的时候, Timer 代表了 real clock
	timer *time.Timer
	*task

	// Stop prevents the Timer from firing.
	// It returns true if the call stops the timer, false if the timer has already expired or been stopped.
	Stop func() bool
	// Reset changes the timer to expire after duration d.
	// It returns true if the timer had been active, false if the timer had expired or been stopped.
	//
	// A negative or zero duration fires the timer immediately.
	Reset func(d time.Duration) bool
}

// Sleep pauses the current goroutine for at least the duration d.
//
// A negative or zero duration causes Sleep to return immediately.
func (s *Simulator) Sleep(d time.Duration) {
	<-s.After(d)
}

// After waits for the duration to elapse and then sends the current time on
// the returned channel.
//
// A negative or zero duration fires the underlying timer immediately.
func (s *Simulator) After(d time.Duration) <-chan time.Time {
	return s.NewTimer(d).C
}

// AfterFunc waits for the duration to elapse and then calls f in its own goroutine.
// It returns a Timer that can be used to cancel the call using its Stop method.
//
// A negative or zero duration fires the timer immediately.
func (s *Simulator) AfterFunc(d time.Duration, f func()) *Timer2 {
	s.Lock()
	defer s.Unlock()
	return s.newTimerFunc(s.now.Add(d), f)
}

// NewTimer creates a new Timer that will send the current time on its channel
// after at least duration d.
//
// A negative or zero duration fires the timer immediately.
func (s *Simulator) NewTimer(d time.Duration) *Timer2 {
	s.Lock()
	defer s.Unlock()
	return s.newTimerFunc(s.now.Add(d), nil)
}

//
func (s *Simulator) newTimerFunc(deadline time.Time, afterFunc func()) *Timer2 {
	c := make(chan time.Time, 1)
	runTask := func(t *task) *task {
		if afterFunc != nil {
			go afterFunc()
		} else {
			// 因为 time.Tick 的处理逻辑也是这样的
			// 有人收就发过去, 每人接收就丢弃.
			// NOTICE: AfterFunc 创建的 *Timer 不会发送 current time
			select {
			case c <- s.now:
			default:
			}
		}
		return nil
	}
	t := &Timer2{
		C:    c,
		task: newTask2(deadline, runTask),
	}
	s.accept(t.task)
	t.Stop = func() bool {
		if t.timer != nil {
			return t.timer.Stop()
		}
		s.Lock()
		defer s.Unlock()
		isActive := !t.task.hasStopped()
		s.heap.remove(t.task)
		return isActive
	}
	t.Reset = func(d time.Duration) bool {
		if t.timer != nil {
			return t.timer.Reset(d)
		}
		s.Lock()
		defer s.Unlock()
		isActive := !t.hasStopped()
		s.heap.remove(t.task)
		t.deadline = s.now.Add(d)
		s.accept(t.task)
		return isActive
	}
	return t
}
