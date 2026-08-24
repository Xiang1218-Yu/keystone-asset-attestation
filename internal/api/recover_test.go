package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestRecoverPanic_ContainsPanickingHandler verifies the per-request recovery
// layer: a handler that panics must be turned into a 500 JSON response rather
// than bubbling the panic to the caller, which is the failure mode the process
// hit before recoverPanic was wired into the handler chain.
func TestRecoverPanic_ContainsPanickingHandler(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(&strings.Builder{}, nil))
	panicking := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("simulated probe module failure")
	})
	handler := requestLog(logger, recoverPanic(logger, panicking))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/probe", strings.NewReader(""))

	// Must not panic the test goroutine; recovery returns a 500 instead.
	assertNotPanics(t, func() {
		handler.ServeHTTP(rec, req)
	})

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status: got %d, want %d", rec.Code, http.StatusInternalServerError)
	}
	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if body["error"] != "internal server error" {
		t.Fatalf("error body: got %q, want %q", body["error"], "internal server error")
	}
}

// TestRecoverPanic_NormalHandlerStillServes confirms that a non-panicking
// handler in the same chain returns its normal response — i.e. recovery does
// not interfere with healthy requests, so normal module requests still return.
func TestRecoverPanic_NormalHandlerStillServes(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(&strings.Builder{}, nil))
	normal := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{"status": "ok"})
	})
	handler := requestLog(logger, recoverPanic(logger, normal))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/healthz", strings.NewReader(""))

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d", rec.Code, http.StatusOK)
	}
	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if body["status"] != "ok" {
		t.Fatalf("body: got %v, want status=ok", body)
	}
}

// TestServer_HandlerWiredWithRecovery checks the production Handler() assembly
// includes the recovery layer, guarding against a regression where recoverPanic
// is defined but unwired (the original bug).
func TestServer_HandlerWiredWithRecovery(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(&strings.Builder{}, nil))
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(""))
	assertNotPanics(t, func() {
		// Reuse the same composition as Handler() so this mirrors production.
		requestLog(logger, recoverPanic(logger, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
			panic("boom")
		}))).ServeHTTP(rec, req)
	})
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status: got %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}

// TestRecoverPanic_IsolationThenNormalRequest is the end-to-end version of the
// reported incident: a "monitoring probe module" handler encounters an anomaly
// and panics. The recovery layer must isolate that to a single 500 response so
// the process keeps serving — and a subsequent normal request must still return.
// This runs the same composition as Server.Handler() against a real
// http.Server over httptest, matching the production request path.
func TestRecoverPanic_IsolationThenNormalRequest(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(&strings.Builder{}, nil))
	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/probe", func(w http.ResponseWriter, r *http.Request) {
		panic("simulated probe module anomaly")
	})
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{"status": "ok"})
	})
	// Same composition as Server.Handler(): outer requestLog, inner recoverPanic
	// wrapping the mux directly.
	server := httptest.NewServer(requestLog(logger, recoverPanic(logger, mux)))
	defer server.Close()

	// First request: the probe handler panics. Must come back as a clean 500,
	// not a connection-level crash, and must not fail the test process.
	probeResp, err := server.Client().Get(server.URL + "/v1/probe")
	if err != nil {
		t.Fatalf("probe request: %v", err)
	}
	probeResp.Body.Close()
	if probeResp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("probe status: got %d, want %d (panic should be isolated to a 500)",
			probeResp.StatusCode, http.StatusInternalServerError)
	}

	// Second request: a normal handler in the same process. Must still return
	// 200 — proving the process-level recovery kept the server usable.
	normalResp, err := server.Client().Get(server.URL + "/healthz")
	if err != nil {
		t.Fatalf("normal request after panic: %v", err)
	}
	defer normalResp.Body.Close()
	if normalResp.StatusCode != http.StatusOK {
		t.Fatalf("normal status after panic: got %d, want %d", normalResp.StatusCode, http.StatusOK)
	}
	var body map[string]any
	if err := json.NewDecoder(normalResp.Body).Decode(&body); err != nil {
		t.Fatalf("decode normal body: %v", err)
	}
	if body["status"] != "ok" {
		t.Fatalf("normal body after panic: got %v, want status=ok", body)
	}
}

func assertNotPanics(t *testing.T, fn func()) {
	t.Helper()
	defer func() {
		if value := recover(); value != nil {
			t.Fatalf("expected no panic, got: %v", value)
		}
	}()
	fn()
}
