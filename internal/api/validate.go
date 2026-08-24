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
	requestContext := r.Context()
	payload := input.Payload
	score, warnings, err := s.engine.ValidatePayload(requestContext, payload)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"score": score, "warnings": warnings})
}
