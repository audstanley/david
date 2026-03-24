package handlers

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/audstanley/david/app/api/errors"
	"github.com/audstanley/david/app/api/models"
	"github.com/audstanley/david/app/auth/common"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	tokenSecret string
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(tokenSecret string) *AuthHandler {
	return &AuthHandler{tokenSecret: tokenSecret}
}

// LoginHandler handles login requests
func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req models.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.ValidationError("body", "invalid JSON").Write(w, http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		errors.ValidationError("request", err.Error()).Write(w, http.StatusBadRequest)
		return
	}

	// TODO: Actually verify credentials from database
	// For now, use hardcoded test credentials
	if req.Username != "admin" || req.Password != "password" {
		errors.Unauthorized("invalid credentials").Write(w, http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"access_token":  generateMockToken(req.Username, common.RoleAdmin),
		"refresh_token": generateMockRefreshToken(req.Username),
		"token_type":    "Bearer",
		"expires_in":    86400,
		"user": map[string]interface{}{
			"id":       "user-1",
			"username": req.Username,
			"role":     common.RoleAdmin,
		},
	})
}

// RefreshHandler handles refresh token requests
func (h *AuthHandler) RefreshHandler(w http.ResponseWriter, r *http.Request) {
	var req models.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.ValidationError("body", "invalid JSON").Write(w, http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		errors.ValidationError("request", err.Error()).Write(w, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"access_token": generateMockToken("user-1", common.RoleUser),
	})
}

// LogoutHandler handles logout requests
func (h *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
}

// VerifyHandler handles token verification
func (h *AuthHandler) VerifyHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		errors.Unauthorized("no token provided").Write(w, http.StatusUnauthorized)
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	if strings.HasPrefix(token, "mock_") {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"valid":   true,
			"user_id": "user-1",
			"role":    common.RoleUser,
		})
		return
	}

	errors.Unauthorized("invalid token").Write(w, http.StatusUnauthorized)
}

// CreateAPIKeyHandler handles API key creation
func (h *AuthHandler) CreateAPIKeyHandler(w http.ResponseWriter, r *http.Request) {
	var req models.CreateAPIKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.ValidationError("body", "invalid JSON").Write(w, http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		errors.ValidationError("request", err.Error()).Write(w, http.StatusBadRequest)
		return
	}

	key := "david_" + time.Now().Format("20060102150405")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":          time.Now().Unix(),
		"name":        req.Name,
		"key":         key,
		"permissions": req.Permissions,
		"expiry":      time.Now().Add(time.Duration(req.ExpiryDays) * 24 * time.Hour).Format(time.RFC3339),
	})
}

// ListAPIKeyHandler handles listing API keys
func (h *AuthHandler) ListAPIKeyHandler(w http.ResponseWriter, r *http.Request) {
	keys := []map[string]interface{}{
		{
			"id":          1,
			"name":        "Main Key",
			"permissions": []string{"calendar.read", "event.read"},
			"created":     time.Now().Format(time.RFC3339),
			"disabled":    false,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(keys)
}

// RevokeAPIKeyHandler handles API key revocation
func (h *AuthHandler) RevokeAPIKeyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
}

// ParseBasicAuth parses HTTP Basic Auth header
func ParseBasicAuth(auth string) (username, password string, ok bool) {
	if !strings.HasPrefix(auth, "Basic ") {
		return "", "", false
	}

	encoded := strings.TrimPrefix(auth, "Basic ")
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", "", false
	}

	creds := string(decoded)
	colonIdx := strings.Index(creds, ":")
	if colonIdx == -1 {
		return "", "", false
	}

	return creds[:colonIdx], creds[colonIdx+1:], true
}

// Helper functions

func generateMockToken(username, role string) string {
	return "mock_" + username + "_" + role + "_" + time.Now().Format("20060102150405")
}

func generateMockRefreshToken(username string) string {
	return "mock_refresh_" + username + "_" + time.Now().Format("20060102150405")
}
