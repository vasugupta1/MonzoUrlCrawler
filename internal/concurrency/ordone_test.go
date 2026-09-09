package concurrency

import (
	"context"
	"testing"
	"time"
)

func TestOrDone(t *testing.T) {
	t.Run("forwards all values from channel", func(t *testing.T) {
		ch := make(chan int, 3)
		ch <- 1
		ch <- 2
		ch <- 3
		close(ch)

		var got []int
		for v := range OrDone(context.Background(), ch) {
			got = append(got, v)
		}
		if len(got) != 3 {
			t.Errorf("got %v, want [1 2 3]", got)
		}
	})

	t.Run("closes output when input channel closes", func(t *testing.T) {
		ch := make(chan int)
		close(ch)

		out := OrDone(context.Background(), ch)
		_, ok := <-out
		if ok {
			t.Error("expected output channel to be closed")
		}
	})

	t.Run("closes output when context is cancelled", func(t *testing.T) {
		ch := make(chan int) // never sends anything
		ctx, cancel := context.WithCancel(context.Background())
		out := OrDone(ctx, ch)

		cancel()

		select {
		case _, ok := <-out:
			if ok {
				t.Error("expected output channel to be closed")
			}
		case <-time.After(time.Second):
			t.Fatal("timed out waiting for output channel to close")
		}
	})

	t.Run("stops forwarding mid-stream on cancel", func(t *testing.T) {
		ch := make(chan int)
		ctx, cancel := context.WithCancel(context.Background())
		out := OrDone(ctx, ch)

		// send one value, then cancel
		ch <- 42
		v := <-out
		if v != 42 {
			t.Errorf("got %d, want 42", v)
		}

		cancel()

		select {
		case _, ok := <-out:
			if ok {
				t.Error("expected output channel to be closed after cancel")
			}
		case <-time.After(time.Second):
			t.Fatal("timed out waiting for output channel to close")
		}
	})
}
