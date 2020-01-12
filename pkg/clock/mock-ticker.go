package clock

import (
	"errors"
	"time"
)

// mockTicker represents a time.mockTicker.
type mockTicker struct {
	C    <-chan time.Time
	fire func() time.Duration
	mock *mockClock
	*timePieceOld
}

// NewTicker returns a new Ticker containing a channel that will send the
// current time with a period specified by the duration d.
func (m *mockClock) NewTicker(d time.Duration) Stopper {
	m.Lock()
	defer m.Unlock()
	if d <= 0 {
		panic(errors.New("non-positive interval for NewTicker"))
	}
	return m.newTicker(d)
}

// Tick is a convenience wrapper for NewTicker providing access to the ticking
// channel only.
func (m *mockClock) Tick(d time.Duration) <-chan time.Time {
	m.Lock()
	defer m.Unlock()
	if d <= 0 {
		return nil
	}
	return m.newTicker(d).C
}

func (m *mockClock) newTicker(d time.Duration) *mockTicker {
	c := make(chan time.Time, 1)
	t := &mockTicker{
		C:            c,
		timePieceOld: newTimePiece(m, m.now.Add(d)),
	}
	t.fire = func() time.Duration {
		select {
		case c <- m.now:
		default:
		}
		return d
	}
	m.start(t.timePieceOld)
	return t
}

// Stop turns off a ticker. After Stop, no more ticks will be sent.
func (t *mockTicker) Stop() {
	t.mock.Lock()
	defer t.mock.Unlock()
	t.mock.stop(t.timePieceOld)
}
