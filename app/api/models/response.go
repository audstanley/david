package models

import "encoding/json"

// SuccessResponse represents a successful API response
type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// ErrorResponse represents an error API response
type ErrorResponse struct {
	Success bool      `json:"success"`
	Error   *APIError `json:"error"`
}

// APIError represents an API error
type APIError struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// Meta contains pagination metadata
type Meta struct {
	Total   int  `json:"total"`
	Page    int  `json:"page"`
	PerPage int  `json:"per_page"`
	HasMore bool `json:"has_more"`
}

// NewSuccessResponse creates a success response
func NewSuccessResponse(data interface{}, meta *Meta) *SuccessResponse {
	return &SuccessResponse{
		Success: true,
		Data:    data,
		Meta:    meta,
	}
}

// NewErrorResponse creates an error response
func NewErrorResponse(code, message string, details interface{}) *ErrorResponse {
	return &ErrorResponse{
		Success: false,
		Error: &APIError{
			Code:    code,
			Message: message,
			Details: details,
		},
	}
}

// ToJSON serializes response to JSON
func (r *SuccessResponse) ToJSON() ([]byte, error) {
	return json.Marshal(r)
}

// ToJSON serializes error response to JSON
func (r *ErrorResponse) ToJSON() ([]byte, error) {
	return json.Marshal(r)
}
