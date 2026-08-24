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
func TestBug014ConcurrentAdvanceKeepsSingleTransition(t *testing.T) {
	handler := bugHandler()
	bugJSONRequest(t, handler, http.MethodPost, "/v1/records", map[string]string{"id": "same-stage", "payload": "asset supplier lineage inspection baseline"})
	const workers = 20
	start := make(chan struct{})
	results := make(chan int, workers)
	for i := 0; i < workers; i++ {
		go func() {
			<-start
			rec := bugJSONRequest(t, handler, http.MethodPost, "/v1/records/same-stage/advance", map[string]string{"stage": "verified"})
			results <- rec.Code
		}()
	}
	close(start)
	success := 0
	for i := 0; i < workers; i++ {
		if <-results == http.StatusOK {
			success++
		}
	}
	if success != 1 {
		t.Fatalf("concurrent advance successes = %d, want 1", success)
	}
}

func TestBug014SingleAdvanceRegression(t *testing.T) {
	handler := bugHandler()
	bugJSONRequest(t, handler, http.MethodPost, "/v1/records", map[string]string{"id": "single-stage", "payload": "asset supplier lineage inspection baseline"})
	rec := bugJSONRequest(t, handler, http.MethodPost, "/v1/records/single-stage/advance", map[string]string{"stage": "verified"})
	if rec.Code != http.StatusOK || !strings.Contains(bugReadBody(t, rec), `"version":2`) {
		t.Fatalf("single advance failed: %d", rec.Code)
	}
}
