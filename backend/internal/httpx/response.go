package httpx

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id,omitempty"`
}

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, r *http.Request, status int, code string, message string) {
	rid, _ := RequestIDFromContext(r.Context())
	WriteJSON(w, status, ErrorResponse{Code: code, Message: message, RequestID: rid})
}
