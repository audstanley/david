package recurrence

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/audstanley/david/app/storage"
)

const (
	keyPrefixRecurrence   = "recurrence:"
	keyPrefixRecurrenceBy = "by:"
)

// Store implements RecurrenceStore
type Store struct {
	db *storage.Storage
}

// New creates a new recurrence store
func New(db *storage.Storage) *Store {
	return &Store{db: db}
}

// Store stores expanded recurrence instances
func (s *Store) Store(eventUID string, instances []storage.RecurrenceInstance) error {
	// Delete existing instances for this event
	existing, err := s.Get(eventUID, time.Time{}, time.Time{})
	if err != nil && !storage.IsNotFound(err) {
		return fmt.Errorf("failed to get existing instances: %w", err)
	}

	var pairs [][2][]byte

	// Remove old instances
	for _, inst := range existing {
		pairs = append(pairs, [][2][]byte{
			{[]byte(keyPrefixRecurrence + inst.InstanceUID), nil},
			{[]byte(keyPrefixRecurrenceBy + eventUID + ":" + inst.InstanceUID), nil},
		}...)
	}

	// Add new instances
	for _, inst := range instances {
		data, err := storage.Encode(inst)
		if err != nil {
			return fmt.Errorf("failed to encode recurrence instance: %w", err)
		}

		pairs = append(pairs, [][2][]byte{
			{[]byte(keyPrefixRecurrence + inst.InstanceUID), data},
			{[]byte(keyPrefixRecurrenceBy + eventUID + ":" + inst.InstanceUID), data},
		}...)
	}

	if len(pairs) > 0 {
		if err := s.db.Batch(pairs); err != nil {
			return err
		}
	}

	return nil
}

// Get retrieves recurrence instances for an event in a date range
func (s *Store) Get(eventUID string, start, end time.Time) ([]storage.RecurrenceInstance, error) {
	prefix := keyPrefixRecurrenceBy + eventUID + ":"

	var instances []storage.RecurrenceInstance

	err := s.db.Iterate([]byte(prefix), func(key, value []byte) error {
		var inst storage.RecurrenceInstance
		if err := storage.Decode(value, &inst); err != nil {
			return err
		}

		// Filter by date range if provided
		if (!start.IsZero() && inst.DTStart.Before(start)) ||
			(!end.IsZero() && inst.DTStart.After(end)) {
			return nil
		}

		instances = append(instances, inst)
		return nil
	})

	if err != nil {
		return nil, err
	}

	if len(instances) == 0 {
		return nil, storage.ErrNotFound
	}

	return instances, nil
}

// Delete deletes all recurrence instances for an event
func (s *Store) Delete(eventUID string) error {
	instances, err := s.Get(eventUID, time.Time{}, time.Time{})
	if err != nil && !storage.IsNotFound(err) {
		return fmt.Errorf("failed to get instances: %w", err)
	}

	var pairs [][2][]byte
	for _, inst := range instances {
		pairs = append(pairs, [][2][]byte{
			{[]byte(keyPrefixRecurrence + inst.InstanceUID), nil},
			{[]byte(keyPrefixRecurrenceBy + eventUID + ":" + inst.InstanceUID), nil},
		}...)
	}

	if len(pairs) > 0 {
		return s.db.Batch(pairs)
	}

	return nil
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
