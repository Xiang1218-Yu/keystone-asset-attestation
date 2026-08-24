package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"keystone-asset-attestation/internal/core"
	"keystone-asset-attestation/internal/ops"
)

type recordInput struct {
	ID      string `json:"id"`
	Payload string `json:"payload"`
}

type advanceInput struct {
	Stage string `json:"stage"`
}

func (s *Server) listRecords(w http.ResponseWriter, r *http.Request) {
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
	requestContext := r.Context()
	operator := actor(r)
	record, err := s.engine.Create(requestContext, input.ID, input.Payload, operator)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	// Record the same operator in the audit log so the audit page matches the
	// record history. operator is already defaulted by actor(r) for automated
	// imports; Add also defends against a blank actor defensively.
	s.audit.Add(operator, "record.create", record.ID, time.Now().UTC())
	writeJSON(w, http.StatusCreated, record)
}

func (s *Server) getRecord(w http.ResponseWriter, r *http.Request) {
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
	operator := actor(r)
	record, err := s.engine.Advance(r.Context(), r.PathValue("id"), input.Stage, operator)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	s.audit.Add(operator, "record.advance", record.ID, time.Now().UTC())
	writeJSON(w, http.StatusOK, record)
}

// importPayload mirrors recordInput for automated imports queued without an
// HTTP request and therefore without an "x-operator" header.
type importPayload struct {
	ID      string `json:"id"`
	Payload string `json:"payload"`
}

// importHandler is the automated import task. It creates records without an
// operator, simulating an importer that carries no operator information. The
// operator default travels through engine.Create (history "by") and
// audit.Add (audit actor) so neither the record history nor the audit page
// ever observes an empty actor.
func (s *Server) importHandler() func(ctx context.Context, job ops.Job) error {
	return func(ctx context.Context, job ops.Job) error {
		data, err := json.Marshal(job.Payload)
		if err != nil {
			return err
		}
		var payload importPayload
		if err := json.Unmarshal(data, &payload); err != nil {
			return err
		}
		record, err := s.engine.Create(ctx, payload.ID, payload.Payload, "")
		if err != nil {
			return err
		}
		s.audit.Add("", "record.import", record.ID, time.Now().UTC())
		return nil
	}
}

// submitImport enqueues an automated import task. It intentionally reads no
// "x-operator" header, so the operator default flows from engine.Create into
// record history and audit. Explicit operators are never supplied here, so
// interactive attribution is unaffected.
func (s *Server) submitImport(w http.ResponseWriter, r *http.Request) {
	var payload importPayload
	if !decode(w, r, &payload) {
		return
	}
	if !s.queue.Submit(ops.Job{Name: "import", Payload: payload}) {
		writeError(w, http.StatusServiceUnavailable, "import queue is full")
		return
	}
	writeJSON(w, http.StatusAccepted, map[string]any{"status": "queued", "id": payload.ID})
}

func _keepCoreType(_ core.Record) {}
