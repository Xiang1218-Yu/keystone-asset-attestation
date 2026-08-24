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
func TestBug029AnonymousActorStaysSafe(t *testing.T) {
	handler := bugHandler()
	rec := bugJSONRequest(t, handler, http.MethodPost, "/v1/records", map[string]string{"id": "anon", "payload": "asset supplier lineage inspection baseline"})
	if !strings.Contains(bugReadBody(t, rec), `"by":"anonymous"`) {
		t.Fatalf("anonymous actor missing: %s", bugReadBody(t, rec))
	}
}

func TestBug029NamedActorRegression(t *testing.T) {
	handler := bugHandler()
	req := httptest.NewRequest(http.MethodPost, "/v1/records", strings.NewReader(`{"id":"named","payload":"asset supplier lineage inspection baseline"}`))
	req.Header.Set("x-operator", "auditor")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if !strings.Contains(bugReadBody(t, rec), `"by":"auditor"`) {
		t.Fatalf("named actor missing")
	}
}
