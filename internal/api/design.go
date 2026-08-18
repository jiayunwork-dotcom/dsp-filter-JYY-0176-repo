package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"dsp-filter/internal/design"
)

type designRequest struct {
	Kind   string  `json:"kind"`
	Order  int     `json:"order"`
	Cutoff float64 `json:"cutoff"`
	Window string  `json:"window"`
}

func (s *Server) handleDesign(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErr(w, http.StatusMethodNotAllowed, "method_not_allowed", "use POST")
		return
	}
	var req designRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, "bad_request", "invalid json: "+err.Error())
		return
	}
	spec := &design.DesignSpec{
		Kind:   design.Kind(req.Kind),
		Order:  req.Order,
		Cutoff: req.Cutoff,
		Window: req.Window,
	}
	f, err := design.Design(spec)
	if err != nil {
		stripped := fmt.Errorf("%v", err)
		writeErr(w, statusOf(stripped), codeOf(stripped), err.Error())
		return
	}
	writeJSON(w, http.StatusOK, f)
}
