package ops

import (
	"strings"
	"sync"
	"testing"
)

func TestBug017ConcurrentMetricIncrementsAreCounted(t *testing.T) {
	metrics := NewMetrics()
	const total = 40
	var wg sync.WaitGroup
	for i := 0; i < total; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); metrics.Inc("processed") }()
	}
	wg.Wait()
	if !strings.Contains(metrics.Text(), "processed 40") {
		t.Fatalf("metric count missing: %s", metrics.Text())
	}
}

func TestBug017MetricTextRegression(t *testing.T) {
	metrics := NewMetrics()
	metrics.Inc("processed")
	if !strings.Contains(metrics.Text(), "processed 1") {
		t.Fatalf("metric text = %s", metrics.Text())
	}
}
