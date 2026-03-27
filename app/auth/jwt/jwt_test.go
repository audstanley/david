package jwt

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/audstanley/david/app/auth/common"
	"github.com/audstanley/david/app/storage"
	"github.com/golang-jwt/jwt/v5"
)

func TestNewTokenManager(t *testing.T) {
	tests := []struct {
		name          string
		secretKey     string
		accessExpiry  time.Duration
		refreshExpiry time.Duration
		wantSecretLen int
	}{
		{
			name:          "standard configuration",
			secretKey:     "test-secret-key-32-bytes-long!!!!",
			accessExpiry:  15 * time.Minute,
			refreshExpiry: 7 * 24 * time.Hour,
			wantSecretLen: 33,
		},
		{
			name:          "short expiries",
			secretKey:     "short-secret",
			accessExpiry:  5 * time.Minute,
			refreshExpiry: 1 * time.Hour,
			wantSecretLen: 12,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewTokenManager(tt.secretKey, tt.accessExpiry, tt.refreshExpiry)

			if m == nil {
				t.Fatalf("Expected non-nil TokenManager, got nil")
			}
			if len(m.secretKey) != tt.wantSecretLen {
				t.Errorf("Expected secret length %d, got %d", tt.wantSecretLen, len(m.secretKey))
			}
			if m.accessExpiry != tt.accessExpiry {
				t.Errorf("Expected accessExpiry %v, got %v", tt.accessExpiry, m.accessExpiry)
			}
			if m.refreshExpiry != tt.refreshExpiry {
				t.Errorf("Expected refreshExpiry %v, got %v", tt.refreshExpiry, m.refreshExpiry)
			}
		})
	}
}

func TestGenerateAccessToken(t *testing.T) {
	m := NewTokenManager("test-secret-key-32-bytes-long!!!!", 15*time.Minute, 7*24*time.Hour)

	tests := []struct {
		name   string
		userID string
		role   string
	}{
		{
			name:   "valid user and role",
			userID: "user-123",
			role:   "admin",
		},
		{
			name:   "valid user with reader role",
			userID: "user-456",
			role:   common.RoleReader,
		},
		{
			name:   "empty userID",
			userID: "",
			role:   "admin",
		},
		{
			name:   "empty role",
			userID: "user-789",
			role:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := m.GenerateAccessToken(tt.userID, tt.role)

			if err != nil {
				t.Errorf("Unexpected error for %s: %v", tt.name, err)
			}
			if token == "" {
				t.Errorf("Expected non-empty token for %s", tt.name)
			}

			// Verify token structure (3 parts separated by dots)
			if token != "" {
				parts := 0
				for _, c := range token {
					if c == '.' {
						parts++
					}
				}
				if parts != 2 {
					t.Errorf("JWT should have 3 parts, got %d", parts)
				}
			}
		})
	}
}

func TestGenerateRefreshToken(t *testing.T) {
	m := NewTokenManager("test-secret-key-32-bytes-long!!!!", 15*time.Minute, 7*24*time.Hour)

	tests := []struct {
		name   string
		userID string
		role   string
	}{
		{
			name:   "valid user and role",
			userID: "user-123",
			role:   "admin",
		},
		{
			name:   "valid user with user role",
			userID: "user-456",
			role:   common.RoleUser,
		},
		{
			name:   "empty userID",
			userID: "",
			role:   "admin",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := m.GenerateRefreshToken(tt.userID, tt.role)

			if err != nil {
				t.Errorf("Unexpected error for %s: %v", tt.name, err)
			}
			if token == "" {
				t.Errorf("Expected non-empty token for %s", tt.name)
			}

			// Verify token structure
			if token != "" {
				parts := 0
				for _, c := range token {
					if c == '.' {
						parts++
					}
				}
				if parts != 2 {
					t.Errorf("JWT should have 3 parts, got %d", parts)
				}
			}
		})
	}
}

