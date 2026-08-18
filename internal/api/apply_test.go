package api

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestAPIApplyStep(t *testing.T) {
	s := newTestServer()
	rec := post(t, s, "/api/apply", `{"b":[1],"a":[1],"signal":"step","count":8}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d: %s", rec.Code, rec.Body.String())
	}
	var out struct {
		Output []float64 `json:"output"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatal(err)
	}
	if len(out.Output) != 8 {
		t.Fatalf("output length %d, want 8", len(out.Output))
	}
}

func TestAPIApplySine(t *testing.T) {
	s := newTestServer()
	rec := post(t, s, "/api/apply", `{"b":[0.5,0.5],"a":[1],"signal":"sine","count":32,"freq":5}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d: %s", rec.Code, rec.Body.String())
	}
}

func TestAPIApplyBadSignal(t *testing.T) {
	s := newTestServer()
	rec := post(t, s, "/api/apply", `{"b":[1],"a":[1],"signal":"noise","count":8}`)
	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status %d, want 422", rec.Code)
	}
}

func TestAPIApplyBadCount(t *testing.T) {
	s := newTestServer()
	rec := post(t, s, "/api/apply", `{"b":[1],"a":[1],"signal":"step","count":0}`)
	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status %d, want 422", rec.Code)
	}
}
