package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"keystone-asset-attestation/internal/core"
)

func decode(w http.ResponseWriter, r *http.Request, target any) bool {
	defer r.Body.Close()
	decoder := json.NewDecoder(io.LimitReader(r.Body, 1<<20))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return false
	}
	return true
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func limit(value string, fallback int) int {
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < 1 {
		return fallback
	}
	if parsed > 1000 {
		return 1000
	}
	return parsed
}

// actor resolves the operator for a request. An explicit "x-operator" header
// is preserved verbatim so interactive callers are always attributed
// correctly. When the header is absent — the case for automated import tasks
// that do not carry operator information — it falls back to core.DefaultActor
// so the operator default travels from the request edge into record history.
func actor(r *http.Request) string {
	if value := strings.TrimSpace(r.Header.Get("x-operator")); value != "" {
		return value
	}
	return core.DefaultActor
}
