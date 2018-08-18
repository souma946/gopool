package gopool

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type (
	// WorkerPool worker pool
	WorkerPool struct {
		size      int
		queue     chan *executedFunction
		executors []*executor
		isClosed  bool
		ctx       context.Context
		cancel    func()
		wg        *sync.WaitGroup
	}

	executedFunction struct {
		f   func(ctx context.Context) error
		ctx context.Context
	}
	executor struct {
		id    int
		queue chan *executedFunction
		wg    *sync.WaitGroup
	}
)

func newExecutor(id int, wg *sync.WaitGroup, queue chan *executedFunction) *executor {
	return &executor{
		id:    id,
		wg:    wg,
		queue: queue,
	}
}

func (e *executor) start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case f := <-e.queue:
				if f == nil {
					return
				}
				f.f(f.ctx)
				e.wg.Done()
				break
			default:
			}
		}
	}()
}

// NewWorkerPool constructor
func NewWorkerPool(size int) *WorkerPool {
	executors := make([]*executor, size)
	queue := make(chan *executedFunction, size)
	wg := new(sync.WaitGroup)
	ctx, cancel := context.WithCancel(context.Background())
	for i := 0; i < size; i++ {
		executors[i] = newExecutor(i, wg, queue)
		executors[i].start(ctx)
	}

	return &WorkerPool{
		size:      size,
		executors: executors,
		queue:     queue,
		ctx:       ctx,
		cancel:    cancel,
		wg:        wg,
	}
}

// Execute execute function async
func (p *WorkerPool) Execute(ctx context.Context, f func(ctx context.Context) error) error {
	if p.isClosed {
		return fmt.Errorf("pool already shutdown")
	}
	p.queue <- &executedFunction{
		f:   f,
		ctx: ctx,
	}
	p.wg.Add(1)
	return nil
}

// Shutdown stop worker pool
func (p *WorkerPool) Shutdown(timeout time.Duration) {
	p.isClosed = true
	done := make(chan struct{})

	go func() {
		p.wg.Wait()
		done <- struct{}{}
	}()
	go func() {
		select {
		case <-time.After(timeout * time.Second):
			done <- struct{}{}
		}
	}()
	close(p.queue)
	for {
		select {
		case <-done:
			p.cancel()
			return
		}
	}
}
