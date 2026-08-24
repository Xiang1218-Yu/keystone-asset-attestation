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

// POST /v1/records/{id}/advance
func TestBug002AdvanceCancellationKeepsStage(t *testing.T) {
	handler := bugHandler()
	create := bugJSONRequest(t, handler, http.MethodPost, "/v1/records", map[string]string{"id": "cancel-advance", "payload": "asset supplier lineage inspection baseline"})
	if create.Code != http.StatusCreated {
		t.Fatalf("create status = %d", create.Code)
	}
	req := httptest.NewRequest(http.MethodPost, "/v1/records/cancel-advance/advance", strings.NewReader(`{"stage":"verified"}`))
	ctx, cancel := context.WithCancel(req.Context())
	cancel()
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("cancelled advance status = %d, want 400", rec.Code)
	}
	got := bugRequest(t, handler, http.MethodGet, "/v1/records/cancel-advance", "")
	if !strings.Contains(bugReadBody(t, got), `"stage":"captured"`) {
		t.Fatalf("cancelled request changed record: %s", bugReadBody(t, got))
	}
}

func TestBug002AdvanceValidRequestRegression(t *testing.T) {
	handler := bugHandler()
	create := bugJSONRequest(t, handler, http.MethodPost, "/v1/records", map[string]string{"id": "valid-advance", "payload": "asset supplier lineage inspection baseline"})
	if create.Code != http.StatusCreated {
		t.Fatalf("create status = %d", create.Code)
	}
	rec := bugJSONRequest(t, handler, http.MethodPost, "/v1/records/valid-advance/advance", map[string]string{"stage": "verified"})
	if rec.Code != http.StatusOK || !strings.Contains(bugReadBody(t, rec), `"stage":"verified"`) {
		t.Fatalf("valid advance failed: %d", rec.Code)
	}
}
