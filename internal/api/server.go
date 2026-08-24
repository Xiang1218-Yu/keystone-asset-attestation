package api

import (
	"log/slog"
	"net/http"
	"time"

	"keystone-asset-attestation/internal/core"
	"keystone-asset-attestation/internal/ops"
)

type Server struct {
	engine *core.Engine
	audit  *ops.AuditLog
	queue  *ops.Queue
	logger *slog.Logger
	mux    *http.ServeMux
}

func New(engine *core.Engine, logger *slog.Logger) *Server {
	if logger == nil {
		logger = slog.Default()
	}
	server := &Server{engine: engine, audit: ops.NewAuditLog(), logger: logger, mux: http.NewServeMux()}
	server.queue = ops.NewQueue(64, 2)
	server.queue.Register("import", server.importHandler())
	server.routes()
	return server
}

func (s *Server) Close() {
	if s.queue != nil {
		s.queue.Stop()
	}
}

func (s *Server) Handler() http.Handler {
	return requestLog(s.logger, recoverPanic(s.mux))
}

func (s *Server) routes() {
	s.mux.HandleFunc("GET /healthz", s.health)
	s.mux.HandleFunc("GET /readyz", s.ready)
	s.mux.HandleFunc("GET /v1/modules", s.modules)
	s.mux.HandleFunc("GET /v1/snapshot", s.snapshot)
	s.mux.HandleFunc("GET /v1/audit", s.listAudit)
	s.mux.HandleFunc("POST /v1/validate", s.validate)
	s.mux.HandleFunc("GET /v1/records", s.listRecords)
	s.mux.HandleFunc("POST /v1/records", s.createRecord)
	s.mux.HandleFunc("GET /v1/records/{id}", s.getRecord)
	s.mux.HandleFunc("POST /v1/records/{id}/advance", s.advanceRecord)
	s.mux.HandleFunc("POST /v1/imports", s.submitImport)
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
