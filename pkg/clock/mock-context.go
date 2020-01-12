package clock

import (
	"context"
	"time"
)

type mockContext struct {
	context.Context
	deadline time.Time
	done     chan struct{}
	err      error
}

func (ctx *mockContext) Deadline() (time.Time, bool) {
	// TODO: 为什么一定是 true
	return ctx.deadline, true
}

func (ctx *mockContext) Done() <-chan struct{} {
	return ctx.done
}

func (ctx *mockContext) Err() error {
	select {
	case <-ctx.done:
		return ctx.err
	default:
		return nil
	}
}
