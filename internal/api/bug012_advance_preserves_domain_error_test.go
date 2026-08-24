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
func TestBug012AdvancePreservesDomainError(t *testing.T) {
	handler := bugHandler()
	bugJSONRequest(t, handler, http.MethodPost, "/v1/records", map[string]string{"id": "advance-body", "payload": "asset supplier lineage inspection baseline"})
	rec := bugJSONRequest(t, handler, http.MethodPost, "/v1/records/advance-body/advance", map[string]string{"stage": "invalid"})
	if !strings.Contains(bugReadBody(t, rec), "stage is not accepted by module") {
		t.Fatalf("domain error lost: %s", bugReadBody(t, rec))
	}
}

func TestBug012AdvanceSuccessRegression(t *testing.T) {
	handler := bugHandler()
	bugJSONRequest(t, handler, http.MethodPost, "/v1/records", map[string]string{"id": "advance-ok", "payload": "asset supplier lineage inspection baseline"})
	rec := bugJSONRequest(t, handler, http.MethodPost, "/v1/records/advance-ok/advance", map[string]string{"stage": "verified"})
	if rec.Code != http.StatusOK || !strings.Contains(bugReadBody(t, rec), `"version":2`) {
		t.Fatalf("advance response incomplete: %d", rec.Code)
	}
}
