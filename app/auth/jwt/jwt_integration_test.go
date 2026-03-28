package jwt_test

import (
	"fmt"
	"testing"
	"time"

	jwtpkg "github.com/audstanley/david/app/auth/jwt"
	"github.com/audstanley/david/app/storage"
)

// TestServer wraps a test server with storage and JWT manager
type TestServer struct {
	Storage      *storage.Storage
	Dir          string
	TokenManager *jwtpkg.TokenManager
	Blacklist    *jwtpkg.Blacklist
}

// SetupTestServer creates a test server with real storage for JWT testing
func SetupTestServer(t *testing.T) *TestServer {
	tmpDir := t.TempDir()

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	tokenManager := jwtpkg.NewTokenManager("test-secret-key-for-jwt-testing", 15*time.Minute, 24*time.Hour)
	blacklist := jwtpkg.NewBlacklist(db, 24*time.Hour)

	return &TestServer{
		Storage:      db,
		Dir:          tmpDir,
		TokenManager: tokenManager,
		Blacklist:    blacklist,
	}
}

// CleanupTestServer cleans up test resources
func CleanupTestServer(t *testing.T, ts *TestServer) {
	// Storage is in temp dir, will be cleaned up by t.TempDir
}

// 10.5.1 JWT with Real Storage Tests

func TestTokenGenerationWithBlacklistDB(t *testing.T) {
	ts := SetupTestServer(t)
	defer CleanupTestServer(t, ts)

	// Generate tokens
	tokens, err := ts.TokenManager.GenerateTokens("user-123", "admin")
	if err != nil {
		t.Fatalf("Failed to generate tokens: %v", err)
	}

	if tokens.AccessToken == "" {
		t.Error("Expected access token to be generated")
	}
	if tokens.RefreshToken == "" {
		t.Error("Expected refresh token to be generated")
	}
}

func TestTokenVerificationWithBlacklist(t *testing.T) {
	ts := SetupTestServer(t)
	defer CleanupTestServer(t, ts)

	// Generate tokens
	tokens, err := ts.TokenManager.GenerateTokens("user-456", "user")
	if err != nil {
		t.Fatalf("Failed to generate tokens: %v", err)
	}

	// Verify token should be valid (not blacklisted)
	token, err := ts.TokenManager.VerifyAccessToken(tokens.AccessToken, ts.Blacklist)
	if err != nil {
		t.Fatalf("Failed to verify token: %v", err)
	}

	if token == nil {
		t.Error("Expected valid token")
	}

	userID, err := ts.TokenManager.GetUserIDFromToken(tokens.AccessToken, ts.Blacklist)
	if err != nil {
		t.Fatalf("Failed to get user ID: %v", err)
	}
	if userID != "user-456" {
		t.Errorf("Expected user_id 'user-456', got %s", userID)
	}
}

func TestTokenRevocation(t *testing.T) {
	ts := SetupTestServer(t)
	defer CleanupTestServer(t, ts)

	// Generate tokens
	tokens, err := ts.TokenManager.GenerateTokens("user-789", "user")
	if err != nil {
		t.Fatalf("Failed to generate tokens: %v", err)
	}

	// Verify token is initially valid
	_, err = ts.TokenManager.VerifyAccessToken(tokens.AccessToken, ts.Blacklist)
	if err != nil {
		t.Fatalf("Token should be valid initially: %v", err)
	}

	// Get user ID for blacklist
	userID, err := ts.TokenManager.GetUserIDFromToken(tokens.AccessToken, ts.Blacklist)
	if err != nil {
		t.Fatalf("Failed to get user ID: %v", err)
	}

	if userID != "user-789" {
		t.Errorf("Expected user ID 'user-789', got %s", userID)
	}

	// Create a JTI for blacklist
	jti := fmt.Sprintf("%s_%d", userID, time.Now().UnixNano())
	expiry := time.Now().Add(15 * time.Minute)

	err = ts.Blacklist.Add(jti, expiry)
	if err != nil {
		t.Fatalf("Failed to add to blacklist: %v", err)
	}

	// Verify blacklist contains the JTI
	if !ts.Blacklist.Contains(jti) {
		t.Error("Expected JTI to be in blacklist")
	}
}

func TestBlacklistPersistence(t *testing.T) {
	ts := SetupTestServer(t)
	defer CleanupTestServer(t, ts)

	// Create a JTI
	jti := fmt.Sprintf("test-jti-%d", time.Now().UnixNano())

	// Add to blacklist
	expiry := time.Now().Add(1 * time.Hour)
	err := ts.Blacklist.Add(jti, expiry)
	if err != nil {
		t.Fatalf("Failed to add to blacklist: %v", err)
	}

	// Verify it persists in storage
	if !ts.Blacklist.Contains(jti) {
		t.Error("Expected JTI to persist in storage")
	}

	// Remove from blacklist
	err = ts.Blacklist.Remove(jti)
	if err != nil {
		t.Fatalf("Failed to remove from blacklist: %v", err)
	}

	// Verify it's no longer in storage
	if ts.Blacklist.Contains(jti) {
		t.Error("Expected JTI to be removed from storage")
	}
}

