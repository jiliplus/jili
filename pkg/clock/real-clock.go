package clock

import (
	"context"
	"time"
)

type realClock struct{}

// New 返回标准库中真实时间的时钟。
func New() Clock {
	return realClock{}
}

func (realClock) After(d time.Duration) <-chan time.Time {
	return time.After(d)
}

func (realClock) AfterFunc(d time.Duration, f func()) ResetStopper {
	return time.AfterFunc(d, f)
}

func (realClock) NewTicker(d time.Duration) Stopper {
	return time.NewTicker(d)
}

func (realClock) NewTimer(d time.Duration) ResetStopper {
	return time.NewTimer(d)
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
	return time.NewTicker(d).C
}

func (realClock) Until(t time.Time) time.Duration {
	return time.Until(t)
}

func (realClock) ContextWithDeadline(parent context.Context, d time.Time) (context.Context, context.CancelFunc) {
	return context.WithDeadline(parent, d)
}

func (realClock) ContextWithTimeout(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(parent, timeout)
}
