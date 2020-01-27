package clock

import (
	"context"
	"time"
)

type clockKey struct{}

// Set 把 Clock 放入 ctx 中
func Set(ctx context.Context, c Clock) context.Context {
	return context.WithValue(ctx, clockKey{}, c)
}

// Get 取出 ctx 中的 Clock
func Get(ctx context.Context) Clock {
	if c, ok := ctx.Value(clockKey{}).(Clock); ok {
		return c
	}
	return NewRealClock()
}

// After is a convenience wrapper for FromContext(ctx).After.
func After(ctx context.Context, d time.Duration) <-chan time.Time {
	return Get(ctx).After(d)
}

// AfterFunc is a convenience wrapper for FromContext(ctx).AfterFunc.
func AfterFunc(ctx context.Context, d time.Duration, f func()) *Timer {
	return Get(ctx).AfterFunc(d, f)
}

// NewTicker is a convenience wrapper for FromContext(ctx).NewTicker.
func NewTicker(ctx context.Context, d time.Duration) *Ticker {
	return Get(ctx).NewTicker(d)
}

// NewTimer is a convenience wrapper for FromContext(ctx).NewTimer.
func NewTimer(ctx context.Context, d time.Duration) *Timer {
	return Get(ctx).NewTimer(d)
}

// Now is a convenience wrapper for FromContext(ctx).Now.
func Now(ctx context.Context) time.Time {
	return Get(ctx).Now()
}

// Since is a convenience wrapper for FromContext(ctx).Since.
func Since(ctx context.Context, t time.Time) time.Duration {
	return Get(ctx).Since(t)
}

// Sleep is a convenience wrapper for FromContext(ctx).Sleep.
func Sleep(ctx context.Context, d time.Duration) {
	Get(ctx).Sleep(d)
}

// Tick is a convenience wrapper for FromContext(ctx).Tick.
func Tick(ctx context.Context, d time.Duration) <-chan time.Time {
	return Get(ctx).Tick(d)
}

// Until is a convenience wrapper for FromContext(ctx).Until.
func Until(ctx context.Context, t time.Time) time.Duration {
	return Get(ctx).Until(t)
}

// ContextWithDeadline is a convenience wrapper for FromContext(ctx).ContextWithDeadline.
func ContextWithDeadline(ctx context.Context, d time.Time) (context.Context, context.CancelFunc) {
	return Get(ctx).ContextWithDeadline(ctx, d)
}

// ContextWithTimeout is a convenience wrapper for FromContext(ctx).ContextWithTimeout.
func ContextWithTimeout(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	return Get(ctx).ContextWithTimeout(ctx, timeout)
}
