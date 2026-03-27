package journals

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/audstanley/david/app/storage"
)

func TestJournalCreate(t *testing.T) {
	tmpDir := "/tmp/david-test-journals-create-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	journal := &storage.Journal{
		UID:          "journal-test-create",
		CalendarUID:  "calendar-123",
		DTStamp:      time.Now().UTC(),
		Summary:      "Test Journal",
		Description:  "This is a test journal entry",
		Status:       "DRAFT",
		Sequence:     0,
		Created:      time.Time{},
		LastModified: time.Time{},
	}

	err = store.Create(journal)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if journal.UID == "" {
		t.Error("UID should be set after creation")
	}
}

func TestJournalCreateDuplicateUID(t *testing.T) {
	tmpDir := "/tmp/david-test-journals-duplicate-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	journal1 := &storage.Journal{
		UID:          "journal-duplicate-test",
		CalendarUID:  "calendar-123",
		DTStamp:      time.Now().UTC(),
		Summary:      "Test Journal 1",
		Description:  "First entry",
		Status:       "DRAFT",
		Sequence:     0,
		Created:      time.Time{},
		LastModified: time.Time{},
	}

	err = store.Create(journal1)
	if err != nil {
		t.Fatalf("First Create failed: %v", err)
	}

	journal2 := &storage.Journal{
		UID:          "journal-duplicate-test",
		CalendarUID:  "calendar-456",
		DTStamp:      time.Now().UTC(),
		Summary:      "Test Journal 2",
		Description:  "Second entry",
		Status:       "DRAFT",
		Sequence:     0,
		Created:      time.Time{},
		LastModified: time.Time{},
	}

	err = store.Create(journal2)
	if err != storage.ErrAlreadyExists {
		t.Errorf("Expected ErrAlreadyExists, got %v", err)
	}
}

func TestJournalCreateWithAllFields(t *testing.T) {
	tmpDir := "/tmp/david-test-journals-fields-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	journal := &storage.Journal{
		UID:          "journal-custom-fields",
		CalendarUID:  "calendar-123",
		DTStamp:      time.Now().UTC(),
		Summary:      "Complete Journal",
		Description:  "Full description with details",
		Status:       "PUBLISHED",
		Sequence:     0,
		Created:      time.Time{},
		LastModified: time.Time{},
	}

	err = store.Create(journal)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	retrieved, err := store.GetByUID(journal.UID)
	if err != nil {
		t.Fatalf("GetByUID failed: %v", err)
	}

	if retrieved.Summary != journal.Summary {
		t.Errorf("Summary mismatch: got %q, want %q", retrieved.Summary, journal.Summary)
	}
	if retrieved.Description != journal.Description {
		t.Errorf("Description mismatch: got %q, want %q", retrieved.Description, journal.Description)
	}
	if retrieved.Status != journal.Status {
		t.Errorf("Status mismatch: got %q, want %q", retrieved.Status, journal.Status)
	}
}

func TestJournalGetByUID(t *testing.T) {
	tmpDir := "/tmp/david-test-journals-getbyuid-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	journal := &storage.Journal{
		UID:          "journal-getbyuid-test",
		CalendarUID:  "calendar-123",
		DTStamp:      time.Now().UTC(),
		Summary:      "Test Journal",
		Description:  "This is a test journal entry",
		Status:       "DRAFT",
		Sequence:     0,
		Created:      time.Time{},
		LastModified: time.Time{},
	}

	err = store.Create(journal)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	retrieved, err := store.GetByUID(journal.UID)
	if err != nil {
		t.Fatalf("GetByUID failed: %v", err)
	}

	if retrieved.UID != journal.UID {
		t.Errorf("UID mismatch: got %q, want %q", retrieved.UID, journal.UID)
	}
	if retrieved.CalendarUID != journal.CalendarUID {
		t.Errorf("CalendarUID mismatch: got %q, want %q", retrieved.CalendarUID, journal.CalendarUID)
	}
	if retrieved.Summary != journal.Summary {
		t.Errorf("Summary mismatch: got %q, want %q", retrieved.Summary, journal.Summary)
	}
}

func TestJournalGetByUIDNonExistent(t *testing.T) {
	tmpDir := "/tmp/david-test-journals-nonexistent-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	_, err = store.GetByUID("non-existent-uid")
	if err != storage.ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

func TestJournalUpdate(t *testing.T) {
	tmpDir := "/tmp/david-test-journals-update-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	journal := &storage.Journal{
		UID:          "journal-update-test",
		CalendarUID:  "calendar-123",
		DTStamp:      time.Now().UTC(),
		Summary:      "Test Journal",
		Description:  "Original description",
		Status:       "DRAFT",
		Sequence:     0,
		Created:      time.Time{},
		LastModified: time.Time{},
	}

	err = store.Create(journal)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	journal.Summary = "Updated Journal"
	journal.Description = "Updated description with changes"
	journal.Status = "PUBLISHED"

	err = store.Update(journal.UID, journal)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	retrieved, err := store.GetByUID(journal.UID)
	if err != nil {
		t.Fatalf("GetByUID failed: %v", err)
	}

	if retrieved.Summary != "Updated Journal" {
		t.Errorf("Summary not updated: got %q", retrieved.Summary)
	}
	if retrieved.Description != "Updated description with changes" {
		t.Errorf("Description not updated: got %q", retrieved.Description)
	}
	if retrieved.Status != "PUBLISHED" {
		t.Errorf("Status not updated: got %q", retrieved.Status)
	}
}

