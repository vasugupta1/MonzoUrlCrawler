package concurrency

import (
	"context"
	"fmt"
	"net/url"
	"testing"
	"time"
)

func mustParse(raw string) *url.URL {
	u, err := url.Parse(raw)
	if err != nil {
		panic(err)
	}
	return u
}

func TestWorkerPool(t *testing.T) {
	t.Run("processes submitted work and returns results", func(t *testing.T) {
		wp := NewWorkerPool(2)
		wp.Submit(mustParse("https://example.com/a"))
		wp.Start(context.Background(), func(ctx context.Context, u *url.URL) ([]*url.URL, error) {
			return []*url.URL{mustParse("/found")}, nil
		})

		result := <-wp.Results()
		if result.Err != nil {
			t.Fatalf("unexpected error: %v", result.Err)
		}
		if result.SourceUrl.String() != "https://example.com/a" {
			t.Errorf("source = %q, want %q", result.SourceUrl, "https://example.com/a")
		}
		if len(result.DiscoverdUrls) != 1 || result.DiscoverdUrls[0].String() != "/found" {
			t.Errorf("discovered = %v, want [/found]", result.DiscoverdUrls)
		}
		wp.Done()
	})

	t.Run("propagates worker errors", func(t *testing.T) {
		wp := NewWorkerPool(1)
		wp.Submit(mustParse("https://example.com"))
		wp.Start(context.Background(), func(ctx context.Context, u *url.URL) ([]*url.URL, error) {
			return nil, fmt.Errorf("fetch failed")
		})

		result := <-wp.Results()
		if result.Err == nil {
			t.Fatal("expected error, got nil")
		}
		wp.Done()
	})

	t.Run("Done shuts down when pending reaches zero", func(t *testing.T) {
		wp := NewWorkerPool(1)
		wp.Submit(mustParse("https://example.com"))
		wp.Start(context.Background(), func(ctx context.Context, u *url.URL) ([]*url.URL, error) {
			return nil, nil
		})

		<-wp.Results()
		wp.Done()

		// results channel should be closed now
		_, ok := <-wp.Results()
		if ok {
			t.Error("expected results channel to be closed")
		}
	})

	t.Run("processes multiple items", func(t *testing.T) {
		wp := NewWorkerPool(2)
		wp.Start(context.Background(), func(ctx context.Context, u *url.URL) ([]*url.URL, error) {
			return nil, nil
		})

		wp.Submit(mustParse("https://example.com/1"))
		wp.Submit(mustParse("https://example.com/2"))
		wp.Submit(mustParse("https://example.com/3"))

		count := 0
		for range wp.Results() {
			count++
			wp.Done()
		}
		if count != 3 {
			t.Errorf("got %d results, want 3", count)
		}
	})

	t.Run("workers exit on cancelled context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		wp := NewWorkerPool(2)
		wp.Submit(mustParse("https://example.com"))
		wp.Start(ctx, func(ctx context.Context, u *url.URL) ([]*url.URL, error) {
			<-ctx.Done()
			return nil, ctx.Err()
		})

		cancel()

		// should not hang — Shutdown closes channels
		time.Sleep(50 * time.Millisecond)
		wp.Shutdown()
	})

	t.Run("Shutdown is safe to call multiple times", func(t *testing.T) {
		wp := NewWorkerPool(1)
		wp.Start(context.Background(), func(ctx context.Context, u *url.URL) ([]*url.URL, error) {
			return nil, nil
		})

		wp.Shutdown()
		wp.Shutdown() // should not panic
	})
}
