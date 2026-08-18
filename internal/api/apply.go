package api

import (
	"encoding/json"
	"net/http"

	"dsp-filter/internal/filter"
)

type applyRequest struct {
	B      []float64 `json:"b"`
	A      []float64 `json:"a"`
	Signal string    `json:"signal"`
	Count  int       `json:"count"`
	Freq   float64   `json:"freq"`
}

type applyResponse struct {
	Signal []float64 `json:"signal"`
	Output []float64 `json:"output"`
}

func (s *Server) handleApply(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErr(w, http.StatusMethodNotAllowed, "method_not_allowed", "use POST")
		return
	}
	var req applyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, "bad_request", "invalid json: "+err.Error())
		return
	}
	if len(req.B) == 0 || len(req.A) == 0 {
		writeErr(w, http.StatusUnprocessableEntity, "bad_coefficients", "b and a must be non-empty")
		return
	}
	if req.Count <= 0 || req.Count > 10000 {
		writeErr(w, http.StatusUnprocessableEntity, "bad_count", "count must be in (0, 10000]")
		return
	}
	var sig []float64
	switch req.Signal {
	case "step":
		sig = make([]float64, req.Count)
		for i := range sig {
			sig[i] = 1
		}
	case "impulse":
		sig = make([]float64, req.Count)
		sig[0] = 1
	case "sine":
		sig = filter.GenSinusoid(req.Count, req.Freq)
	default:
		writeErr(w, http.StatusUnprocessableEntity, "bad_signal", "signal must be step, impulse or sine")
		return
	}
	f := filter.New(req.B, req.A)
	out := f.Process(sig)
	writeJSON(w, http.StatusOK, applyResponse{Signal: sig, Output: out})
}
