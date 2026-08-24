package core

import (
	"context"
	"testing"
)

// POST /v1/records
func TestBug021EvidenceSliceDoesNotPolluteRecord(t *testing.T) {
	engine := NewEngine()
	created, err := engine.Create(context.Background(), "evidence", "asset supplier lineage inspection baseline", "tester")
	if err != nil || len(created.Evidence) == 0 {
		t.Fatalf("create failed: %v", err)
	}
	created.Evidence[0].Detail = "changed"
	got, err := engine.Get(context.Background(), "evidence")
	if err != nil || got.Evidence[0].Detail == "changed" {
		t.Fatalf("evidence mutation polluted stored record")
	}
}

func TestBug021EvidenceSliceRegression(t *testing.T) {
	record, err := NewEngine().Create(context.Background(), "evidence-ok", "asset supplier lineage inspection baseline", "tester")
	if err != nil || len(record.Evidence) == 0 {
		t.Fatalf("evidence was not created: %v", err)
	}
}
