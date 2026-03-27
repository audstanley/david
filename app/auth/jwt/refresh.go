package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"
)

// RefreshManager handles refresh token operations
type RefreshManager struct {
	tokenManager *TokenManager
	blacklist    *Blacklist
	rotateTokens bool
}

// NewRefreshManager creates a new refresh token manager
func NewRefreshManager(tokenManager *TokenManager, blacklist *Blacklist, rotateTokens bool) *RefreshManager {
	return &RefreshManager{
		tokenManager: tokenManager,
		blacklist:    blacklist,
		rotateTokens: rotateTokens,
	}
}

// Refresh exchanges a refresh token for a new access token
func (m *RefreshManager) Refresh(refreshToken string) (*Tokens, error) {
	// Verify the refresh token
	token, err := m.tokenManager.VerifyRefreshToken(refreshToken, m.blacklist)
	if err != nil {
		return nil, err
	}

	claims := token.Claims.(jwt.MapClaims)

	// Extract user info
	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, errors.New("invalid refresh token claims")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return nil, errors.New("invalid refresh token claims")
	}

	// Generate new tokens
	newTokens, err := m.tokenManager.GenerateTokens(userID, role)
	if err != nil {
		return nil, err
	}

	// Rotate tokens (invalidate old refresh token)
	if m.rotateTokens && m.blacklist != nil {
		jti, ok := claims["jti"].(string)
		if ok {
			if err := m.blacklist.Add(jti, time.Time{}); err != nil {
				log.WithError(err).Warn("Failed to add token to blacklist during rotation")
			}
		}
	}

	return newTokens, nil
}
