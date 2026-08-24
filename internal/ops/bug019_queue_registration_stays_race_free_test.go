package ops

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestBug019QueueRegistrationStaysRaceFree(t *testing.T) {
	queue := NewQueue(32, 2)
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			queue.Register("job", func(context.Context, Job) error { return nil })
			queue.Submit(Job{Name: "job"})
		}(i)
	}
	wg.Wait()
	time.Sleep(10 * time.Millisecond)
	queue.Stop()
}

func TestBug019QueueWorkerRegression(t *testing.T) {
	queue := NewQueue(1, 1)
	done := make(chan struct{})
	queue.Register("job", func(context.Context, Job) error { close(done); return nil })
	queue.Submit(Job{Name: "job"})
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("job did not run")
	}
	queue.Stop()
}
