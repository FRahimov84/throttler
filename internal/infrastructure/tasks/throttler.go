package tasks

import (
	"context"
	"github.com/FRahimov84/throttler/internal/usecase"
	"time"
)

type ThrottlerTask struct {
	n int
	k int
	x int
}

func New(n int, k int, x int) *ThrottlerTask {
	return &ThrottlerTask{n: n, k: k, x: x}
}

func (t *ThrottlerTask) Do(ctx context.Context, caller usecase.Caller) {
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
