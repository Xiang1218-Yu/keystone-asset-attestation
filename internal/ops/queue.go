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
	case <-q.ctx.Done():
		return false
	default:
	}
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
			q.mu.RLock()
			handler := q.handlers[job.Name]
			q.mu.RUnlock()
			if handler != nil {
				// Run the handler under the queue's context so that Stop
				// cancels in-flight tasks instead of leaving them on a
				// never-cancelable context.Background() to time out.
				_ = handler(q.ctx, job)
			}
		}
	}
}
