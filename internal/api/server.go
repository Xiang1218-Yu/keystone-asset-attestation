package api

import (
	"log/slog"
	"net/http"
	"time"

	"keystone-asset-attestation/internal/core"
)

type Server struct {
	engine *core.Engine
	logger *slog.Logger
	mux    *http.ServeMux
}

func New(engine *core.Engine, logger *slog.Logger) *Server {
	if logger == nil {
		logger = slog.Default()
	}
	server := &Server{engine: engine, logger: logger, mux: http.NewServeMux()}
	server.routes()
	return server
}

func (s *Server) Handler() http.Handler {
	// Order matters: requestLog is the outermost layer (it records the total
	// request duration even when the handler panics), and recoverPanic wraps the
	// mux directly so any panic from a module handler is contained to that one
	// request instead of crashing the process. Before this, recoverPanic existed
	// but was never wired in, so a panicking handler bubbled the panic to the
	// caller and left the shared engine unusable for later requests.
	return requestLog(s.logger, recoverPanic(s.logger, s.mux))
}

func (s *Server) routes() {
	s.mux.HandleFunc("GET /healthz", s.health)
	s.mux.HandleFunc("GET /readyz", s.ready)
	s.mux.HandleFunc("GET /v1/modules", s.modules)
	s.mux.HandleFunc("GET /v1/snapshot", s.snapshot)
	s.mux.HandleFunc("POST /v1/validate", s.validate)
	s.mux.HandleFunc("GET /v1/records", s.listRecords)
	s.mux.HandleFunc("POST /v1/records", s.createRecord)
	s.mux.HandleFunc("GET /v1/records/{id}", s.getRecord)
	s.mux.HandleFunc("POST /v1/records/{id}/advance", s.advanceRecord)
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"status": "ok", "service": "keystone-asset-attestation", "time": time.Now().UTC()})
}

func (s *Server) ready(w http.ResponseWriter, r *http.Request) {
	if s.engine == nil {
		writeError(w, http.StatusServiceUnavailable, "engine is not configured")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ready": true})
}
