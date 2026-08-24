package ops

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// TestQueueConcurrentRegisterAndConsume exercises the race between an
// administrator dynamically registering handlers while workers are draining
// jobs. Before the fix the worker read the handler map without the read lock,
// which `go test -race` flagged as a concurrent map read and map write; under
// load that race also terminates the worker goroutine mid-job and drops the
// in-flight job. After the fix the workers stay alive and every accepted job
// executes, even while handlers are being (re)registered concurrently.
func TestQueueConcurrentRegisterAndConsume(t *testing.T) {
	queue := NewQueue(64, 4)
	t.Cleanup(queue.Stop)

	// Pre-register a baseline handler so jobs submitted before the dynamic
	// registrations are still executed.
	var executed atomic.Int64
	queue.Register("baseline", func(ctx context.Context, job Job) error {
		executed.Add(1)
		return nil
	})

	// Hammer the queue from one side and register handlers from the other.
	const rounds = 2000
	var accepted atomic.Int64
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < rounds; i++ {
			// Register handlers with the same name repeatedly so the map write
			// keeps racing the worker's map read on the same key — this is
			// what turns an incidental read into a live concurrent read/write
			// pair on the same map.
			queue.Register("dynamic", func(ctx context.Context, job Job) error {
				return nil
			})
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < rounds; i++ {
			if queue.Submit(Job{Name: "baseline", Payload: i}) {
				accepted.Add(1)
			}
		}
	}()
	wg.Wait()

	// Every accepted job must execute — this is the "preserve normal
	// execution of already-registered jobs" guarantee. The workers must not
	// die under the race, which would strand accepted jobs unprocessed.
	deadline := time.Now().Add(5 * time.Second)
	for {
		got := executed.Load()
		want := accepted.Load()
		if got >= want {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("only %d/%d accepted jobs executed; worker died under concurrent Register", got, want)
		}
		time.Sleep(time.Millisecond)
	}
}

func TestQueueRegisterBeforeSubmitExecutes(t *testing.T) {
	queue := NewQueue(4, 1)
	t.Cleanup(queue.Stop)

	var done atomic.Bool
	queue.Register("work", func(ctx context.Context, job Job) error {
		done.Store(true)
		return nil
	})

	if !queue.Submit(Job{Name: "work"}) {
		t.Fatal("submit rejected")
	}

	deadline := time.Now().Add(2 * time.Second)
	for {
		if done.Load() {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("handler never executed")
		}
		time.Sleep(time.Millisecond)
	}
}
