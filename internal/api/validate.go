package api

import "net/http"

type validateInput struct {
	Payload string `json:"payload"`
}

func (s *Server) validate(w http.ResponseWriter, r *http.Request) {
	var input validateInput
	if !decode(w, r, &input) {
		return
	}
	if !s.requireEngine(w) {
		return
	}
	score, warnings, err := s.engine.ValidatePayload(r.Context(), input.Payload)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"score": score, "warnings": warnings})
}
