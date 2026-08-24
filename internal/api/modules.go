package api

import (
	"net/http"
)

func (s *Server) modules(w http.ResponseWriter, r *http.Request) {
	values := make([]map[string]any, 0)
	for _, module := range s.engine.Modules() {
		values = append(values, map[string]any{
			"key": module.Key(), "description": module.Description(),
			"priority": module.Priority(), "family": module.Family(), "enabled": module.Enabled(),
		})
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": values, "count": len(values)})
}

func (s *Server) snapshot(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.engine.Snapshot())
}

// listAudit renders the audit page. Because Add always resolves a blank actor
// to core.DefaultActor, every entry here carries a non-empty actor and the
// page can render it directly without special-casing empty data.
func (s *Server) listAudit(w http.ResponseWriter, r *http.Request) {
	entries := s.audit.List(limit(r.URL.Query().Get("limit"), 100))
	writeJSON(w, http.StatusOK, map[string]any{"items": entries, "count": len(entries)})
}