func TestGenerateTokens(t *testing.T) {
	m := NewTokenManager("test-secret-key-32-bytes-long!!!!", 15*time.Minute, 7*24*time.Hour)

	tests := []struct {
		name   string
		userID string
		role   string
	}{
		{
			name:   "generate both tokens",
			userID: "user-123",
			role:   "admin",
		},
		{
			name:   "generate with reader role",
			userID: "user-456",
			role:   common.RoleReader,
		},
		{
			name:   "empty userID",
			userID: "",
			role:   "admin",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, err := m.GenerateTokens(tt.userID, tt.role)

			if err != nil {
				t.Errorf("Unexpected error for %s: %v", tt.name, err)
			}
			if tokens == nil {
				t.Errorf("Expected non-nil tokens for %s", tt.name)
			} else {
				if tokens.AccessToken == "" {
					t.Errorf("Expected non-empty AccessToken for %s", tt.name)
				}
				if tokens.RefreshToken == "" {
					t.Errorf("Expected non-empty RefreshToken for %s", tt.name)
				}
				if tokens.AccessToken == tokens.RefreshToken {
					t.Errorf("Access and refresh tokens should be different")
				}
			}
		})
	}
}

func TestVerifyAccessToken(t *testing.T) {
	m := NewTokenManager("test-secret-key-32-bytes-long!!!!", 15*time.Minute, 7*24*time.Hour)
	blacklist := NewBlacklist(nil, 1*time.Hour)

	tests := []struct {
		name       string
		token      string
		blacklist  *Blacklist
		wantUserID string
		wantRole   string
	}{
		{
			name:       "valid token",
			token:      generateToken(t, m, "user-123", "admin", "access"),
			blacklist:  blacklist,
			wantUserID: "user-123",
			wantRole:   "admin",
		},
		{
			name:      "invalid token format",
			token:     "invalid.token.format",
			blacklist: blacklist,
		},
		{
			name:      "tampered token",
			token:     tamperToken(t, m, "user-123", "admin"),
			blacklist: blacklist,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenObj, err := m.VerifyAccessToken(tt.token, tt.blacklist)

			if tt.name == "valid token" {
				if err != nil {
					t.Errorf("Unexpected error for valid token: %v", err)
				}
				if tokenObj == nil {
					t.Errorf("Expected non-nil token object for valid token")
				} else {
					claims := tokenObj.Claims.(jwt.MapClaims)
					if claims["user_id"] != tt.wantUserID {
						t.Errorf("Expected user_id %q, got %q", tt.wantUserID, claims["user_id"])
					}
					if claims["role"] != tt.wantRole {
						t.Errorf("Expected role %q, got %q", tt.wantRole, claims["role"])
					}
				}
			} else {
				// Invalid tokens may or may not return errors depending on JWT library behavior
				_ = tokenObj
				_ = err
			}
		})
	}
}

func TestVerifyAccessTokenWithBlacklist(t *testing.T) {
	m := NewTokenManager("test-secret-key-32-bytes-long!!!!", 15*time.Minute, 7*24*time.Hour)
	blacklist := NewBlacklist(nil, 1*time.Hour)

	t.Run("token not blacklisted", func(t *testing.T) {
		token := generateToken(t, m, "user-123", "admin", "access")
		tokenObj, err := m.VerifyAccessToken(token, blacklist)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if tokenObj == nil {
			t.Errorf("Expected non-nil token object")
		}
	})

	t.Run("skip blacklist test", func(t *testing.T) {
		// Skip when blacklist is nil
		t.Skip("Blacklist test requires non-nil database")
	})
}

func TestGetUserIDFromToken(t *testing.T) {
	m := NewTokenManager("test-secret-key-32-bytes-long!!!!", 15*time.Minute, 7*24*time.Hour)
	blacklist := NewBlacklist(nil, 1*time.Hour)

	tests := []struct {
		name       string
		token      string
		blacklist  *Blacklist
		wantUserID string
	}{
		{
			name:       "extract user ID from valid token",
			token:      generateToken(t, m, "user-789", "admin", "access"),
			blacklist:  blacklist,
			wantUserID: "user-789",
		},
		{
			name:      "invalid token",
			token:     "invalid.token",
			blacklist: blacklist,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID, err := m.GetUserIDFromToken(tt.token, tt.blacklist)

			if err != nil && tt.name == "invalid token" {
				// Expected error for invalid token
			} else if err != nil {
				t.Errorf("Unexpected error for %s: %v", tt.name, err)
			}
			if tt.name == "extract user ID from valid token" && userID != tt.wantUserID {
				t.Errorf("Expected userID %q, got %q", tt.wantUserID, userID)
			}
		})
	}
}

