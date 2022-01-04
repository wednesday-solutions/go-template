package testutls

import (
	"context"
	"time"
)

type MockCtx struct{}

func (ctx MockCtx) Deadline() (deadline time.Time, ok bool) {
	return deadline, ok
}

func (ctx MockCtx) Done() <-chan struct{} {
	ch := make(chan struct{})
	close(ch)
	return ch
}

func (ctx MockCtx) Err() error {
	return context.DeadlineExceeded
}

func (ctx MockCtx) Value(key interface{}) interface{} {
	return nil
}
