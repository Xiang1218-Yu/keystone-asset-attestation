package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

// closableBody mimics the real http server request body: after Close, further
// reads return an error (real connections behave this way), so a decoder that
// reads from a body closed before Decode reports EOF and treats the request as
// empty. bytes.Reader / strings.Reader used by httptest.NewRequest are no-ops
// on Close and would hide the regression, so we use this instead.
type closableBody struct {
	data   []byte
	offset int
	closed bool
}

func (c *closableBody) Read(p []byte) (int, error) {
	if c.closed {
		return 0, errors.New("read after close: body already closed")
	}
	if c.offset >= len(c.data) {
		return 0, errEOF
	}
	n := copy(p, c.data[c.offset:])
	c.offset += n
	return n, nil
}

func (c *closableBody) Close() error { c.closed = true; return nil }

var errEOF = errors.New("EOF")

func newRequest(body string) *http.Request {
	req := httptest.NewRequest(http.MethodPost, "/v1/records", &closableBody{data: []byte(body)})
	req.Header.Set("content-type", "application/json")
	return req
}

// TestDecodeDoesNotCloseBodyBeforeReading guards against the regression where
// r.Body.Close() ran before json.Decode, so every submission (including batch
// imports) read EOF and was rejected as an empty request.
func TestDecodeDoesNotCloseBodyBeforeReading(t *testing.T) {
	req := newRequest(`{"id":"demo-1","payload":"asset turbine supplier lineage"}`)
	w := httptest.NewRecorder()

	var input recordInput
	if !decode(w, req, &input) {
		t.Fatalf("decode returned false; body closed before read: status=%d body=%s", w.Code, w.Body.String())
	}
	if input.ID != "demo-1" || input.Payload != "asset turbine supplier lineage" {
		t.Fatalf("unexpected decode result: %+v", input)
	}
}

// TestDecodeRejectsBadJSON ensures the deferred close does not mask real errors.
func TestDecodeRejectsBadJSON(t *testing.T) {
	req := newRequest("{not json")
	w := httptest.NewRecorder()

	var input recordInput
	if decode(w, req, &input) {
		t.Fatalf("decode should have failed for invalid JSON")
	}
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

// TestDecodeHandlesMultipleRequests simulates a batch import: many submissions
// over separate requests must all decode, not just the first.
func TestDecodeHandlesMultipleRequests(t *testing.T) {
	for i := 0; i < 5; i++ {
		req := newRequest(`{"id":"batch","payload":"p"}`)
		w := httptest.NewRecorder()
		var input recordInput
		if !decode(w, req, &input) {
			t.Fatalf("request %d failed: %s", i, w.Body.String())
		}
	}
}
