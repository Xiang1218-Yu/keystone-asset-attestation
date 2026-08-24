package ops

import (
	"testing"
	"time"
)

func TestBug030AuditListDoesNotPolluteLog(t *testing.T) {
	log := NewAuditLog()
	log.Add("operator", "inspect", "asset", time.Unix(1, 2))
	entries := log.List(10)
	entries[0].Action = "rewritten"
	if log.List(10)[0].Action == "rewritten" {
		t.Fatal("audit list mutation polluted log")
	}
}

func TestBug030AuditListRegression(t *testing.T) {
	log := NewAuditLog()
	log.Add("operator", "inspect", "asset", time.Unix(1, 2))
	if len(log.List(10)) != 1 {
		t.Fatalf("audit list length = %d, want 1", len(log.List(10)))
	}
}
