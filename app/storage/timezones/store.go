package timezones

import (
	"encoding/json"
	"fmt"

	"github.com/audstanley/david/app/storage"
)

const (
	keyPrefixTZ       = "tz:"
	keyPrefixTZByTzID = "tzid:"
)

// Store implements TimeZoneStore
type Store struct {
	db *storage.Storage
}

// New creates a new timezone store
func New(db *storage.Storage) *Store {
	return &Store{db: db}
}

// Store stores a VTIMEZONE definition
func (s *Store) Store(tzID string, vtimezone *storage.VTimeZone) error {
	// Check if timezone already exists
	exists, err := s.db.Exists([]byte(keyPrefixTZByTzID + tzID))
	if err != nil {
		return fmt.Errorf("failed to check timezone: %w", err)
	}
	if exists {
		return storage.ErrAlreadyExists
	}

	data, err := storage.Encode(vtimezone)
	if err != nil {
		return fmt.Errorf("failed to encode timezone: %w", err)
	}

	pairs := [][2][]byte{
		{[]byte(keyPrefixTZ + tzID), data},
		{[]byte(keyPrefixTZByTzID + tzID), []byte(tzID)},
	}

	return s.db.Batch(pairs)
}

// Get retrieves a timezone by TZID
func (s *Store) Get(tzID string) (*storage.VTimeZone, error) {
	data, err := s.db.Get([]byte(keyPrefixTZ + tzID))
	if err != nil {
		return nil, err
	}

	var vtimezone storage.VTimeZone
	if err := storage.Decode(data, &vtimezone); err != nil {
		return nil, fmt.Errorf("failed to decode timezone: %w", err)
	}

	return &vtimezone, nil
}

// Update updates a timezone
func (s *Store) Update(tzID string, vtimezone *storage.VTimeZone) error {
	existing, err := s.Get(tzID)
	if err != nil {
		return err
	}

	data, err := storage.Encode(vtimezone)
	if err != nil {
		return fmt.Errorf("failed to encode timezone: %w", err)
	}

	pairs := [][2][]byte{
		{[]byte(keyPrefixTZ + tzID), data},
	}

	// Remove old TZID index if TZID changed
	if existing.TZID != tzID {
		pairs = append(pairs, [][2][]byte{
			{[]byte(keyPrefixTZByTzID + existing.TZID), nil},
		}...)
	}

	// Add new TZID index
	pairs = append(pairs, [][2][]byte{
		{[]byte(keyPrefixTZByTzID + tzID), []byte(tzID)},
	}...)

	return s.db.Batch(pairs)
}

// Delete deletes a timezone
func (s *Store) Delete(tzID string) error {
	existing, err := s.Get(tzID)
	if err != nil {
		return err
	}

	pairs := [][2][]byte{
		{[]byte(keyPrefixTZ + tzID), nil},
		{[]byte(keyPrefixTZByTzID + existing.TZID), nil},
	}

	return s.db.Batch(pairs)
}

// Encode serializes a slice of components
func encodeComponents(comps []interface{}) (string, error) {
	data, err := json.Marshal(comps)
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
