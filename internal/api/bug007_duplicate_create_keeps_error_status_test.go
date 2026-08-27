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
func TestBug007DuplicateCreateKeepsErrorStatus(t *testing.T) {
	handler := bugHandler()
	payload := map[string]string{"id": "duplicate", "payload": "asset supplier lineage inspection baseline"}
	first := bugJSONRequest(t, handler, http.MethodPost, "/v1/records", payload)
	if first.Code != http.StatusCreated {
		t.Fatalf("first status = %d", first.Code)
	}
	second := bugJSONRequest(t, handler, http.MethodPost, "/v1/records", payload)
	if second.Code != http.StatusBadRequest {
		t.Fatalf("duplicate status = %d, want 400", second.Code)
	}
}

func TestBug007FirstCreateRegression(t *testing.T) {
	rec := bugJSONRequest(t, bugHandler(), http.MethodPost, "/v1/records", map[string]string{"id": "first", "payload": "asset supplier lineage inspection baseline"})
	if rec.Code != http.StatusCreated {
		t.Fatalf("first create status = %d, want 201", rec.Code)
	}
}
