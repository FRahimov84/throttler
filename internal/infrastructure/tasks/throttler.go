package tasks

import (
	"context"
	"fmt"
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

func New(n int, k int, x int) (*ThrottlerTask, error) {
	if n < 1 || k < 1 {
		return nil, fmt.Errorf("bad task options n: %d, k: %d. Options should be greater than zero", n, k)
	}
	return &ThrottlerTask{n: n, k: k, x: x}, nil
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
