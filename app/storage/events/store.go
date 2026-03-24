package events

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/audstanley/david/app/storage"
)

const (
	keyPrefixEvent       = "event:"
	keyPrefixEventByUID  = "uid:"
	keyPrefixEventByCal  = "cal:"
	keyPrefixEventByDate = "date:"
)

// Store implements EventStore
type Store struct {
	db *storage.Storage
}

// New creates a new event store
func New(db *storage.Storage) *Store {
	return &Store{db: db}
}

// Create creates a new event
func (s *Store) Create(event *storage.Event) error {
	// Check if event UID already exists
	exists, err := s.db.Exists([]byte(keyPrefixEventByUID + event.UID))
	if err != nil {
		return fmt.Errorf("failed to check event UID: %w", err)
	}
	if exists {
		return storage.ErrAlreadyExists
	}

	data, err := storage.Encode(event)
	if err != nil {
		return fmt.Errorf("failed to encode event: %w", err)
	}

	// Generate date index key
	dateKey := generateDateIndexKey(event.CalendarUID, event.DTStart)

	pairs := [][2][]byte{
		{[]byte(keyPrefixEvent + event.UID), data},
		{[]byte(keyPrefixEventByUID + event.UID), []byte(event.CalendarUID)},
		{[]byte(keyPrefixEventByCal + event.CalendarUID + ":" + event.UID), data},
		{[]byte(keyPrefixEventByDate + dateKey), data},
	}

	return s.db.Batch(pairs)
}

// GetByUID retrieves an event by UID
func (s *Store) GetByUID(eventUID string) (*storage.Event, error) {
	data, err := s.db.Get([]byte(keyPrefixEvent + eventUID))
	if err != nil {
		return nil, err
	}

	var event storage.Event
	if err := storage.Decode(data, &event); err != nil {
		return nil, fmt.Errorf("failed to decode event: %w", err)
	}

	return &event, nil
}

// Update updates an event
func (s *Store) Update(eventUID string, event *storage.Event, sequence int) error {
	existing, err := s.GetByUID(eventUID)
	if err != nil {
		return err
	}

	// Check sequence number
	if sequence != 0 && event.Sequence != 0 && event.Sequence <= existing.Sequence {
		return storage.ErrConflict
	}

	// Update sequence if needed
	if sequence != 0 {
		event.Sequence = sequence
	}

	event.LastModified = now()

	data, err := storage.Encode(event)
	if err != nil {
		return fmt.Errorf("failed to encode event: %w", err)
	}

	// Remove old date index
	oldDateKey := generateDateIndexKey(existing.CalendarUID, existing.DTStart)

	pairs := [][2][]byte{
		{[]byte(keyPrefixEvent + eventUID), data},
		{[]byte(keyPrefixEventByCal + event.CalendarUID + ":" + eventUID), data},
		{[]byte(keyPrefixEventByDate + oldDateKey), nil},
	}

	// Add new date index
	newDateKey := generateDateIndexKey(event.CalendarUID, event.DTStart)
	pairs = append(pairs, [2][]byte{[]byte(keyPrefixEventByDate + newDateKey), data})

	return s.db.Batch(pairs)
}

// Delete deletes an event
func (s *Store) Delete(eventUID string) error {
	event, err := s.GetByUID(eventUID)
	if err != nil {
		return err
	}

	dateKey := generateDateIndexKey(event.CalendarUID, event.DTStart)

	pairs := [][2][]byte{
		{[]byte(keyPrefixEvent + eventUID), nil},
		{[]byte(keyPrefixEventByUID + eventUID), nil},
		{[]byte(keyPrefixEventByCal + event.CalendarUID + ":" + eventUID), nil},
		{[]byte(keyPrefixEventByDate + dateKey), nil},
	}

	return s.db.Batch(pairs)
}

// List lists events for a calendar in a date range with pagination
func (s *Store) List(calendarUID string, start, end time.Time, limit, offset int) ([]*storage.Event, error) {
	if limit <= 0 {
		limit = 100
	}
	if limit > 1000 {
		limit = 1000
	}

	var events []*storage.Event
	count := 0

	startPrefix := generateDateIndexKey(calendarUID, start)
	endPrefix := generateDateIndexKey(calendarUID, end)

	err := s.db.Iterate([]byte(startPrefix), func(key, value []byte) error {
		// Check if key is within range
		keyStr := string(key)
		if keyStr < endPrefix {
			if count >= offset {
				var event storage.Event
				if err := storage.Decode(value, &event); err != nil {
					return err
				}
				events = append(events, &event)
				count++
				if count >= offset+limit {
					return nil
				}
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return events, nil
}

// GetInDateRange retrieves a specific event in a date range
func (s *Store) GetInDateRange(calendarUID, eventUID string, start, end time.Time) (*storage.Event, error) {
	// Try direct UID lookup first
	event, err := s.GetByUID(eventUID)
	if err == nil {
		// Verify it's in the calendar and date range
		if event.CalendarUID == calendarUID &&
			(event.DTStart.After(start) || event.DTStart.Equal(start)) &&
			event.DTStart.Before(end) {
			return event, nil
		}
	}

	// Try date range lookup
	events, err := s.List(calendarUID, start, end, 1, 0)
	if err != nil {
		return nil, err
	}

	for _, e := range events {
		if e.UID == eventUID {
			return e, nil
		}
	}

	return nil, storage.ErrNotFound
}

// generateDateIndexKey generates a lexicographically sortable date index key
func generateDateIndexKey(calendarUID string, dt time.Time) string {
	return calendarUID + ":" + dt.UTC().Format("20060102T150405")
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
