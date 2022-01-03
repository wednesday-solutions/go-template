package testutls

import (
	"context"
	"time"
)

type MockContext struct{}

func (ctx MockContext) Deadline() (deadline time.Time, ok bool) {
	return deadline, ok
}

func (ctx MockContext) Done() <-chan struct{} {
	ch := make(chan struct{})
	close(ch)
	return ch
}

func (ctx MockContext) Err() error {
	return context.DeadlineExceeded
}

func (ctx MockContext) Value(key interface{}) interface{} {
	return nil
}
