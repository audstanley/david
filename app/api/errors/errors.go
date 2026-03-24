package errors

import (
	"encoding/json"
	"net/http"
)

// API Error codes
const (
	CodeNotFound          = "not_found"
	CodeUnauthorized      = "unauthorized"
	CodeForbidden         = "forbidden"
	CodeValidationError   = "validation_error"
	CodeConflict          = "conflict"
	CodeRateLimitExceeded = "rate_limit_exceeded"
	CodeInternalError     = "internal_error"
	CodeInvalidAuth       = "invalid_auth"
	CodeAuthExpired       = "auth_expired"
)

// HTTP status codes
var ErrorToStatus = map[string]int{
	CodeNotFound:          http.StatusNotFound,
	CodeUnauthorized:      http.StatusUnauthorized,
	CodeForbidden:         http.StatusForbidden,
	CodeValidationError:   http.StatusBadRequest,
	CodeConflict:          http.StatusConflict,
	CodeRateLimitExceeded: http.StatusTooManyRequests,
	CodeInternalError:     http.StatusInternalServerError,
	CodeInvalidAuth:       http.StatusUnauthorized,
	CodeAuthExpired:       http.StatusUnauthorized,
}

// Error represents an API error
type Error struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// NewError creates a new API error
func NewError(code, message string, details interface{}) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// WithDetails adds details to an error
func (e *Error) WithDetails(details interface{}) *Error {
	e.Details = details
	return e
}

// Write writes the error to the HTTP response
func (e *Error) Write(w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := map[string]interface{}{
		"success": false,
		"error": map[string]interface{}{
			"code":    e.Code,
			"message": e.Message,
		},
	}

	if e.Details != nil {
		response["error"].(map[string]interface{})["details"] = e.Details
	}

	json.NewEncoder(w).Encode(response)
}

// NotFound returns a not found error
func NotFound(resource, identifier string) *Error {
	return NewError(CodeNotFound, resource+" "+identifier+" not found", nil)
}

// Unauthorized returns an unauthorized error
func Unauthorized(message string) *Error {
	if message == "" {
		message = "authentication required"
	}
	return NewError(CodeUnauthorized, message, nil)
}

// Forbidden returns a forbidden error
func Forbidden(message string) *Error {
	if message == "" {
		message = "access denied"
	}
	return NewError(CodeForbidden, message, nil)
}

// ValidationError returns a validation error
func ValidationError(field, message string) *Error {
	return NewError(CodeValidationError, "validation failed", map[string]string{"field": field, "message": message})
}

// ConflictError returns a conflict error
func Conflict(message string) *Error {
	return NewError(CodeConflict, message, nil)
}

// RateLimitError returns a rate limit error
func RateLimitError(message string) *Error {
	return NewError(CodeRateLimitExceeded, message, nil)
}

// InternalError returns an internal error
func InternalError(err error) *Error {
	return NewError(CodeInternalError, "internal server error", map[string]string{"original_error": err.Error()})
}
