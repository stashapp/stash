package utils

import (
	"context"
	"time"
)

type ValueOnlyContext struct {
	context.Context
}

func (ValueOnlyContext) Deadline() (deadline time.Time, ok bool) {
	return
}

func (ValueOnlyContext) Done() <-chan struct{} {
	return nil
}

func (ValueOnlyContext) Err() error {
	return nil
}
