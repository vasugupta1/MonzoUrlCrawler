package concurrency

import (
	"context"
	"net/url"
	"sync/atomic"
)

type WorkResult struct {
	discoverdUrls []*url.URL
	sourceUrl     *url.URL
}

type WorkerPool struct {
	workerCount int
	queue       chan *url.URL
	results     chan WorkResult
	pendingWork atomic.Int64
}

func NewWorkerPool(workerCount int) *WorkerPool {
	return &WorkerPool{
		workerCount: workerCount,
		queue: make(chan *url.URL, workerCount*10),
		results:     make(chan WorkResult)

	}
}

func (wp *WorkerPool) Submit(url *url.URL){
	wp.pendingWork.Add(1)
	wp.queue <- url
}

func (wp *WorkerPool) Results() <-chan WorkResult{
	return wp.results
}


func (wp *WorkerPool) Start(ctx context.Context, function func(ctx context.Context, url *url.URL) [] *url.URL){
	for i := 0; i < wp.workerCount; i++{
		go wp.worker(ctx, function)
	}
}

