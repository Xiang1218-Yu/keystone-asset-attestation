package ops

import (
	"context"
	"testing"
	"time"
)

func TestBug018QueueStopCancelsRunningHandler(t *testing.T) {
	queue := NewQueue(1, 1)
	started := make(chan struct{})
	cancelled := make(chan bool, 1)
	queue.Register("wait", func(ctx context.Context, _ Job) error {
		close(started)
		select {
		case <-ctx.Done():
			cancelled <- true
		case <-time.After(100 * time.Millisecond):
			cancelled <- false
		}
		return nil
	})
	if !queue.Submit(Job{Name: "wait"}) {
		t.Fatal("submit failed")
	}
	<-started
	queue.Stop()
	if !<-cancelled {
		t.Fatal("running handler did not receive queue cancellation")
	}
}

func TestBug018QueueSubmitRegression(t *testing.T) {
	queue := NewQueue(1, 1)
	done := make(chan struct{})
	queue.Register("work", func(context.Context, Job) error { close(done); return nil })
	if !queue.Submit(Job{Name: "work"}) {
		t.Fatal("submit failed")
	}
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("job did not run")
	}
	queue.Stop()
}
