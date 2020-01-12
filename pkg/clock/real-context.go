package clock

import (
	"context"
	"time"
)

var clockKey = struct{}{}

// SetClock returns a copy of parent in which the Clock is associated with.
func SetClock(parent context.Context, c Clock) context.Context {
	return context.WithValue(parent, clockKey, c)
}

// GetClock returns the Clock associated with the context, or Realtime().
func GetClock(ctx context.Context) Clock {
	if c, ok := ctx.Value(clockKey).(Clock); ok {
		return c
	}
	return New()
}

// After is a convenience wrapper for FromContext(ctx).After.
func After(ctx context.Context, d time.Duration) <-chan time.Time {
	return GetClock(ctx).After(d)
}

// AfterFunc is a convenience wrapper for FromContext(ctx).AfterFunc.
func AfterFunc(ctx context.Context, d time.Duration, f func()) ResetStopper {
	return GetClock(ctx).AfterFunc(d, f)
}

// NewTicker is a convenience wrapper for FromContext(ctx).NewTicker.
func NewTicker(ctx context.Context, d time.Duration) Stopper {
	return GetClock(ctx).NewTicker(d)
}

// NewTimer is a convenience wrapper for FromContext(ctx).NewTimer.
func NewTimer(ctx context.Context, d time.Duration) ResetStopper {
	return GetClock(ctx).NewTimer(d)
}

// Now is a convenience wrapper for FromContext(ctx).Now.
func Now(ctx context.Context) time.Time {
	return GetClock(ctx).Now()
}

// Since is a convenience wrapper for FromContext(ctx).Since.
func Since(ctx context.Context, t time.Time) time.Duration {
	return GetClock(ctx).Since(t)
}

// Sleep is a convenience wrapper for FromContext(ctx).Sleep.
func Sleep(ctx context.Context, d time.Duration) {
	GetClock(ctx).Sleep(d)
}

// Tick is a convenience wrapper for FromContext(ctx).Tick.
func Tick(ctx context.Context, d time.Duration) <-chan time.Time {
	return GetClock(ctx).Tick(d)
}

// Until is a convenience wrapper for FromContext(ctx).Until.
func Until(ctx context.Context, t time.Time) time.Duration {
	return GetClock(ctx).Until(t)
}

// ContextWithDeadline is a convenience wrapper for FromContext(ctx).DeadlineContext.
func ContextWithDeadline(ctx context.Context, d time.Time) (context.Context, context.CancelFunc) {
	return GetClock(ctx).ContextWithDeadline(ctx, d)
}

// ContextWithTimeout is a convenience wrapper for FromContext(ctx).TimeoutContext.
func ContextWithTimeout(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	return GetClock(ctx).ContextWithTimeout(ctx, timeout)
}
