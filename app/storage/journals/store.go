package journals

import (
	"fmt"
	"time"

	"github.com/audstanley/david/app/storage"
)

const (
	keyPrefixJournal      = "journal:"
	keyPrefixJournalByUID = "uid:"
	keyPrefixJournalByCal = "cal:"
)

// Store implements JournalStore
type Store struct {
	db *storage.Storage
}

// New creates a new journal store
func New(db *storage.Storage) *Store {
	return &Store{db: db}
}

// Create creates a new journal
func (s *Store) Create(journal *storage.Journal) error {
	// Check if journal UID already exists
	exists, err := s.db.Exists([]byte(keyPrefixJournalByUID + journal.UID))
	if err != nil {
		return fmt.Errorf("failed to check journal UID: %w", err)
	}
	if exists {
		return storage.ErrAlreadyExists
	}

	data, err := storage.Encode(journal)
	if err != nil {
		return fmt.Errorf("failed to encode journal: %w", err)
	}

	pairs := [][2][]byte{
		{[]byte(keyPrefixJournal + journal.UID), data},
		{[]byte(keyPrefixJournalByUID + journal.UID), []byte(journal.CalendarUID)},
		{[]byte(keyPrefixJournalByCal + journal.CalendarUID + ":" + journal.UID), data},
	}

	return s.db.Batch(pairs)
}

// GetByUID retrieves a journal by UID
func (s *Store) GetByUID(journalUID string) (*storage.Journal, error) {
	data, err := s.db.Get([]byte(keyPrefixJournal + journalUID))
	if err != nil {
		return nil, err
	}

	var journal storage.Journal
	if err := storage.Decode(data, &journal); err != nil {
		return nil, fmt.Errorf("failed to decode journal: %w", err)
	}

	return &journal, nil
}

// Update updates a journal
func (s *Store) Update(journalUID string, journal *storage.Journal) error {
	existing, err := s.GetByUID(journalUID)
	if err != nil {
		return err
	}

	journal.Sequence = existing.Sequence + 1
	journal.LastModified = now()

	data, err := storage.Encode(journal)
	if err != nil {
		return fmt.Errorf("failed to encode journal: %w", err)
	}

	pairs := [][2][]byte{
		{[]byte(keyPrefixJournal + journalUID), data},
		{[]byte(keyPrefixJournalByCal + journal.CalendarUID + ":" + journalUID), data},
	}

	return s.db.Batch(pairs)
}

// Delete deletes a journal
func (s *Store) Delete(journalUID string) error {
	journal, err := s.GetByUID(journalUID)
	if err != nil {
		return err
	}

	pairs := [][2][]byte{
		{[]byte(keyPrefixJournal + journalUID), nil},
		{[]byte(keyPrefixJournalByUID + journalUID), nil},
		{[]byte(keyPrefixJournalByCal + journal.CalendarUID + ":" + journalUID), nil},
	}

	return s.db.Batch(pairs)
}

// List lists journals for a calendar with pagination
func (s *Store) List(calendarUID string, limit, offset int) ([]*storage.Journal, error) {
	if limit <= 0 {
		limit = 100
	}
	if limit > 1000 {
		limit = 1000
	}

	var journals []*storage.Journal
	count := 0

	err := s.db.Iterate([]byte(keyPrefixJournal), func(key, value []byte) error {
		var journal storage.Journal
		if err := storage.Decode(value, &journal); err != nil {
			return err
		}

		if journal.CalendarUID == calendarUID {
			if count < offset {
				count++
				return nil
			}
			journals = append(journals, &journal)
			count++
			if len(journals) >= limit {
				return storage.ErrLimitReached
			}
		}
		return nil
	})

	if err != nil && err != storage.ErrLimitReached {
		return nil, err
	}

	return journals, nil
}

// now returns the current time
func now() time.Time {
	return time.Now().UTC()
}
