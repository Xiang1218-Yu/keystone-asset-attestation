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

// GET /v1/snapshot
func TestBug024SnapshotCountsRecordsNotVersions(t *testing.T) {
	handler := bugHandler()
	bugJSONRequest(t, handler, http.MethodPost, "/v1/records", map[string]string{"id": "snapshot-state", "payload": "asset supplier lineage inspection baseline"})
	bugJSONRequest(t, handler, http.MethodPost, "/v1/records/snapshot-state/advance", map[string]string{"stage": "verified"})
	rec := bugRequest(t, handler, http.MethodGet, "/v1/snapshot", "")
	if !strings.Contains(bugReadBody(t, rec), `"stage_verified":1`) {
		t.Fatalf("snapshot count was version-weighted: %s", bugReadBody(t, rec))
	}
}

func TestBug024SnapshotRegression(t *testing.T) {
	rec := bugRequest(t, bugHandler(), http.MethodGet, "/v1/snapshot", "")
	if rec.Code != http.StatusOK || !strings.Contains(bugReadBody(t, rec), `"records":0`) {
		t.Fatalf("empty snapshot regression failed: %d", rec.Code)
	}
}
