package api_test

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"keystone-asset-attestation/internal/api"
	"keystone-asset-attestation/internal/core"
)

// GET /v1/modules
func TestBug026ModulePanicStaysIsolated(t *testing.T) {
	handler := api.New(nil, slog.Default()).Handler()
	defer func() {
		if recovered := recover(); recovered != nil {
			t.Fatalf("panic escaped handler: %v", recovered)
		}
	}()
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/v1/modules", nil))
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("panic response status = %d, want 500", rec.Code)
	}
}

func TestBug026ModulePanicRegression(t *testing.T) {
	handler := api.New(core.NewEngine(), slog.Default()).Handler()
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/v1/modules", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("normal module status = %d", rec.Code)
	}
}
