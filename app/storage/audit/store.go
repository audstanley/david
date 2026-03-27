package audit

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/audstanley/david/app/storage"
)

const (
	keyPrefixAudit       = "audit:"
	keyPrefixAuditBy     = "by:"
	keyPrefixAuditByDate = "date:"
)

// Store implements AuditStore
type Store struct {
	db *storage.Storage
}

// New creates a new audit store
func New(db *storage.Storage) *Store {
	return &Store{db: db}
}

// Log logs an audit entry
func (s *Store) Log(userID, action, entityType, entityID, details string) error {
	entry := &storage.AuditEntry{
		ID:         generateAuditID(),
		Timestamp:  now(),
		UserID:     userID,
		Action:     action,
		EntityType: entityType,
		EntityID:   entityID,
		Details:    details,
	}

	data, err := storage.Encode(entry)
	if err != nil {
		return fmt.Errorf("failed to encode audit entry: %w", err)
	}

	dateKey := entry.Timestamp.Format("20060102")
	pairs := [][2][]byte{
		{[]byte(keyPrefixAudit + entry.ID), data},
		{[]byte(keyPrefixAuditBy + userID + ":" + entry.ID), data},
		{[]byte(keyPrefixAuditBy + entityType + ":" + entityID + ":" + entry.ID), data},
		{[]byte(keyPrefixAuditByDate + dateKey + ":" + entry.ID), data},
	}

	return s.db.Batch(pairs)
}

// GetByUser retrieves audit entries for a user with pagination
func (s *Store) GetByUser(userID string, limit, offset int) ([]*storage.AuditEntry, error) {
	if limit <= 0 {
		limit = 100
	}
	if limit > 1000 {
		limit = 1000
	}

	var entries []*storage.AuditEntry
	skipped := 0
	done := false

	err := s.db.Iterate([]byte(keyPrefixAuditBy+userID+":"), func(key, value []byte) error {
		if done {
			return nil
		}
		if skipped < offset {
			skipped++
			return nil
		}
		var entry storage.AuditEntry
		if err := storage.Decode(value, &entry); err != nil {
			return err
		}
		entries = append(entries, &entry)
		if len(entries) >= limit {
			done = true
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return entries, nil
}

// GetByEntity retrieves audit entries for an entity with pagination
func (s *Store) GetByEntity(entityType, entityID string, limit, offset int) ([]*storage.AuditEntry, error) {
	if limit <= 0 {
		limit = 100
	}
	if limit > 1000 {
		limit = 1000
	}

	var entries []*storage.AuditEntry
	skipped := 0
	done := false

	err := s.db.Iterate([]byte(keyPrefixAuditBy+entityType+":"+entityID+":"), func(key, value []byte) error {
		if done {
			return nil
		}
		if skipped < offset {
			skipped++
			return nil
		}
		var entry storage.AuditEntry
		if err := storage.Decode(value, &entry); err != nil {
			return err
		}
		entries = append(entries, &entry)
		if len(entries) >= limit {
			done = true
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return entries, nil
}

// GetByDateRange retrieves audit entries in a date range with pagination
func (s *Store) GetByDateRange(start, end time.Time, limit, offset int) ([]*storage.AuditEntry, error) {
	if limit <= 0 {
		limit = 100
	}
	if limit > 1000 {
		limit = 1000
	}

	var entries []*storage.AuditEntry
	skipped := 0
	done := false

	startPrefix := keyPrefixAuditByDate + start.Format("20060102") + ":"
	endPrefix := keyPrefixAuditByDate + end.Add(24*time.Hour).Format("20060102") + ":"

	err := s.db.Iterate([]byte(startPrefix), func(key, value []byte) error {
		if done {
			return nil
		}
		keyStr := string(key)
		if keyStr >= endPrefix {
			done = true
			return nil
		}
		if skipped < offset {
			skipped++
			return nil
		}
		var entry storage.AuditEntry
		if err := storage.Decode(value, &entry); err != nil {
			return err
		}
		entries = append(entries, &entry)
		if len(entries) >= limit {
			done = true
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return entries, nil
}

// generateAuditID generates a unique audit entry ID
func generateAuditID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// now returns the current time
func now() time.Time {
	return time.Now().UTC()
}

// Encode serializes a slice of strings
func encodeStringSlice(items []string) (string, error) {
	data, err := json.Marshal(items)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// DecodeStringSlice deserializes a string slice
func decodeStringSlice(s string) ([]string, error) {
	var items []string
	if err := json.Unmarshal([]byte(s), &items); err != nil {
		return nil, err
	}
	return items, nil
}
