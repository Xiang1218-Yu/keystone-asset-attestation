package api_test

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"keystone-asset-attestation/internal/api"
	"keystone-asset-attestation/internal/core"
)

// GET /v1/modules
func TestBug025NilEngineIsRecovered(t *testing.T) {
	handler := api.New(nil, slog.Default()).Handler()
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/v1/modules", nil))
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("nil engine status = %d, want 500", rec.Code)
	}
}

func TestBug025ModulesRegression(t *testing.T) {
	handler := api.New(core.NewEngine(), slog.Default()).Handler()
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/v1/modules", nil))
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"count":62`) {
		t.Fatalf("modules response invalid: %d", rec.Code)
	}
}
