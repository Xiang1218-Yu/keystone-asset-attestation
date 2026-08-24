package core

import "testing"

func TestBug020ModuleSliceDoesNotPolluteEngine(t *testing.T) {
	engine := NewEngine()
	modules := engine.Modules()
	modules[0] = nil
	if engine.Modules()[0] == nil {
		t.Fatal("module slice mutation polluted engine")
	}
}

func TestBug020ModuleSliceRegression(t *testing.T) {
	if len(NewEngine().Modules()) != 62 {
		t.Fatalf("module count = %d, want 62", len(NewEngine().Modules()))
	}
}
