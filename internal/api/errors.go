package api

import (
	"encoding/json"
	"net/http"

	"dsp-filter/internal/design"
	"dsp-filter/internal/response"
)

type ErrResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, ErrResponse{Code: code, Message: message})
}

func statusOf(err error) int {
	switch err.(type) {
	case *design.Error:
		return http.StatusUnprocessableEntity
	case *response.Error:
		return http.StatusUnprocessableEntity
	}
	return http.StatusBadRequest
}

func codeOf(err error) string {
	if e, ok := err.(*design.Error); ok {
		return e.Code
	}
	if e, ok := err.(*response.Error); ok {
		return e.Code
	}
	return "internal_error"
}
