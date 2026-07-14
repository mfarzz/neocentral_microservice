package response

import (
	"encoding/json"
	"errors"
	"net/http"

	"neocentral-go/pkg/apperror"
)

// APIResponse is the standard envelope for all JSON responses.
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

// JSON writes a success response.
func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(APIResponse{
		Success: true,
		Data:    data,
	})
}

// Success writes a success response with a message.
func Success(w http.ResponseWriter, status int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Error writes an error response. If the error is an *apperror.AppError its
// status code is used; otherwise the provided fallback status is used.
func Error(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(APIResponse{
		Success: false,
		Message: msg,
	})
}

// FromError inspects err for an *apperror.AppError and responds accordingly.
func FromError(w http.ResponseWriter, err error) {
	var appErr *apperror.AppError
	if errors.As(err, &appErr) {
		Error(w, appErr.Code, appErr.Message)
		return
	}
	Error(w, http.StatusInternalServerError, "Internal server error")
}

// ValidationError writes a 422 response with field-level validation errors.
func ValidationError(w http.ResponseWriter, errs interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnprocessableEntity)
	json.NewEncoder(w).Encode(APIResponse{
		Success: false,
		Message: "Validation failed",
		Errors:  errs,
	})
}
