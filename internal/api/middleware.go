package api

import (
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"
)

func requestLog(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		started := time.Now()
		next.ServeHTTP(w, r)
		logger.Info("request", "method", r.Method, "path", r.URL.Path, "duration", time.Since(started))
	})
}

// recoverPanic isolates a panic raised by a single request handler so it is
// translated into a 500 response instead of bubbling up through the HTTP
// server (which, because handlers share the engine's global state, would
// otherwise crash the process and break every subsequent request). It is the
// per-request recovery layer; the deferred recover must live inside the handler
// goroutine that runs the panic, so recoverPanic must wrap the mux directly.
func recoverPanic(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if value := recover(); value != nil {
				if logger != nil {
					logger.Error("panic recovered",
						"method", r.Method, "path", r.URL.Path, "panic", value, "stack", string(debug.Stack()))
				}
				// A panic may have happened after a partial write. If the
				// response has not yet started, emit a clean JSON error so the
				// caller gets a well-formed body; if headers were already sent
				// we can only close the connection, which http.Server handles.
				if !headersWritten(w) {
					writeError(w, http.StatusInternalServerError, "internal server error")
				}
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// headersWritten reports whether the response already has a status line or
// headers committed to the wire, so the recovery layer knows whether a clean
// error response is still possible.
func headersWritten(w http.ResponseWriter) bool {
	type written interface{ Written() bool }
	if w, ok := w.(written); ok {
		return w.Written()
	}
	return false
}
