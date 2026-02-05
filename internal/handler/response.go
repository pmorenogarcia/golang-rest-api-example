package handler

import (
	"encoding/json"
	"net/http"

	"github.com/polgarcia/golang-rest-api/pkg/logger"
	"go.uber.org/zap"
)

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error     string `json:"error"`
	Message   string `json:"message"`
	Code      int    `json:"code"`
	RequestID string `json:"request_id,omitempty"`
}

// WriteJSON writes a JSON response
func WriteJSON(w http.ResponseWriter, status int, data interface{}, log *logger.Logger) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Error("Failed to encode JSON response", zap.Error(err))
	}
}

// WriteError writes an error response
func WriteError(w http.ResponseWriter, status int, message string, log *logger.Logger) {
	errResp := ErrorResponse{
		Error:   http.StatusText(status),
		Message: message,
		Code:    status,
	}

	WriteJSON(w, status, errResp, log)
}
