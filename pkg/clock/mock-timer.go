package clock

import "time"

// mockTimer represents a time.mockTimer.
type mockTimer struct {
	C <-chan time.Time
	*timePiece
}

// After waits for the duration to elapse and then sends the current time on
// the returned channel.
//
// A negative or zero duration fires the underlying timer immediately.
func (m *mockClock) After(d time.Duration) <-chan time.Time {
	return m.newTimer(d).C
}

// AfterFunc waits for the duration to elapse and then calls f in its own goroutine.
// It returns a Timer that can be used to cancel the call using its Stop method.
//
// A negative or zero duration fires the timer immediately.
func (m *mockClock) AfterFunc(d time.Duration, f func()) ResetStopper {
	m.Lock()
	defer m.Unlock()
	return m.newTimerFunc(m.now.Add(d), f)
}

// NewTimer creates a new Timer that will send the current time on its channel
// after at least duration d.
//
// A negative or zero duration fires the timer immediately.
func (m *mockClock) NewTimer(d time.Duration) ResetStopper {
	m.Lock()
	defer m.Unlock()
	return m.newTimer(d)
}

func (m *mockClock) newTimer(d time.Duration) *mockTimer {
	return m.newTimerFunc(m.now.Add(d), nil)
}

// Sleep pauses the current goroutine for at least the duration d.
//
// A negative or zero duration causes Sleep to return immediately.
func (m *mockClock) Sleep(d time.Duration) {
	<-m.After(d)
}

func (m *mockClock) newTimerFunc(deadline time.Time, afterFunc func()) *mockTimer {
	t := &mockTimer{
		timePiece: newTimePiece(m, deadline),
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
		m.start(t.timePiece)
	}
	return t
}

// Stop prevents the Timer from firing.
// It returns true if the call stops the timer, false if the timer has already
// expired or been stopped.
func (t *mockTimer) Stop() bool {
	t.mock.Lock()
	defer t.mock.Unlock()
	wasActive := !t.timePiece.hasStopped()
	t.mock.stop(t.timePiece)
	return wasActive
}

// Reset changes the timer to expire after duration d.
// It returns true if the timer had been active, false if the timer had
// expired or been stopped.
//
// A negative or zero duration fires the timer immediately.
func (t *mockTimer) Reset(d time.Duration) bool {
	t.mock.Lock()
	defer t.mock.Unlock()
	wasActive := !t.timePiece.hasStopped()
	t.deadline = t.mock.now.Add(d)
	if !t.deadline.After(t.mock.now) {
		t.fire()
		t.mock.stop(t.timePiece)
	} else {
		t.mock.reset(t.timePiece)
	}
	return wasActive
}