func TestJournalUpdateNonExistent(t *testing.T) {
	tmpDir := "/tmp/david-test-journals-update-nonexistent-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	journal := &storage.Journal{
		UID:          "journal-update-test",
		CalendarUID:  "calendar-123",
		DTStamp:      time.Now().UTC(),
		Summary:      "Test Journal",
		Description:  "Description",
		Status:       "DRAFT",
		Sequence:     0,
		Created:      time.Time{},
		LastModified: time.Time{},
	}

	err = store.Update("non-existent-uid", journal)
	if err != storage.ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

func TestJournalDelete(t *testing.T) {
	tmpDir := "/tmp/david-test-journals-delete-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	journal := &storage.Journal{
		UID:          "journal-delete-test",
		CalendarUID:  "calendar-123",
		DTStamp:      time.Now().UTC(),
		Summary:      "Test Journal",
		Description:  "This is a test journal entry",
		Status:       "DRAFT",
		Sequence:     0,
		Created:      time.Time{},
		LastModified: time.Time{},
	}

	err = store.Create(journal)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	err = store.Delete(journal.UID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err = store.GetByUID(journal.UID)
	if err != storage.ErrNotFound {
		t.Errorf("Expected ErrNotFound after delete, got %v", err)
	}
}

func TestJournalDeleteNonExistent(t *testing.T) {
	tmpDir := "/tmp/david-test-journals-delete-nonexistent-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	err = store.Delete("non-existent-uid")
	if err != storage.ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

func TestJournalList(t *testing.T) {
	tmpDir := "/tmp/david-test-journals-list-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	calendar1 := "calendar-123"
	calendar2 := "calendar-456"

	journal1 := &storage.Journal{
		UID:          "journal-list-1",
		CalendarUID:  calendar1,
		DTStamp:      time.Now().UTC(),
		Summary:      "Journal 1",
		Description:  "Description 1",
		Status:       "DRAFT",
		Sequence:     0,
		Created:      time.Time{},
		LastModified: time.Time{},
	}

	journal2 := &storage.Journal{
		UID:          "journal-list-2",
		CalendarUID:  calendar1,
		DTStamp:      time.Now().UTC(),
		Summary:      "Journal 2",
		Description:  "Description 2",
		Status:       "DRAFT",
		Sequence:     0,
		Created:      time.Time{},
		LastModified: time.Time{},
	}

	journal3 := &storage.Journal{
		UID:          "journal-list-3",
		CalendarUID:  calendar2,
		DTStamp:      time.Now().UTC(),
		Summary:      "Journal 3",
		Description:  "Description 3",
		Status:       "DRAFT",
		Sequence:     0,
		Created:      time.Time{},
		LastModified: time.Time{},
	}

	err = store.Create(journal1)
	if err != nil {
		t.Fatalf("Create journal1 failed: %v", err)
	}

	err = store.Create(journal2)
	if err != nil {
		t.Fatalf("Create journal2 failed: %v", err)
	}

	err = store.Create(journal3)
	if err != nil {
		t.Fatalf("Create journal3 failed: %v", err)
	}

	journals, err := store.List(calendar1, 100, 0)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(journals) != 2 {
		t.Errorf("Expected 2 journals for calendar1, got %d", len(journals))
	}

	journals, err = store.List(calendar2, 100, 0)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(journals) != 1 {
		t.Errorf("Expected 1 journal for calendar2, got %d", len(journals))
	}
}

func TestJournalListPagination(t *testing.T) {
	tmpDir := "/tmp/david-test-journals-pagination-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	calendarUID := "calendar-123"

	var journalUIDs []string
	for i := 0; i < 5; i++ {
		journal := &storage.Journal{
			UID:          fmt.Sprintf("journal-pagination-%d", i),
			CalendarUID:  calendarUID,
			DTStamp:      time.Now().UTC(),
			Summary:      fmt.Sprintf("Journal %d", i),
			Description:  "Description",
			Status:       "DRAFT",
			Sequence:     0,
			Created:      time.Time{},
			LastModified: time.Time{},
		}
		err := store.Create(journal)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}
		journalUIDs = append(journalUIDs, journal.UID)
	}

	journals, err := store.List(calendarUID, 2, 0)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(journals) != 2 {
		t.Errorf("Expected 2 journals, got %d", len(journals))
	}

	journals, err = store.List(calendarUID, 2, 2)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(journals) != 2 {
		t.Errorf("Expected 2 journals, got %d", len(journals))
	}

	journals, err = store.List(calendarUID, 2, 4)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(journals) != 1 {
		t.Errorf("Expected 1 journal, got %d", len(journals))
	}

	journals, err = store.List(calendarUID, 2, 5)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(journals) != 0 {
		t.Errorf("Expected 0 journals, got %d", len(journals))
	}
}