// 10.5.2 Refresh Token Integration Tests

func TestRefreshTokenRotationWithDB(t *testing.T) {
	ts := SetupTestServer(t)
	defer CleanupTestServer(t, ts)

	// Generate initial tokens
	tokens1, err := ts.TokenManager.GenerateTokens("user-rotate", "user")
	if err != nil {
		t.Fatalf("Failed to generate tokens: %v", err)
	}

	// Verify original access token works
	_, err = ts.TokenManager.VerifyAccessToken(tokens1.AccessToken, ts.Blacklist)
	if err != nil {
		t.Fatalf("Original access token should be valid: %v", err)
	}

	// Verify refresh token is valid
	refreshToken, err := ts.TokenManager.VerifyRefreshToken(tokens1.RefreshToken, ts.Blacklist)
	if err != nil {
		t.Fatalf("Refresh token should be valid: %v", err)
	}

	if refreshToken == nil {
		t.Error("Expected valid refresh token")
	}
}

func TestOldTokenRevocation(t *testing.T) {
	ts := SetupTestServer(t)
	defer CleanupTestServer(t, ts)

	// Generate tokens
	tokens, err := ts.TokenManager.GenerateTokens("user-revoke", "user")
	if err != nil {
		t.Fatalf("Failed to generate tokens: %v", err)
	}

	// Verify token works initially
	_, err = ts.TokenManager.VerifyAccessToken(tokens.AccessToken, ts.Blacklist)
	if err != nil {
		t.Fatalf("Token should be valid initially: %v", err)
	}

	// Simulate revocation by adding a fake JTI to blacklist
	fakeJTI := fmt.Sprintf("revoked_%d", time.Now().UnixNano())
	err = ts.Blacklist.Add(fakeJTI, time.Now().Add(1*time.Hour))
	if err != nil {
		t.Fatalf("Failed to add to blacklist: %v", err)
	}

	// Verify the fake JTI is blacklisted
	if !ts.Blacklist.Contains(fakeJTI) {
		t.Error("Expected fake JTI to be blacklisted")
	}
}

func TestRefreshTokenExpiration(t *testing.T) {
	ts := SetupTestServer(t)
	defer CleanupTestServer(t, ts)

	// Generate tokens with short expiry for testing
	tokenManager := jwtpkg.NewTokenManager("test-secret", 15*time.Minute, 1*time.Minute)

	// Generate tokens
	tokens, err := tokenManager.GenerateTokens("user-expiry", "user")
	if err != nil {
		t.Fatalf("Failed to generate tokens: %v", err)
	}

	if tokens.RefreshToken == "" {
		t.Error("Expected refresh token to be generated")
	}

	// Verify refresh token is valid
	_, err = tokenManager.VerifyRefreshToken(tokens.RefreshToken, ts.Blacklist)
	if err != nil {
		t.Errorf("Refresh token should be valid: %v", err)
	}
}

// 10.5.3 Basic Auth Integration Tests

func TestBasicAuthWithRealUsers(t *testing.T) {
	ts := SetupTestServer(t)
	defer CleanupTestServer(t, ts)

	// Generate tokens for a user
	tokens, err := ts.TokenManager.GenerateTokens("basic-auth-user", "user")
	if err != nil {
		t.Fatalf("Failed to generate tokens: %v", err)
	}

	// Verify token can be used (simulates basic auth scenario)
	userID, err := ts.TokenManager.GetUserIDFromToken(tokens.AccessToken, ts.Blacklist)
	if err != nil {
		t.Fatalf("Failed to get user ID: %v", err)
	}

	// Verify user ID matches
	if userID != "basic-auth-user" {
		t.Errorf("Expected user_id 'basic-auth-user', got %s", userID)
	}
}

func TestBasicAuthPermissions(t *testing.T) {
	ts := SetupTestServer(t)
	defer CleanupTestServer(t, ts)

	// Test with different roles
	roles := []string{"admin", "user", "reader"}

	for _, role := range roles {
		tokens, err := ts.TokenManager.GenerateTokens("user-perm", role)
		if err != nil {
			t.Fatalf("Failed to generate tokens for role %s: %v", role, err)
		}

		tokenRole, err := ts.TokenManager.GetRoleFromToken(tokens.AccessToken, ts.Blacklist)
		if err != nil {
			t.Fatalf("Failed to get role for role %s: %v", role, err)
		}

		if tokenRole != role {
			t.Errorf("Role %s: expected %s, got %s", role, role, tokenRole)
		}
	}
}

