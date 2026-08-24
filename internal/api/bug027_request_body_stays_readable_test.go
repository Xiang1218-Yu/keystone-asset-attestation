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
type bug027Body struct {
	reader *strings.Reader
	closed bool
}

func (b *bug027Body) Read(p []byte) (int, error) {
	if b.closed {
		return 0, io.ErrClosedPipe
	}
	return b.reader.Read(p)
}

func (b *bug027Body) Close() error { b.closed = true; return nil }

func TestBug027RequestBodyStaysReadable(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/v1/records", strings.NewReader(`{"id":"body","payload":"asset supplier lineage inspection baseline"}`))
	req.Body = &bug027Body{reader: strings.NewReader(`{"id":"body","payload":"asset supplier lineage inspection baseline"}`)}
	recorder := httptest.NewRecorder()
	bugHandler().ServeHTTP(recorder, req)
	rec := recorder
	if rec.Code != http.StatusCreated {
		t.Fatalf("readable body status = %d, want 201", rec.Code)
	}
}

// GET /healthz
func TestBug027RequestBodyRegression(t *testing.T) {
	rec := bugRequest(t, bugHandler(), http.MethodGet, "/healthz", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("health status = %d, want 200", rec.Code)
	}
}
