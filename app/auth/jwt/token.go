package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/audstanley/david/app/auth/common"
	"github.com/golang-jwt/jwt/v5"
)

// TokenManager handles JWT token operations
type TokenManager struct {
	secretKey     []byte
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

// NewTokenManager creates a new JWT token manager
func NewTokenManager(secretKey string, accessExpiry, refreshExpiry time.Duration) *TokenManager {
	return &TokenManager{
		secretKey:     []byte(secretKey),
		accessExpiry:  accessExpiry,
		refreshExpiry: refreshExpiry,
	}
}

// GenerateTokens generates both access and refresh tokens for a user
func (m *TokenManager) GenerateTokens(userID, role string) (*Tokens, error) {
	accessToken, err := m.GenerateAccessToken(userID, role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := m.GenerateRefreshToken(userID, role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// GenerateAccessToken generates a short-lived access token
func (m *TokenManager) GenerateAccessToken(userID, role string) (string, error) {
	now := time.Now()
	expiresAt := now.Add(m.accessExpiry)

	jwtID := fmt.Sprintf("%d", now.UnixNano())

	claims := jwt.MapClaims{
		"user_id":    userID,
		"role":       role,
		"token_type": common.TokenTypeAccess,
		"jti":        jwtID,
		"exp":        expiresAt.Unix(),
		"iat":        now.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secretKey)
}

// GenerateRefreshToken generates a long-lived refresh token
func (m *TokenManager) GenerateRefreshToken(userID, role string) (string, error) {
	now := time.Now()
	expiresAt := now.Add(m.refreshExpiry)

	jwtID := fmt.Sprintf("%d", now.UnixNano())

	claims := jwt.MapClaims{
		"user_id":    userID,
		"role":       role,
		"token_type": common.TokenTypeRefresh,
		"jti":        jwtID,
		"exp":        expiresAt.Unix(),
		"iat":        now.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secretKey)
}

// VerifyAccessToken verifies an access token and returns claims
func (m *TokenManager) VerifyAccessToken(tokenString string, blacklist *Blacklist) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	// Check token type
	tokenType, ok := claims["token_type"].(string)
	if !ok || tokenType != common.TokenTypeAccess {
		return nil, errors.New("invalid token type")
	}

	// Check blacklist
	jti, ok := claims["jti"].(string)
	if ok && blacklist != nil && blacklist.Contains(jti) {
		return nil, errors.New("token has been revoked")
	}

	// Check expiry
	expiresAt, ok := claims["exp"].(float64)
	if !ok {
		return nil, errors.New("missing expiration claim")
	}

	if time.Now().Unix() > int64(expiresAt) {
		return nil, errors.New("token has expired")
	}

	return token, nil
}

// GetUserIDFromToken extracts user ID from a verified access token
func (m *TokenManager) GetUserIDFromToken(tokenString string, blacklist *Blacklist) (string, error) {
	token, err := m.VerifyAccessToken(tokenString, blacklist)
	if err != nil {
		return "", err
	}

	claims := token.Claims.(jwt.MapClaims)
	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", errors.New("missing user_id claim")
	}

	return userID, nil
}

// GetRoleFromToken extracts role from a verified access token
func (m *TokenManager) GetRoleFromToken(tokenString string, blacklist *Blacklist) (string, error) {
	token, err := m.VerifyAccessToken(tokenString, blacklist)
	if err != nil {
		return "", err
	}

	claims := token.Claims.(jwt.MapClaims)
	role, ok := claims["role"].(string)
	if !ok {
		return "", errors.New("missing role claim")
	}

	return role, nil
}

// VerifyRefreshToken verifies a refresh token and returns claims
func (m *TokenManager) VerifyRefreshToken(tokenString string, blacklist *Blacklist) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	// Check token type
	tokenType, ok := claims["token_type"].(string)
	if !ok || tokenType != common.TokenTypeRefresh {
		return nil, errors.New("invalid token type")
	}

	// Check blacklist
	jti, ok := claims["jti"].(string)
	if ok && blacklist != nil && blacklist.Contains(jti) {
		return nil, errors.New("token has been revoked")
	}

	// Check expiry
	expiresAt, ok := claims["exp"].(float64)
	if !ok {
		return nil, errors.New("missing expiration claim")
	}

	if time.Now().Unix() > int64(expiresAt) {
		return nil, errors.New("token has expired")
	}

	return token, nil
}

// Tokens holds both access and refresh tokens
type Tokens struct {
	AccessToken  string
	RefreshToken string
}