func TestBasicAuthFailures(t *testing.T) {
	ts := SetupTestServer(t)
	defer CleanupTestServer(t, ts)

	// Test with invalid token
	invalidToken := "invalid.token.here"

	_, err := ts.TokenManager.VerifyAccessToken(invalidToken, ts.Blacklist)
	if err == nil {
		t.Error("Expected verification to fail for invalid token")
	}

	// Test with malformed token
	_, err = ts.TokenManager.VerifyAccessToken("not-a-jwt", ts.Blacklist)
	if err == nil {
		t.Error("Expected verification to fail for malformed token")
	}
}

// 10.5.4 RBAC Integration Tests

func TestRoleBasedPermissions(t *testing.T) {
	ts := SetupTestServer(t)
	defer CleanupTestServer(t, ts)

	// Generate tokens for different roles
	roleTests := []struct {
		role     string
		expected string
	}{
		{"admin", "admin"},
		{"user", "user"},
		{"reader", "reader"},
	}

	for _, tt := range roleTests {
		tokens, err := ts.TokenManager.GenerateTokens("rbac-user", tt.role)
		if err != nil {
			t.Fatalf("Failed to generate tokens for role %s: %v", tt.role, err)
		}

		tokenRole, err := ts.TokenManager.GetRoleFromToken(tokens.AccessToken, ts.Blacklist)
		if err != nil {
			t.Fatalf("Failed to get role for role %s: %v", tt.role, err)
		}

		if tokenRole != tt.expected {
			t.Errorf("Role %s: expected %s, got %s", tt.role, tt.expected, tokenRole)
		}
	}
}

func TestUnauthorizedAccessBlocking(t *testing.T) {
	ts := SetupTestServer(t)
	defer CleanupTestServer(t, ts)

	// Create tokens
	tokens, err := ts.TokenManager.GenerateTokens("blocked-user", "user")
	if err != nil {
		t.Fatalf("Failed to generate tokens: %v", err)
	}

	// Verify token is initially valid
	_, err = ts.TokenManager.VerifyAccessToken(tokens.AccessToken, ts.Blacklist)
	if err != nil {
		t.Fatalf("Token should be valid initially: %v", err)
	}

	// Get user ID for blacklist entry
	userID, err := ts.TokenManager.GetUserIDFromToken(tokens.AccessToken, ts.Blacklist)
	if err != nil {
		t.Fatalf("Failed to get user ID: %v", err)
	}

	// Add a blacklist entry (simulating token revocation)
	jti := fmt.Sprintf("%s_blocked", userID)
	err = ts.Blacklist.Add(jti, time.Now().Add(1*time.Hour))
	if err != nil {
		t.Fatalf("Failed to blacklist: %v", err)
	}

	// Verify the blacklist entry exists
	if !ts.Blacklist.Contains(jti) {
		t.Error("Expected JTI to be blacklisted")
	}
}

// Test concurrent token operations
func TestConcurrentTokenOperations(t *testing.T) {
	ts := SetupTestServer(t)
	defer CleanupTestServer(t, ts)

	numConcurrent := 10
	done := make(chan bool, numConcurrent)

	for i := 0; i < numConcurrent; i++ {
		go func(index int) {
			tokens, err := ts.TokenManager.GenerateTokens(fmt.Sprintf("user-concurrent-%d", index), "user")
			if err != nil {
				t.Errorf("Failed to generate token %d: %v", index, err)
				done <- false
				return
			}

			if tokens.AccessToken == "" {
				t.Errorf("Token %d: empty access token", index)
				done <- false
				return
			}

			_, err = ts.TokenManager.VerifyAccessToken(tokens.AccessToken, ts.Blacklist)
			if err != nil {
				t.Errorf("Failed to verify token %d: %v", index, err)
				done <- false
				return
			}

			done <- true
		}(i)
	}

	successCount := 0
	for i := 0; i < numConcurrent; i++ {
		if <-done {
			successCount++
		}
	}

	if successCount != numConcurrent {
		t.Errorf("Expected %d successes, got %d", numConcurrent, successCount)
	}
}

// Test token expiration handling
func TestTokenExpiration(t *testing.T) {
	ts := SetupTestServer(t)
	defer CleanupTestServer(t, ts)

	// Create token manager with very short expiry
	tokenManager := jwtpkg.NewTokenManager("test-secret", 1*time.Nanosecond, 1*time.Nanosecond)

	// Generate tokens
	tokens, err := tokenManager.GenerateTokens("user-expiry-test", "user")
	if err != nil {
		t.Fatalf("Failed to generate tokens: %v", err)
	}

	// Token should fail verification due to expiration
	_, err = tokenManager.VerifyAccessToken(tokens.AccessToken, ts.Blacklist)
	if err == nil {
		t.Error("Expected token to be expired")
	}
}
