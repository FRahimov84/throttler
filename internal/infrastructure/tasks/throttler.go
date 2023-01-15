package tasks

import (
	"context"
	"time"
)

type ThrottlerTask struct {
	n int
	k int
	x int
}

type Caller interface {
	Call(ctx context.Context)
}

func New(n int, k int, x int) *ThrottlerTask {
	return &ThrottlerTask{n: n, k: k, x: x}
}

func (t *ThrottlerTask) Do(ctx context.Context, caller Caller) {
	rate := time.Second * time.Duration(t.k/t.n)
	burstLimit := t.n
	tick := time.NewTicker(rate)
	defer tick.Stop()
	throttle := make(chan struct{}, burstLimit)
	go func() {
		for _ = range tick.C {
			throttle <- struct{}{}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case <-throttle:
			go caller.Call(ctx)
		}
	}
}
