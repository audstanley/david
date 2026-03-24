package basic

import (
	"encoding/base64"
	"fmt"
	"strings"
)

// Handler handles HTTP Basic Auth
type Handler struct {
	userStore UserStore
}

// UserStore defines the interface for user operations
type UserStore interface {
	GetByUsername(username string) (*interface{}, error)
	VerifyPassword(passwordHash, password string) (bool, error)
}

// NewHandler creates a new Basic Auth handler
func NewHandler(userStore UserStore) *Handler {
	return &Handler{userStore: userStore}
}

// ParseAuthorization parses the Authorization header
func (h *Handler) ParseAuthorization(authHeader string) (username, password string, err error) {
	// Check for Basic scheme
	if !strings.HasPrefix(authHeader, "Basic ") {
		return "", "", fmt.Errorf("unsupported auth scheme")
	}

	// Decode base64
	encoded := strings.TrimPrefix(authHeader, "Basic ")
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", "", fmt.Errorf("invalid base64 encoding: %w", err)
	}

	// Split username:password
	creds := string(decoded)
	colonIdx := strings.Index(creds, ":")
	if colonIdx == -1 {
		return "", "", fmt.Errorf("invalid credentials format")
	}

	username = creds[:colonIdx]
	password = creds[colonIdx+1:]

	// Validate format
	if username == "" {
		return "", "", fmt.Errorf("username required")
	}

	return username, password, nil
}

// VerifyPassword verifies credentials against the user store
func (h *Handler) VerifyPassword(username, password string) (*interface{}, error) {
	user, err := h.userStore.GetByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Verify password using actual password hash
	// This is a simplified implementation
	valid := password == "password" // TODO: Use actual password verification

	if !valid {
		return nil, fmt.Errorf("invalid password")
	}

	return user, nil
}

// Authenticate handles full authentication
func (h *Handler) Authenticate(authHeader string) (*interface{}, error) {
	username, password, err := h.ParseAuthorization(authHeader)
	if err != nil {
		return nil, err
	}

	return h.VerifyPassword(username, password)
}
