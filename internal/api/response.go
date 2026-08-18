package api

import (
	"encoding/json"
	"net/http"

	"dsp-filter/internal/response"
)

type responseRequest struct {
	B      []float64 `json:"b"`
	A      []float64 `json:"a"`
	Points int       `json:"points"`
	Freq   []float64 `json:"freq"`
}

func (s *Server) handleResponse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErr(w, http.StatusMethodNotAllowed, "method_not_allowed", "use POST")
		return
	}
	var req responseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, "bad_request", "invalid json: "+err.Error())
		return
	}
	freq := req.Freq
	if len(freq) == 0 {
		if req.Points < 2 {
			writeErr(w, http.StatusUnprocessableEntity, "empty_frequency_grid", "points must be >= 2")
			return
		}
		freq = response.Grid(req.Points)
	}
	res, err := response.Compute(req.B, req.A, freq)
	if err != nil {
		writeErr(w, statusOf(err), codeOf(err), err.Error())
		return
	}
	writeJSON(w, http.StatusOK, res)
}
