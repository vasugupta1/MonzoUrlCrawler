package concurrency

import "context"

type RateLimiter[T any] struct {
	ch chan struct{}
}

func NewRateLimiter[T any](limit int) *RateLimiter[T] {
	return &RateLimiter[T]{
		ch: make(chan struct{}, limit),
	}
}

func (rt *RateLimiter[T]) Process(ctx context.Context, fun func() (T, error)) (T, error) {
	select {
	case <-ctx.Done():
		var zero T
		return zero, ctx.Err()
	case rt.ch <- struct{}{}:
		res, err := fun()
		<-rt.ch
		return res, err
	}
}
