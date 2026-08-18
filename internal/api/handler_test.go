package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"
)

func newTestServer() *Server {
	web := fstest.MapFS{
		"index.html": &fstest.MapFile{Data: []byte("<html></html>")},
	}
	examples := fstest.MapFS{
		"iir_bw4.json": &fstest.MapFile{Data: []byte(`{"kind":"iir","order":4,"cutoff":0.2}`)},
	}
	return New(web, examples)
}

func post(t *testing.T, s *Server, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)
	return rec
}

func TestAPIDesignFIR(t *testing.T) {
	s := newTestServer()
	rec := post(t, s, "/api/design", `{"kind":"fir","order":30,"cutoff":0.2,"window":"hamming"}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d: %s", rec.Code, rec.Body.String())
	}
	var out struct {
		Kind string    `json:"kind"`
		B    []float64 `json:"b"`
		A    []float64 `json:"a"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatal(err)
	}
	if out.Kind != "fir" || len(out.B) != 31 {
		t.Fatalf("unexpected design output: %+v", out)
	}
}

func TestAPIDesignIIR(t *testing.T) {
	s := newTestServer()
	rec := post(t, s, "/api/design", `{"kind":"iir","order":4,"cutoff":0.2}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d: %s", rec.Code, rec.Body.String())
	}
	var out struct {
		Kind string    `json:"kind"`
		A    []float64 `json:"a"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatal(err)
	}
	if out.Kind != "iir" || len(out.A) != 5 {
		t.Fatalf("unexpected iir output: %+v", out)
	}
}

func TestAPIResponse(t *testing.T) {
	s := newTestServer()
	rec := post(t, s, "/api/response", `{"b":[1,0],"a":[1],"points":64}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d: %s", rec.Code, rec.Body.String())
	}
	var out struct {
		Freq        []float64 `json:"freq"`
		MagnitudeDB []float64 `json:"mag_db"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatal(err)
	}
	if len(out.Freq) != 64 || len(out.MagnitudeDB) != 64 {
		t.Fatalf("response length mismatch")
	}
}

func TestAPIZPlane(t *testing.T) {
	s := newTestServer()
	rec := post(t, s, "/api/zplane", `{"b":[1,0,0],"a":[1,-2.1]}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d: %s", rec.Code, rec.Body.String())
	}
	var out struct {
		Stable bool `json:"stable"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatal(err)
	}
	if out.Stable {
		t.Fatal("pole at 2.1 must be unstable")
	}
}

func TestAPIErrorJSON(t *testing.T) {
	s := newTestServer()
	rec := post(t, s, "/api/design", `{"kind":"fir","order":30,"cutoff":0.2,"window":"kaiser"}`)
	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status %d, want 422", rec.Code)
	}
	var out ErrResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatal(err)
	}
	if out.Code != "unknown_window" {
		t.Fatalf("code = %q, want unknown_window", out.Code)
	}
}

func TestAPIDesignBadCutoff(t *testing.T) {
	s := newTestServer()
	rec := post(t, s, "/api/design", `{"kind":"iir","order":4,"cutoff":0.9}`)
	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status %d, want 422", rec.Code)
	}
}

func TestAPIStaticIndex(t *testing.T) {
	s := newTestServer()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d", rec.Code)
	}
}