func TestJournalListDefaultLimit(t *testing.T) {
	tmpDir := "/tmp/david-test-journals-defaultlimit-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	calendarUID := "calendar-123"

	for i := 0; i < 3; i++ {
		journal := &storage.Journal{
			UID:          fmt.Sprintf("journal-defaultlimit-%d", i),
			CalendarUID:  calendarUID,
			DTStamp:      time.Now().UTC(),
			Summary:      fmt.Sprintf("Journal %d", i),
			Description:  "Description",
			Status:       "DRAFT",
			Sequence:     0,
			Created:      time.Time{},
			LastModified: time.Time{},
		}
		err := store.Create(journal)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}
	}

	journals, err := store.List(calendarUID, 0, 0)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(journals) != 3 {
		t.Errorf("Expected 3 journals with default limit, got %d", len(journals))
	}
}

func TestJournalListEmpty(t *testing.T) {
	tmpDir := "/tmp/david-test-journals-empty-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	journals, err := store.List("non-existent-calendar", 100, 0)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(journals) != 0 {
		t.Errorf("Expected 0 journals, got %d", len(journals))
	}
}

func TestJournalCreateGetUpdateDeleteCycle(t *testing.T) {
	tmpDir := "/tmp/david-test-journals-cycle-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	journal := &storage.Journal{
		UID:          "journal-cycle-test",
		CalendarUID:  "calendar-123",
		DTStamp:      time.Now().UTC(),
		Summary:      "Test Journal",
		Description:  "This is a test journal entry",
		Status:       "DRAFT",
		Sequence:     0,
		Created:      time.Time{},
		LastModified: time.Time{},
	}

	err = store.Create(journal)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	retrieved, err := store.GetByUID(journal.UID)
	if err != nil {
		t.Fatalf("GetByUID failed: %v", err)
	}
	if retrieved.Summary != "Test Journal" {
		t.Error("Summary mismatch after retrieval")
	}

	retrieved.Summary = "Modified Journal"
	retrieved.Description = "Modified description"
	retrieved.Status = "PUBLISHED"
	err = store.Update(journal.UID, retrieved)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	updated, err := store.GetByUID(journal.UID)
	if err != nil {
		t.Fatalf("GetByUID failed: %v", err)
	}
	if updated.Summary != "Modified Journal" {
		t.Errorf("Summary not updated: got %q", updated.Summary)
	}

	err = store.Delete(journal.UID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err = store.GetByUID(journal.UID)
	if err != storage.ErrNotFound {
		t.Errorf("Expected ErrNotFound after delete, got %v", err)
	}
}

func TestJournalListOrderConsistency(t *testing.T) {
	tmpDir := "/tmp/david-test-journals-consistency-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	calendarUID := "calendar-123"

	for i := 0; i < 10; i++ {
		journal := &storage.Journal{
			UID:          fmt.Sprintf("journal-consistency-%d", i),
			CalendarUID:  calendarUID,
			DTStamp:      time.Now().UTC(),
			Summary:      fmt.Sprintf("Journal %d", i),
			Description:  "Description",
			Status:       "DRAFT",
			Sequence:     0,
			Created:      time.Time{},
			LastModified: time.Time{},
		}
		err := store.Create(journal)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}
	}

	var firstList []*storage.Journal
	for iteration := 0; iteration < 3; iteration++ {
		journals, err := store.List(calendarUID, 100, 0)
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}

		if firstList == nil {
			firstList = journals
		} else {
			if len(journals) != len(firstList) {
				t.Errorf("Iteration %d: length mismatch, got %d, want %d", iteration, len(journals), len(firstList))
			}
		}
	}
}

func TestJournalWithDifferentStatuses(t *testing.T) {
	tmpDir := "/tmp/david-test-journals-statuses-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	journal := &storage.Journal{
		UID:          "journal-with-status",
		CalendarUID:  "calendar-123",
		DTStamp:      time.Now().UTC(),
		Summary:      "Journal with status",
		Description:  "Main description",
		Status:       "PUBLISHED",
		Sequence:     0,
		Created:      time.Time{},
		LastModified: time.Time{},
	}

	err = store.Create(journal)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	retrieved, err := store.GetByUID(journal.UID)
	if err != nil {
		t.Fatalf("GetByUID failed: %v", err)
	}

	if retrieved.Status != "PUBLISHED" {
		t.Errorf("Expected status PUBLISHED, got %q", retrieved.Status)
	}

	// Update to different status
	retrieved.Status = "DRAFT"
	err = store.Update(journal.UID, retrieved)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	updated, err := store.GetByUID(journal.UID)
	if err != nil {
		t.Fatalf("GetByUID failed: %v", err)
	}

	if updated.Status != "DRAFT" {
		t.Errorf("Expected status DRAFT, got %q", updated.Status)
	}
}
