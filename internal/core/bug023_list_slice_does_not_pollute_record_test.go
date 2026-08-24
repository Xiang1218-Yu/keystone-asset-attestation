package core

import (
	"context"
	"testing"
)

// GET /v1/records
func TestBug023ListSliceDoesNotPolluteRecord(t *testing.T) {
	engine := NewEngine()
	_, _ = engine.Create(context.Background(), "list", "asset supplier lineage inspection baseline", "tester")
	records, err := engine.List(context.Background(), "", 10)
	if err != nil || len(records) != 1 {
		t.Fatalf("list failed: %v", err)
	}
	records[0].History[0].To = "rewritten"
	got, _ := engine.Get(context.Background(), "list")
	if got.History[0].To == "rewritten" {
		t.Fatal("list result polluted stored record")
	}
}

func TestBug023ListSliceRegression(t *testing.T) {
	engine := NewEngine()
	_, _ = engine.Create(context.Background(), "list-ok", "asset supplier lineage inspection baseline", "tester")
	records, err := engine.List(context.Background(), "", 10)
	if err != nil || len(records) != 1 {
		t.Fatalf("list regression failed: %v", err)
	}
}
