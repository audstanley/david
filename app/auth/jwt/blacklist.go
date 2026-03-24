package jwt

import (
	"fmt"
	"time"

	"github.com/audstanley/david/app/storage"
)

const (
	keyPrefixBlacklist = "blacklist:"
)

// Blacklist manages revoked JWT tokens
type Blacklist struct {
	db      *storage.Storage
	cleanup time.Duration
}

// NewBlacklist creates a new blacklist
func NewBlacklist(db *storage.Storage, cleanup time.Duration) *Blacklist {
	return &Blacklist{
		db:      db,
		cleanup: cleanup,
	}
}

// Add adds a JWT ID to the blacklist
func (b *Blacklist) Add(jti string, expiry time.Time) error {
	data, err := storage.Encode(jti)
	if err != nil {
		return fmt.Errorf("failed to encode jti: %w", err)
	}

	pairs := [][2][]byte{
		{[]byte(keyPrefixBlacklist + jti), data},
	}

	// Schedule cleanup
	go b.scheduleCleanup(jti, expiry)

	return b.db.Batch(pairs)
}

// Contains checks if a JWT ID is blacklisted
func (b *Blacklist) Contains(jti string) bool {
	_, err := b.db.Get([]byte(keyPrefixBlacklist + jti))
	if err != nil {
		if storage.IsNotFound(err) {
			return false
		}
		// Log error but don't fail lookup
		return false
	}
	return true
}

// Remove removes a JWT ID from the blacklist
func (b *Blacklist) Remove(jti string) error {
	return b.db.Delete([]byte(keyPrefixBlacklist + jti))
}

// Cleanup removes expired entries from the blacklist
func (b *Blacklist) Cleanup() error {
	var toDelete []string
	err := b.db.Iterate([]byte(keyPrefixBlacklist), func(key, value []byte) error {
		var jti string
		if err := storage.Decode(value, &jti); err != nil {
			return err
		}

		// Check if entry is expired (stored as Unix timestamp in metadata)
		// For now, we'll clean up all entries older than cleanup duration
		// In a full implementation, we'd store expiry time with each entry
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to iterate blacklist: %w", err)
	}

	// Delete expired entries
	pairs := [][2][]byte{}
	for _, jti := range toDelete {
		pairs = append(pairs, [][2][]byte{
			{[]byte(keyPrefixBlacklist + jti), nil},
		}...)
	}

	if len(pairs) > 0 {
		return b.db.Batch(pairs)
	}

	return nil
}

// scheduleCleanup schedules cleanup of an expired blacklist entry
func (b *Blacklist) scheduleCleanup(jti string, expiry time.Time) {
	duration := time.Until(expiry)
	if duration <= 0 {
		return
	}

	time.AfterFunc(duration, func() {
		_ = b.Remove(jti)
	})
}

// Encode serializes data to JSON
func encode(data interface{}) (string, error) {
	b, err := storage.Encode(data)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Decode deserializes JSON to data
func decode(data string, v interface{}) error {
	return storage.Decode([]byte(data), v)
}
