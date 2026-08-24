package api

import (
	"net/http"

	"keystone-asset-attestation/internal/core"
)

type recordInput struct {
	ID      string `json:"id"`
	Payload string `json:"payload"`
}

type advanceInput struct {
	Stage string `json:"stage"`
}

func (s *Server) listRecords(w http.ResponseWriter, r *http.Request) {
	if !s.requireEngine(w) {
		return
	}
	records, err := s.engine.List(r.Context(), r.URL.Query().Get("stage"), limit(r.URL.Query().Get("limit"), 100))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": records, "count": len(records)})
}

func (s *Server) createRecord(w http.ResponseWriter, r *http.Request) {
	var input recordInput
	if !decode(w, r, &input) {
		return
	}
	if !s.requireEngine(w) {
		return
	}
	record, err := s.engine.Create(r.Context(), input.ID, input.Payload, actor(r))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, record)
}

func (s *Server) getRecord(w http.ResponseWriter, r *http.Request) {
	if !s.requireEngine(w) {
		return
	}
	record, err := s.engine.Get(r.Context(), r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, record)
}

func (s *Server) advanceRecord(w http.ResponseWriter, r *http.Request) {
	var input advanceInput
	if !decode(w, r, &input) {
		return
	}
	if !s.requireEngine(w) {
		return
	}
	record, err := s.engine.Advance(r.Context(), r.PathValue("id"), input.Stage, actor(r))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, record)
}

func _keepCoreType(_ core.Record) {}
