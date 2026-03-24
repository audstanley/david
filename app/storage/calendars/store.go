package calendars

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/audstanley/david/app/storage"
)

const (
	keyPrefixCalendar    = "cal:"
	keyPrefixCalendarUID = "uid:"
	keyPrefixCalendarBy  = "by:"
	keyPrefixCalendarPub = "pub:"
)

// Store implements CalendarStore
type Store struct {
	db *storage.Storage
}

// New creates a new calendar store
func New(db *storage.Storage) *Store {
	return &Store{db: db}
}

// Create creates a new calendar
func (s *Store) Create(calendarUID, ownerID, displayName, description, color, timezone string, isPublic bool) error {
	// Check if calendar UID already exists
	exists, err := s.db.Exists([]byte(keyPrefixCalendarUID + calendarUID))
	if err != nil {
		return fmt.Errorf("failed to check calendar UID: %w", err)
	}
	if exists {
		return storage.ErrAlreadyExists
	}

	cal := &storage.Calendar{
		UID:         calendarUID,
		OwnerID:     ownerID,
		DisplayName: displayName,
		Description: description,
		Color:       color,
		IsPublic:    isPublic,
		PublicHash:  generatePublicHash(calendarUID),
		Timezone:    timezone,
		Created:     now(),
		Updated:     now(),
	}

	data, err := storage.Encode(cal)
	if err != nil {
		return fmt.Errorf("failed to encode calendar: %w", err)
	}

	pairs := [][2][]byte{
		{[]byte(keyPrefixCalendar + calendarUID), data},
		{[]byte(keyPrefixCalendarUID + calendarUID), []byte(ownerID)},
		{[]byte(keyPrefixCalendarBy + ownerID + ":" + displayName), []byte(calendarUID)},
	}

	if isPublic {
		pairs = append(pairs, [][2][]byte{
			{[]byte(keyPrefixCalendarPub + cal.PublicHash), []byte(calendarUID)},
		}...)
	}

	return s.db.Batch(pairs)
}

// GetByUID retrieves a calendar by UID
func (s *Store) GetByUID(calendarUID string) (*storage.Calendar, error) {
	data, err := s.db.Get([]byte(keyPrefixCalendar + calendarUID))
	if err != nil {
		return nil, err
	}

	var cal storage.Calendar
	if err := storage.Decode(data, &cal); err != nil {
		return nil, fmt.Errorf("failed to decode calendar: %w", err)
	}

	return &cal, nil
}

// GetByName retrieves a calendar by name (owner + name)
func (s *Store) GetByName(ownerID, name string) (*storage.Calendar, error) {
	calendarUIDData, err := s.db.Get([]byte(keyPrefixCalendarBy + ownerID + ":" + name))
	if err != nil {
		return nil, err
	}

	return s.GetByUID(string(calendarUIDData))
}

// Update updates a calendar
func (s *Store) Update(calendarUID string, updates map[string]interface{}) error {
	cal, err := s.GetByUID(calendarUID)
	if err != nil {
		return err
	}

	// Apply updates
	if v, ok := updates["displayName"]; ok {
		cal.DisplayName = v.(string)
	}
	if v, ok := updates["description"]; ok {
		cal.Description = v.(string)
	}
	if v, ok := updates["color"]; ok {
		cal.Color = v.(string)
	}
	if v, ok := updates["timezone"]; ok {
		cal.Timezone = v.(string)
	}
	if v, ok := updates["isPublic"]; ok {
		cal.IsPublic = v.(bool)
	}

	cal.Updated = now()

	// Regenerate public hash if isPublic changed
	if _, ok := updates["isPublic"]; ok {
		cal.PublicHash = generatePublicHash(calendarUID)
	}

	data, err := storage.Encode(cal)
	if err != nil {
		return fmt.Errorf("failed to encode calendar: %w", err)
	}

	// Update name index if name changed
	oldName := cal.DisplayName
	if _, ok := updates["displayName"]; !ok {
		return s.db.Delete([]byte(keyPrefixCalendarBy + cal.OwnerID + ":" + oldName))
	}

	var batches [][2][]byte
	batches = append(batches, [2][]byte{[]byte(keyPrefixCalendar + calendarUID), data})
	batches = append(batches, [2][]byte{[]byte(keyPrefixCalendarBy + cal.OwnerID + cal.DisplayName), []byte(calendarUID)})

	if cal.DisplayName != oldName {
		batches = append(batches, [2][]byte{[]byte(keyPrefixCalendarBy + cal.OwnerID + oldName), nil})
	}

	// Update public index if isPublic changed
	if _, ok := updates["isPublic"]; ok {
		if cal.IsPublic {
			batches = append(batches, [2][]byte{[]byte(keyPrefixCalendarPub + cal.PublicHash), []byte(calendarUID)})
		} else {
			batches = append(batches, [2][]byte{[]byte(keyPrefixCalendarPub + cal.PublicHash), nil})
		}
	}

	if err := s.db.Batch(batches); err != nil {
		return err
	}

	return nil
}

// Delete deletes a calendar
func (s *Store) Delete(calendarUID string) error {
	cal, err := s.GetByUID(calendarUID)
	if err != nil {
		return err
	}

	pairs := [][2][]byte{
		{[]byte(keyPrefixCalendar + calendarUID), nil},
		{[]byte(keyPrefixCalendarUID + calendarUID), nil},
		{[]byte(keyPrefixCalendarBy + cal.OwnerID + cal.DisplayName), nil},
	}

	if cal.IsPublic {
		pairs = append(pairs, [][2][]byte{
			{[]byte(keyPrefixCalendarPub + cal.PublicHash), nil},
		}...)
	}

	return s.db.Batch(pairs)
}

// List lists calendars for a user with pagination
func (s *Store) List(ownerID string, limit, offset int) ([]*storage.Calendar, error) {
	if limit <= 0 {
		limit = 100
	}
	if limit > 1000 {
		limit = 1000
	}

	var calendars []*storage.Calendar
	count := 0

	err := s.db.Iterate([]byte(keyPrefixCalendar), func(key, value []byte) error {
		var cal storage.Calendar
		if err := storage.Decode(value, &cal); err != nil {
			return err
		}

		if cal.OwnerID == ownerID {
			if count >= offset {
				calendars = append(calendars, &cal)
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

	return calendars, nil
}

// IsPublic checks if a calendar is public
func (s *Store) IsPublic(calendarUID string) (bool, error) {
	cal, err := s.GetByUID(calendarUID)
	if err != nil {
		return false, err
	}
	return cal.IsPublic, nil
}

// generatePublicHash generates a public hash for a calendar
func generatePublicHash(calendarUID string) string {
	hash := sha256.Sum256([]byte(calendarUID))
	return hex.EncodeToString(hash[:])[:12]
}

// now returns the current time
func now() time.Time {
	return time.Now().UTC()
}
