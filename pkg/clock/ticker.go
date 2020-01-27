package clock

import (
	"errors"
	"time"
)

// Ticker represents a time.Ticker.
type Ticker struct {
	C <-chan time.Time
	// TODO: 修改 Stop2 为 Stop
	Stop2 func()
	// 当 ticker != nil 的时候, Ticker 代表了 real clock
	ticker *time.Ticker
	// TODO: 删除此处内容,使用 Stop2 以后,可以不用保留 task 属性了
	*task
}

// NewTicker returns a new Ticker containing a channel that will send the
// current time with a period specified by the duration d.
func (m *Mock) NewTicker(d time.Duration) *Ticker {
	m.Lock()
	defer m.Unlock()
	if d <= 0 {
		panic(errors.New("non-positive interval for NewTicker"))
	}
	return m.newTicker(d)
}

// Tick is a convenience wrapper for NewTicker providing access to the ticking
// channel only.
func (m *Mock) Tick(d time.Duration) <-chan time.Time {
	m.Lock()
	defer m.Unlock()
	if d <= 0 {
		return nil
	}
	return m.newTicker(d).C
}

// TODO: 删除此处内容
func (m *Mock) newTicker(d time.Duration) *Ticker {
	c := make(chan time.Time, 1)
	t := &Ticker{
		C:    c,
		task: newTask(m, m.now.Add(d)),
	}
	t.fire = func() time.Duration {
		select {
		case c <- m.now:
		default:
		}
		return d
	}
	m.start(t.task)
	return t
}

func (m *Mock) newTicker2(d time.Duration) *Ticker {
	c := make(chan time.Time, 1)
	run := func(t *task) *task {
		// 因为 time.Tick 的处理逻辑也是这样的
		// 有人收就发过去, 每人接收就丢弃.
		select {
		case c <- m.now:
		default:
		}
		t.deadline = t.deadline.Add(d)
		return t
	}
	t := &Ticker{
		C:    c,
		task: newTask2(m.now.Add(d), run),
	}
	t.Stop2 = func() {
		if t.ticker != nil {
			t.ticker.Stop()
			return
		}
		m.Lock()
		m.taskManager.remove(t.task)
		m.Unlock()
	}
	m.start(t.task)
	return t
}

// Stop turns off a ticker. After Stop, no more ticks will be sent.
// TODO: 删除此处内容,把 stop2 修改成 stop
func (t *Ticker) Stop() {
	if t.ticker != nil {
		t.ticker.Stop()
		return
	}
	t.mock.Lock()
	defer t.mock.Unlock()
	t.mock.stop(t.task)
}
