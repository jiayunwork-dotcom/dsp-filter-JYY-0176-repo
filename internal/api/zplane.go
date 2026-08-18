package api

import (
	"encoding/json"
	"net/http"

	"dsp-filter/internal/zplane"
)

type zplaneRequest struct {
	B []float64 `json:"b"`
	A []float64 `json:"a"`
}

func (s *Server) handleZPlane(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErr(w, http.StatusMethodNotAllowed, "method_not_allowed", "use POST")
		return
	}
	var req zplaneRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, "bad_request", "invalid json: "+err.Error())
		return
	}
	if len(req.B) == 0 || len(req.A) == 0 {
		writeErr(w, http.StatusUnprocessableEntity, "bad_coefficients", "b and a must be non-empty")
		return
	}
	zp := zplane.ZeroPoles(req.B, req.A)
	writeJSON(w, http.StatusOK, zp)
}
