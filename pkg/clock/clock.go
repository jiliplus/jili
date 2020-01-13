package clock

import (
	"context"
	"time"
)

type realClock struct{}

// New 返回标准库中真实时间的时钟。
// 并实现了 Clock 接口
func New() Clock {
	return realClock{}
}

func (realClock) After(d time.Duration) <-chan time.Time {
	return time.After(d)
}

func (realClock) AfterFunc(d time.Duration, f func()) *Timer {
	return &Timer{timer: time.AfterFunc(d, f)}
}

func (realClock) NewTicker(d time.Duration) *Ticker {
	t := time.NewTicker(d)
	return &Ticker{
		C:      t.C,
		ticker: t,
	}
}

func (realClock) NewTimer(d time.Duration) *Timer {
	t := time.NewTimer(d)
	return &Timer{
		C:     t.C,
		timer: t,
	}
}

func (realClock) Now() time.Time {
	return time.Now()
}

func (realClock) Since(t time.Time) time.Duration {
	return time.Since(t)
}

func (realClock) Sleep(d time.Duration) {
	time.Sleep(d)
}

func (realClock) Tick(d time.Duration) <-chan time.Time {
	// Using time.Tick would trigger a vet tool warning.
	// if d <= 0 {
	// return nil
	// }
	// TODO: 把以下内容放入 mockTicker.Tick 中
	// panic(errors.New("non-positive interval for NewTicker"))
	return time.NewTicker(d).C
}

func (realClock) Until(t time.Time) time.Duration {
	return time.Until(t)
}

func (realClock) DeadlineContext(parent context.Context, d time.Time) (context.Context, context.CancelFunc) {
	return context.WithDeadline(parent, d)
}

func (realClock) ContextWithTimeout(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(parent, timeout)
}
