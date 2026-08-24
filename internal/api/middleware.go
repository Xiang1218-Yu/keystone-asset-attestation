package api

import (
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"
)

// statusRecorder captures the status code and whether the response header has
// already been committed. It lets the logging and recovery middleware report
// the real outcome of a request even when the handler panics partway through.
type statusRecorder struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
	panicRecovered bool
}

func newStatusRecorder(w http.ResponseWriter) *statusRecorder {
	return &statusRecorder{ResponseWriter: w, status: http.StatusOK}
}

func (r *statusRecorder) WriteHeader(code int) {
	if r.wroteHeader {
		return
	}
	r.status = code
	r.wroteHeader = true
	r.ResponseWriter.WriteHeader(code)
}

func (r *statusRecorder) Write(b []byte) (int, error) {
	if !r.wroteHeader {
		r.wroteHeader = true
	}
	return r.ResponseWriter.Write(b)
}

func requestLog(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		started := time.Now()
		rec := newStatusRecorder(w)
		// Defer the log so a panic in the wrapped handler still produces a
		// record; without this, the eager log after ServeHTTP is skipped on
		// the exact requests that most need leaving a trail.
		defer func() {
			attrs := []any{
				"method", r.Method,
				"path", r.URL.Path,
				"status", rec.status,
				"duration", time.Since(started),
			}
			if rec.panicRecovered {
				attrs = append(attrs, "panic", true)
			}
			logger.Info("request", attrs...)
		}()
		next.ServeHTTP(rec, r)
	})
}

func recoverPanic(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			value := recover()
			if value == nil {
				return
			}
			// Capture the stack at the panic site before anything else.
			stack := debug.Stack()
			// recoverPanic runs directly inside requestLog, which passes us
			// the statusRecorder it owns. Mark it so the deferred log line
			// records that this 500 came from a recovered panic, not a normal
			// error response, and force the reported status to 500.
			rec, _ := w.(*statusRecorder)
			logger.Error("panic recovered",
				"method", r.Method,
				"path", r.URL.Path,
				"value", value,
				"stack", string(stack),
			)
			// If the handler already committed the response header we cannot
			// replace the body with a clean 500; the connection is left in a
			// partially-written state and the existing status is reported.
			if rec != nil {
				rec.panicRecovered = true
				if rec.wroteHeader {
					return
				}
			}
			writeError(w, http.StatusInternalServerError, "internal server error")
		}()
		next.ServeHTTP(w, r)
	})
}