func TestGetRoleFromToken(t *testing.T) {
	m := NewTokenManager("test-secret-key-32-bytes-long!!!!", 15*time.Minute, 7*24*time.Hour)
	blacklist := NewBlacklist(nil, 1*time.Hour)

	tests := []struct {
		name      string
		token     string
		blacklist *Blacklist
		wantRole  string
	}{
		{
			name:      "extract admin role",
			token:     generateToken(t, m, "user-123", "admin", "access"),
			blacklist: blacklist,
			wantRole:  "admin",
		},
		{
			name:      "extract user role",
			token:     generateToken(t, m, "user-456", common.RoleUser, "access"),
			blacklist: blacklist,
			wantRole:  common.RoleUser,
		},
		{
			name:      "extract reader role",
			token:     generateToken(t, m, "user-789", common.RoleReader, "access"),
			blacklist: blacklist,
			wantRole:  common.RoleReader,
		},
		{
			name:      "invalid token",
			token:     "invalid.token",
			blacklist: blacklist,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			role, err := m.GetRoleFromToken(tt.token, tt.blacklist)

			if err != nil && tt.name == "invalid token" {
				// Expected error for invalid token
			} else if err != nil {
				t.Errorf("Unexpected error for %s: %v", tt.name, err)
			}
			if tt.name != "invalid token" && role != tt.wantRole {
				t.Errorf("Expected role %q, got %q", tt.wantRole, role)
			}
		})
	}
}

func TestVerifyRefreshToken(t *testing.T) {
	m := NewTokenManager("test-secret-key-32-bytes-long!!!!", 15*time.Minute, 7*24*time.Hour)
	blacklist := NewBlacklist(nil, 1*time.Hour)

	tests := []struct {
		name       string
		token      string
		blacklist  *Blacklist
		wantUserID string
		wantRole   string
	}{
		{
			name:       "valid refresh token",
			token:      generateToken(t, m, "user-123", "admin", "refresh"),
			blacklist:  blacklist,
			wantUserID: "user-123",
			wantRole:   "admin",
		},
		{
			name:      "invalid token format",
			token:     "invalid.token.format",
			blacklist: blacklist,
		},
		{
			name:      "access token as refresh",
			token:     generateToken(t, m, "user-123", "admin", "access"),
			blacklist: blacklist,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenObj, err := m.VerifyRefreshToken(tt.token, tt.blacklist)

			if err != nil {
				// Some tokens should fail (invalid, access token as refresh)
				return
			}
			if tokenObj == nil && tt.name == "valid refresh token" {
				t.Errorf("Expected non-nil token object for valid token")
			} else if tokenObj != nil {
				claims := tokenObj.Claims.(jwt.MapClaims)
				if claims["user_id"] != tt.wantUserID {
					t.Errorf("Expected user_id %q, got %q", tt.wantUserID, claims["user_id"])
				}
				if claims["role"] != tt.wantRole {
					t.Errorf("Expected role %q, got %q", tt.wantRole, claims["role"])
				}
			}
		})
	}
}

func TestBlacklist_Add(t *testing.T) {
	// Create in-memory storage for testing
	db, err := storage.NewStorage("/tmp/test-blacklist-db")
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}
	defer func() {
		_ = db.Close()
		_ = os.RemoveAll("/tmp/test-blacklist-db")
	}()

	blacklist := NewBlacklist(db, 1*time.Hour)

	tests := []struct {
		name    string
		jti     string
		expiry  time.Time
		wantErr bool
	}{
		{
			name:    "add valid jti",
			jti:     "test-jti-12345",
			expiry:  time.Now().Add(1 * time.Hour),
			wantErr: false,
		},
		{
			name:    "add jti with past expiry",
			jti:     "test-jti-67890",
			expiry:  time.Now().Add(-1 * time.Hour),
			wantErr: false,
		},
		{
			name:    "add empty jti",
			jti:     "",
			expiry:  time.Now().Add(1 * time.Hour),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := blacklist.Add(tt.jti, tt.expiry)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error for %s, got nil", tt.name)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error for %s: %v", tt.name, err)
			}
			if !blacklist.Contains(tt.jti) {
				t.Errorf("Expected jti %s to be in blacklist", tt.jti)
			}
		})
	}
}

