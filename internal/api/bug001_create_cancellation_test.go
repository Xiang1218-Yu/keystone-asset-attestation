package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"keystone-asset-attestation/internal/api"
	"keystone-asset-attestation/internal/core"
)

func bugHandler() http.Handler {
	return api.New(core.NewEngine(), slog.Default()).Handler()
}

func bugRequest(t *testing.T, handler http.Handler, method, target, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, target, strings.NewReader(body))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec
}

func bugJSONRequest(t *testing.T, handler http.Handler, method, target string, value any) *httptest.ResponseRecorder {
	t.Helper()
	raw, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	return bugRequest(t, handler, method, target, string(raw))
}

func bugReadBody(t *testing.T, rec *httptest.ResponseRecorder) string {
	t.Helper()
	raw, err := io.ReadAll(rec.Result().Body)
	if err != nil {
		t.Fatal(err)
	}
	return string(raw)
}

var _ = bytes.NewBuffer
var _ = context.Background
var _ = fmt.Sprintf

// POST /v1/records
func TestBug001CreateCancellationStopsMutation(t *testing.T) {
	handler := bugHandler()

	req := httptest.NewRequest(http.MethodPost, "/v1/records", strings.NewReader(`{"id":"cancelled","payload":"asset supplier lineage inspection baseline"}`))
	ctx, cancel := context.WithCancel(req.Context())
	cancel()
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("cancelled create status = %d, want 400", rec.Code)
	}
}

func TestBug001CreateCancellationValidRequestRegression(t *testing.T) {
	handler := bugHandler()
	rec := bugJSONRequest(t, handler, http.MethodPost, "/v1/records", map[string]string{"id": "valid", "payload": "asset supplier lineage inspection baseline"})
	if rec.Code != http.StatusCreated {
		t.Fatalf("valid request status = %d, want 201", rec.Code)
	}

}
