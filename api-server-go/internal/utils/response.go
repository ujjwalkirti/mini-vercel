package utils

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func RespondJSON(w http.ResponseWriter, statusCode int, response APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func Success(w http.ResponseWriter, data interface{}, message ...string) {
	msg := ""
	if len(message) > 0 {
		msg = message[0]
	}

	RespondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Message: msg,
		Data:    data,
	})
}

func Created(w http.ResponseWriter, data interface{}, message string) {
	RespondJSON(w, http.StatusCreated, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Error(w http.ResponseWriter, statusCode int, message string, errorDetail string) {
	RespondJSON(w, statusCode, APIResponse{
		Success: false,
		Message: message,
		Error:   errorDetail,
	})
}

func NotFound(w http.ResponseWriter, message string) {
	Error(w, http.StatusNotFound, message, "Resource not found")
}

func Unauthorized(w http.ResponseWriter, message string) {
	Error(w, http.StatusUnauthorized, message, "Unauthorized access")
}

func Forbidden(w http.ResponseWriter, message string) {
	Error(w, http.StatusForbidden, message, "Access forbidden")
}

func BadRequest(w http.ResponseWriter, message string) {
	Error(w, http.StatusBadRequest, message, "Invalid request")
}

func InternalServerError(w http.ResponseWriter, message string) {
	Error(w, http.StatusInternalServerError, message, "Internal server error")
}