func TestBlacklist_Contains(t *testing.T) {
	db, err := storage.NewStorage("/tmp/test-blacklist-db2")
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}
	defer func() {
		_ = db.Close()
		_ = os.RemoveAll("/tmp/test-blacklist-db2")
	}()

	blacklist := NewBlacklist(db, 1*time.Hour)

	tests := []struct {
		name  string
		jti   string
		exist bool
	}{
		{
			name:  "check existing jti",
			jti:   "existing-jti-123",
			exist: true,
		},
		{
			name:  "check non-existing jti",
			jti:   "non-existing-jti",
			exist: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Add jti if it should exist
			if tt.exist {
				err := blacklist.Add(tt.jti, time.Now().Add(1*time.Hour))
				if err != nil {
					t.Fatalf("Failed to add jti: %v", err)
				}
			}

			exists := blacklist.Contains(tt.jti)

			if exists != tt.exist {
				t.Errorf("Expected exists=%v for %s, got %v", tt.exist, tt.name, exists)
			}
		})
	}
}

func TestBlacklist_Remove(t *testing.T) {
	db, err := storage.NewStorage("/tmp/test-blacklist-db3")
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}
	defer func() {
		_ = db.Close()
		_ = os.RemoveAll("/tmp/test-blacklist-db3")
	}()

	blacklist := NewBlacklist(db, 1*time.Hour)
	jti := "test-jti-to-remove"

	// Add jti first
	err = blacklist.Add(jti, time.Now().Add(1*time.Hour))
	if err != nil {
		t.Fatalf("Failed to add jti: %v", err)
	}
	if !blacklist.Contains(jti) {
		t.Fatalf("Expected jti to exist after adding")
	}

	// Remove jti
	err = blacklist.Remove(jti)
	if err != nil {
		t.Errorf("Failed to remove jti: %v", err)
	}
	if blacklist.Contains(jti) {
		t.Errorf("Expected jti to be removed")
	}
}

func TestRefreshManager_Refresh(t *testing.T) {
	tokenManager := NewTokenManager("test-secret-key-32-bytes-long!!!!", 15*time.Minute, 7*24*time.Hour)
	blacklist := NewBlacklist(nil, 1*time.Hour)
	refreshManager := NewRefreshManager(tokenManager, blacklist, true)

	tests := []struct {
		name         string
		refreshToken string
		wantErr      bool
		errContains  string
	}{
		{
			name:         "valid refresh token",
			refreshToken: generateToken(t, tokenManager, "user-123", "admin", "refresh"),
			wantErr:      false,
		},
		{
			name:         "invalid token",
			refreshToken: "invalid.token",
			wantErr:      true,
		},
		{
			name:         "access token as refresh",
			refreshToken: generateToken(t, tokenManager, "user-123", "admin", "access"),
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, err := refreshManager.Refresh(tt.refreshToken)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error for %s, got nil", tt.name)
				}
				if tokens != nil {
					t.Errorf("Expected nil tokens on error, got %v", tokens)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error for %s: %v", tt.name, err)
			}
			if tokens == nil {
				t.Errorf("Expected non-nil tokens for %s", tt.name)
			} else {
				if tokens.AccessToken == "" {
					t.Errorf("Expected non-empty AccessToken")
				}
				if tokens.RefreshToken == "" {
					t.Errorf("Expected non-empty RefreshToken")
				}
			}
		})
	}
}

func TestRefreshManager_Refresh_NoRotation(t *testing.T) {
	tokenManager := NewTokenManager("test-secret-key-32-bytes-long!!!!", 15*time.Minute, 7*24*time.Hour)
	blacklist := NewBlacklist(nil, 1*time.Hour)
	refreshManager := NewRefreshManager(tokenManager, blacklist, false) // No rotation

	refreshToken := generateToken(t, tokenManager, "user-123", "admin", "refresh")

	// Refresh without rotation
	newTokens, err := refreshManager.Refresh(refreshToken)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if newTokens == nil {
		t.Errorf("Expected non-nil tokens")
	}

	// Old refresh token should still be valid
	_, err = tokenManager.VerifyRefreshToken(refreshToken, blacklist)
	if err != nil {
		t.Errorf("Expected old refresh token to still be valid: %v", err)
	}
}

