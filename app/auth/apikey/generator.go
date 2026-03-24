package apikey

import (
	"crypto/rand"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

const (
	keyPrefix = "david_"
	keyLength = 32
)

// Generator generates API keys
type Generator struct {
	hashCost int
}

// NewGenerator creates a new API key generator
func NewGenerator(hashCost int) *Generator {
	if hashCost < 4 {
		hashCost = bcrypt.DefaultCost
	}
	return &Generator{hashCost: hashCost}
}

// Generate creates a new API key and returns both the plain key and its hash
func (g *Generator) Generate() (string, string, error) {
	// Generate random bytes
	bytes := make([]byte, keyLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	// Convert to hex string
	key := fmt.Sprintf("%x", bytes)
	key = keyPrefix + key

	// Hash the key
	hash, err := bcrypt.GenerateFromPassword([]byte(key), g.hashCost)
	if err != nil {
		return "", "", fmt.Errorf("failed to hash key: %w", err)
	}

	return key, string(hash), nil
}

// Hash verifies a key against its hash
func (g *Generator) Hash(key string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(key), g.hashCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash key: %w", err)
	}
	return string(hash), nil
}

// Verify compares a key with its hash
func (g *Generator) Verify(key, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(key))
	return err == nil
}

// ValidateFormat checks if a key has the correct format
func (g Generator) ValidateFormat(key string) bool {
	return strings.HasPrefix(key, keyPrefix) && len(key) >= len(keyPrefix)+4
}
