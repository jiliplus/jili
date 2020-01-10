package fakecontext

import (
	"context"
	"time"
)

func init() {
	// 确保 *fake 实现了 context.Context 接口
	var _ context.Context = new(fake)
}

// fake 实现了 context.Context 接口
type fake struct {
}

func (f *fake) Deadline() (deadline time.Time, ok bool) {
	return time.Now(), false
}

func (f *fake) Done() <-chan struct{} {
	return nil
}

func (f *fake) Err() error {
	return nil
}

func (f *fake) Value(key interface{}) interface{} {
	return nil
}
