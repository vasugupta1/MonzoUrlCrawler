package concurrency

import (
	"context"
	"net/url"
	"sync"
	"sync/atomic"
)

type WorkResult struct {
	DiscoverdUrls []*url.URL
	SourceUrl     *url.URL
	Err           error
}

type WorkerPool struct {
	workerCount int
	queue       chan *url.URL
	results     chan WorkResult
	pendingWork atomic.Int64
	closeOnce   sync.Once
}

func NewWorkerPool(workerCount int) *WorkerPool {
	return &WorkerPool{
		workerCount: workerCount,
		queue:       make(chan *url.URL, workerCount),
		results:     make(chan WorkResult),
	}
}

func (wp *WorkerPool) Submit(url *url.URL) {
	wp.pendingWork.Add(1)
	wp.queue <- url
}

func (wp *WorkerPool) Results() <-chan WorkResult {
	return wp.results
}

func (wp *WorkerPool) Start(ctx context.Context, function func(ctx context.Context, url *url.URL) ([]*url.URL, error)) {
	for i := 0; i < wp.workerCount; i++ {
		go wp.worker(ctx, function)
	}
}

func (wp *WorkerPool) Done() {
	//subtracting and checking is 2 atomic operations, its better read the return from .Add
	if wp.pendingWork.Add(-1) == 0 {
		wp.shutdown()
		return
	}
}

func (wp *WorkerPool) Shutdown() {
	wp.shutdown()
}

// closeOnce will use sync.Once to execute closing of the channels at once and its atomic
func (wp *WorkerPool) shutdown() {
	wp.closeOnce.Do(func() {
		close(wp.queue)
		close(wp.results)
	})
}

func (wp *WorkerPool) worker(ctx context.Context, function func(ctx context.Context, url *url.URL) ([]*url.URL, error)) {
	for {
		select {
		case <-ctx.Done():
			return
		case targetUrl, ok := <-wp.queue:
			if !ok {
				//queue is closed here so exit cleanly
				return
			}
			urls, err := function(ctx, targetUrl)

			// edge case can happen where while we were trying to add a result to chan ctx comes and cancel so we should check that
			select {
			case <-ctx.Done():
				return
			case wp.results <- WorkResult{DiscoverdUrls: urls, SourceUrl: targetUrl, Err: err}:
			}
		}
	}
}