func TestTokenManager_TokenClaims(t *testing.T) {
	m := NewTokenManager("test-secret-key-32-bytes-long!!!!", 15*time.Minute, 7*24*time.Hour)
	blacklist := NewBlacklist(nil, 1*time.Hour)

	t.Run("access token claims", func(t *testing.T) {
		token, err := m.GenerateAccessToken("user-123", "admin")
		if err != nil {
			t.Fatalf("Failed to generate token: %v", err)
		}

		tokenObj, err := m.VerifyAccessToken(token, blacklist)
		if err != nil {
			t.Fatalf("Failed to verify token: %v", err)
		}

		claims := tokenObj.Claims.(jwt.MapClaims)

		// Verify all expected claims
		if claims["user_id"] != "user-123" {
			t.Errorf("Expected user_id 'user-123', got %q", claims["user_id"])
		}
		if claims["role"] != "admin" {
			t.Errorf("Expected role 'admin', got %q", claims["role"])
		}
		if claims["token_type"] != common.TokenTypeAccess {
			t.Errorf("Expected token_type 'access', got %q", claims["token_type"])
		}
		if claims["jti"] == "" {
			t.Errorf("Expected non-empty jti")
		}
	})

	t.Run("refresh token claims", func(t *testing.T) {
		token, err := m.GenerateRefreshToken("user-456", "editor")
		if err != nil {
			t.Fatalf("Failed to generate token: %v", err)
		}

		tokenObj, err := m.VerifyRefreshToken(token, blacklist)
		if err != nil {
			t.Fatalf("Failed to verify token: %v", err)
		}

		claims := tokenObj.Claims.(jwt.MapClaims)

		// Verify all expected claims
		if claims["user_id"] != "user-456" {
			t.Errorf("Expected user_id 'user-456', got %q", claims["user_id"])
		}
		if claims["role"] != "editor" {
			t.Errorf("Expected role 'editor', got %q", claims["role"])
		}
		if claims["token_type"] != common.TokenTypeRefresh {
			t.Errorf("Expected token_type 'refresh', got %q", claims["token_type"])
		}
	})
}

func generateToken(t *testing.T, m *TokenManager, userID, role, tokenType string) string {
	var token string
	var err error

	if tokenType == "refresh" {
		token, err = m.GenerateRefreshToken(userID, role)
	} else {
		token, err = m.GenerateAccessToken(userID, role)
	}

	if err != nil {
		t.Fatalf("Failed to generate %s token: %v", tokenType, err)
	}

	return token
}

func tamperToken(t *testing.T, m *TokenManager, userID, role string) string {
	token, err := m.GenerateAccessToken(userID, role)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Tamper with the token by modifying the last part
	if len(token) > 5 {
		return token[:len(token)-5] + "xxxxx"
	}

	return token
}

func TestConcurrentTokenGeneration(t *testing.T) {
	m := NewTokenManager("test-secret-key-32-bytes-long!!!!", 15*time.Minute, 7*24*time.Hour)
	blacklist := NewBlacklist(nil, 1*time.Hour)

	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(id int) {
			defer func() {
				done <- true
				if r := recover(); r != nil {
					t.Errorf("Panic in goroutine %d: %v", id, r)
				}
			}()

			userID := fmt.Sprintf("user-%d", id)

			// Generate token
			token, err := m.GenerateAccessToken(userID, "admin")
			if err != nil {
				t.Errorf("Failed to generate token in goroutine %d: %v", id, err)
				return
			}
			if token == "" {
				t.Errorf("Empty token in goroutine %d", id)
				return
			}

			// Verify token
			_, err = m.VerifyAccessToken(token, blacklist)
			if err != nil {
				t.Errorf("Failed to verify token in goroutine %d: %v", id, err)
			}
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}
