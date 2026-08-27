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

// GET /v1/records/{id}
func TestBug005MissingRecordKeepsNotFoundStatus(t *testing.T) {
	rec := bugRequest(t, bugHandler(), http.MethodGet, "/v1/records/missing", "")
	if rec.Code != http.StatusNotFound {
		t.Fatalf("missing record status = %d, want 404", rec.Code)
	}
}

func TestBug005ExistingRecordRegression(t *testing.T) {
	handler := bugHandler()
	bugJSONRequest(t, handler, http.MethodPost, "/v1/records", map[string]string{"id": "existing", "payload": "asset supplier lineage inspection baseline"})
	rec := bugRequest(t, handler, http.MethodGet, "/v1/records/existing", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("existing record status = %d, want 200", rec.Code)
	}
}
