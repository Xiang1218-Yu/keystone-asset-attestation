package core

import (
	"context"
	"testing"
)

func TestBug022HistorySliceDoesNotPolluteRecord(t *testing.T) {
	engine := NewEngine()
	_, _ = engine.Create(context.Background(), "history", "asset supplier lineage inspection baseline", "tester")
	record, err := engine.Advance(context.Background(), "history", "verified", "tester")
	if err != nil || len(record.History) != 2 {
		t.Fatalf("advance history = %d, err=%v", len(record.History), err)
	}
	record.History[0].To = "rewritten"
	got, _ := engine.Get(context.Background(), "history")
	if got.History[0].To == "rewritten" {
		t.Fatal("history mutation polluted stored record")
	}
}

func TestBug022HistorySliceRegression(t *testing.T) {
	engine := NewEngine()
	_, _ = engine.Create(context.Background(), "history-ok", "asset supplier lineage inspection baseline", "tester")
	record, err := engine.Advance(context.Background(), "history-ok", "verified", "tester")
	if err != nil || len(record.History) != 2 {
		t.Fatalf("history regression failed: %v", err)
	}
}
