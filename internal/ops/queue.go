package ops

import (
	"context"
	"sync"
)

type Job struct {
	Name    string
	Payload any
}

type Handler func(context.Context, Job) error

type Queue struct {
	jobs     chan Job
	handlers map[string]Handler
	mu       sync.RWMutex
	wg       sync.WaitGroup
	ctx      context.Context
	cancel   context.CancelFunc
}

func NewQueue(buffer, workers int) *Queue {
	if buffer < 1 {
		buffer = 32
	}
	if workers < 1 {
		workers = 1
	}
	ctx, cancel := context.WithCancel(context.Background())
	queue := &Queue{jobs: make(chan Job, buffer), handlers: make(map[string]Handler), ctx: ctx, cancel: cancel}
	for index := 0; index < workers; index++ {
		queue.wg.Add(1)
		go queue.worker()
	}
	return queue
}

func (q *Queue) Register(name string, handler Handler) {
	q.mu.Lock()
	q.handlers[name] = handler
	q.mu.Unlock()
}

func (q *Queue) Submit(job Job) bool {
	select {
	case q.jobs <- job:
		return true
	default:
		return false
	}
}

func (q *Queue) Stop() {
	q.cancel()
	q.wg.Wait()
}

func (q *Queue) worker() {
	defer q.wg.Done()
	for {
		select {
		case <-q.ctx.Done():
			return
		case job := <-q.jobs:
			// Snapshot the handler under the read lock so the map lookup is
			// never concurrent with Register's write. The lock is released
			// before invoking the handler so that an administrator can keep
			// (re)registering handlers during a long-running job without
			// blocking on it, and an already-dispatched job still executes.
			q.mu.RLock()
			handler := q.handlers[job.Name]
			q.mu.RUnlock()
			if handler != nil {
				_ = handler(q.ctx, job)
			}
		}
	}
}
