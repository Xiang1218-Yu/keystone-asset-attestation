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
func TestBug015ConcurrentSnapshotReadsStaySafe(t *testing.T) {
	handler := bugHandler()
	start := make(chan struct{})
	done := make(chan struct{}, 20)
	for i := 0; i < 20; i++ {
		go func(i int) {
			<-start
			bugJSONRequest(t, handler, http.MethodPost, "/v1/records", map[string]string{"id": fmt.Sprintf("snap-%d", i), "payload": "asset supplier lineage inspection baseline"})
			done <- struct{}{}
		}(i)
		go func() {
			<-start
			rec := bugRequest(t, handler, http.MethodGet, "/v1/snapshot", "")
			if rec.Code != http.StatusOK {
				t.Errorf("snapshot status = %d", rec.Code)
			}
			done <- struct{}{}
		}()
	}
	close(start)
	for i := 0; i < 40; i++ {
		<-done
	}
}

func TestBug015SnapshotRegression(t *testing.T) {
	rec := bugRequest(t, bugHandler(), http.MethodGet, "/v1/snapshot", "")
	if rec.Code != http.StatusOK || !strings.Contains(bugReadBody(t, rec), `"modules":62`) {
		t.Fatalf("snapshot response invalid: %d", rec.Code)
	}
}
