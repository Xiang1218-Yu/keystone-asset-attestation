package ops

import (
	"sync"
	"testing"
	"time"
)

func TestBug016ConcurrentAuditIdsStayUnique(t *testing.T) {
	log := NewAuditLog()
	const total = 40
	entries := make(chan AuditEntry, total)
	var wg sync.WaitGroup
	for i := 0; i < total; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); entries <- log.Add("operator", "inspect", "asset", time.Unix(1, 2)) }()
	}
	wg.Wait()
	close(entries)
	seen := map[string]bool{}
	for entry := range entries {
		if seen[entry.ID] {
			t.Fatalf("duplicate audit id: %s", entry.ID)
		}
		seen[entry.ID] = true
	}
}

func TestBug016AuditListRegression(t *testing.T) {
	log := NewAuditLog()
	log.Add("operator", "inspect", "asset", time.Unix(1, 2))
	if len(log.List(10)) != 1 {
		t.Fatalf("audit list length = %d, want 1", len(log.List(10)))
	}
}
