package recurrence

import (
	"os"
	"testing"
	"time"

	"github.com/audstanley/david/app/storage"
)

func TestRecurrenceStoreStore(t *testing.T) {
	tmpDir := "/tmp/david-test-rec-store-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	eventUID := "test-event-recurrence"
	instances := []storage.RecurrenceInstance{
		{
			EventUID:    eventUID,
			InstanceUID: "instance-1",
			DTStart:     time.Date(2026, 3, 25, 10, 0, 0, 0, time.UTC),
			DTEnd:       time.Date(2026, 3, 25, 11, 0, 0, 0, time.UTC),
		},
		{
			EventUID:    eventUID,
			InstanceUID: "instance-2",
			DTStart:     time.Date(2026, 4, 1, 10, 0, 0, 0, time.UTC),
			DTEnd:       time.Date(2026, 4, 1, 11, 0, 0, 0, time.UTC),
		},
	}

	err = store.Store(eventUID, instances)
	if err != nil {
		t.Fatalf("Store failed: %v", err)
	}

	retrieved, err := store.Get(eventUID, time.Time{}, time.Time{})
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if len(retrieved) != 2 {
		t.Errorf("Expected 2 instances, got %d", len(retrieved))
	}
}

func TestRecurrenceStoreStoreEmpty(t *testing.T) {
	tmpDir := "/tmp/david-test-rec-store-empty-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	eventUID := "test-event-empty-rec"
	instances := []storage.RecurrenceInstance{}

	err = store.Store(eventUID, instances)
	if err != nil {
		t.Fatalf("Store failed: %v", err)
	}
}

func TestRecurrenceStoreStoreUpdate(t *testing.T) {
	tmpDir := "/tmp/david-test-rec-store-update-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	eventUID := "test-event-update-rec"
	instances1 := []storage.RecurrenceInstance{
		{
			EventUID:    eventUID,
			InstanceUID: "instance-1",
			DTStart:     time.Date(2026, 3, 25, 10, 0, 0, 0, time.UTC),
			DTEnd:       time.Date(2026, 3, 25, 11, 0, 0, 0, time.UTC),
		},
	}

	err = store.Store(eventUID, instances1)
	if err != nil {
		t.Fatalf("First Store failed: %v", err)
	}

	instances2 := []storage.RecurrenceInstance{
		{
			EventUID:    eventUID,
			InstanceUID: "instance-2",
			DTStart:     time.Date(2026, 4, 1, 10, 0, 0, 0, time.UTC),
			DTEnd:       time.Date(2026, 4, 1, 11, 0, 0, 0, time.UTC),
		},
		{
			EventUID:    eventUID,
			InstanceUID: "instance-3",
			DTStart:     time.Date(2026, 4, 8, 10, 0, 0, 0, time.UTC),
			DTEnd:       time.Date(2026, 4, 8, 11, 0, 0, 0, time.UTC),
		},
	}

	err = store.Store(eventUID, instances2)
	if err != nil {
		t.Fatalf("Second Store failed: %v", err)
	}

	retrieved, err := store.Get(eventUID, time.Time{}, time.Time{})
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if len(retrieved) != 2 {
		t.Errorf("Expected 2 instances after update, got %d", len(retrieved))
	}

	// Verify old instance was removed
	for _, inst := range retrieved {
		if inst.InstanceUID == "instance-1" {
			t.Error("Old instance should have been removed")
		}
	}
}

func TestRecurrenceStoreGetWithDateRange(t *testing.T) {
	tmpDir := "/tmp/david-test-rec-get-range-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	eventUID := "test-event-date-range"
	instances := []storage.RecurrenceInstance{
		{
			EventUID:    eventUID,
			InstanceUID: "instance-1",
			DTStart:     time.Date(2026, 3, 25, 10, 0, 0, 0, time.UTC),
			DTEnd:       time.Date(2026, 3, 25, 11, 0, 0, 0, time.UTC),
		},
		{
			EventUID:    eventUID,
			InstanceUID: "instance-2",
			DTStart:     time.Date(2026, 4, 1, 10, 0, 0, 0, time.UTC),
			DTEnd:       time.Date(2026, 4, 1, 11, 0, 0, 0, time.UTC),
		},
		{
			EventUID:    eventUID,
			InstanceUID: "instance-3",
			DTStart:     time.Date(2026, 5, 1, 10, 0, 0, 0, time.UTC),
			DTEnd:       time.Date(2026, 5, 1, 11, 0, 0, 0, time.UTC),
		},
	}

	err = store.Store(eventUID, instances)
	if err != nil {
		t.Fatalf("Store failed: %v", err)
	}

	// Get instances in April only
	start := time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 4, 30, 23, 59, 59, 0, time.UTC)

	retrieved, err := store.Get(eventUID, start, end)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if len(retrieved) != 1 {
		t.Errorf("Expected 1 instance in April, got %d", len(retrieved))
	}

	if retrieved[0].InstanceUID != "instance-2" {
		t.Errorf("Expected instance-2, got %s", retrieved[0].InstanceUID)
	}
}

func TestRecurrenceStoreGetNotFound(t *testing.T) {
	tmpDir := "/tmp/david-test-rec-get-notfound-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	_, err = store.Get("non-existent-event", time.Time{}, time.Time{})
	if err != storage.ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

func TestRecurrenceStoreDelete(t *testing.T) {
	tmpDir := "/tmp/david-test-rec-delete-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	eventUID := "test-event-delete-rec"
	instances := []storage.RecurrenceInstance{
		{
			EventUID:    eventUID,
			InstanceUID: "instance-1",
			DTStart:     time.Date(2026, 3, 25, 10, 0, 0, 0, time.UTC),
			DTEnd:       time.Date(2026, 3, 25, 11, 0, 0, 0, time.UTC),
		},
	}

	err = store.Store(eventUID, instances)
	if err != nil {
		t.Fatalf("Store failed: %v", err)
	}

	err = store.Delete(eventUID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err = store.Get(eventUID, time.Time{}, time.Time{})
	if err != storage.ErrNotFound {
		t.Errorf("Expected ErrNotFound after delete, got %v", err)
	}
}

func TestRecurrenceStoreDeleteNonExistent(t *testing.T) {
	tmpDir := "/tmp/david-test-rec-delete-notfound-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	err = store.Delete("non-existent-event")
	if err != nil {
		t.Errorf("Expected no error for non-existent event, got %v", err)
	}
}

func TestRecurrenceStoreDeleteWithMultipleInstances(t *testing.T) {
	tmpDir := "/tmp/david-test-rec-delete-multi-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	eventUID := "test-event-delete-multi"
	instances := []storage.RecurrenceInstance{
		{
			EventUID:    eventUID,
			InstanceUID: "instance-1",
			DTStart:     time.Date(2026, 3, 25, 10, 0, 0, 0, time.UTC),
			DTEnd:       time.Date(2026, 3, 25, 11, 0, 0, 0, time.UTC),
		},
		{
			EventUID:    eventUID,
			InstanceUID: "instance-2",
			DTStart:     time.Date(2026, 4, 1, 10, 0, 0, 0, time.UTC),
			DTEnd:       time.Date(2026, 4, 1, 11, 0, 0, 0, time.UTC),
		},
		{
			EventUID:    eventUID,
			InstanceUID: "instance-3",
			DTStart:     time.Date(2026, 5, 1, 10, 0, 0, 0, time.UTC),
			DTEnd:       time.Date(2026, 5, 1, 11, 0, 0, 0, time.UTC),
		},
	}

	err = store.Store(eventUID, instances)
	if err != nil {
		t.Fatalf("Store failed: %v", err)
	}

	err = store.Delete(eventUID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err = store.Get(eventUID, time.Time{}, time.Time{})
	if err != storage.ErrNotFound {
		t.Errorf("Expected ErrNotFound after delete, got %v", err)
	}
}

func TestRecurrenceStoreCRUDCycle(t *testing.T) {
	tmpDir := "/tmp/david-test-rec-cycle-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	eventUID := "test-event-cycle"
	instances := []storage.RecurrenceInstance{
		{
			EventUID:    eventUID,
			InstanceUID: "instance-1",
			DTStart:     time.Date(2026, 3, 25, 10, 0, 0, 0, time.UTC),
			DTEnd:       time.Date(2026, 3, 25, 11, 0, 0, 0, time.UTC),
		},
	}

	err = store.Store(eventUID, instances)
	if err != nil {
		t.Fatalf("Store failed: %v", err)
	}

	retrieved, err := store.Get(eventUID, time.Time{}, time.Time{})
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if len(retrieved) != 1 {
		t.Errorf("Expected 1 instance, got %d", len(retrieved))
	}

	updatedInstances := []storage.RecurrenceInstance{
		{
			EventUID:    eventUID,
			InstanceUID: "instance-2",
			DTStart:     time.Date(2026, 4, 1, 10, 0, 0, 0, time.UTC),
			DTEnd:       time.Date(2026, 4, 1, 11, 0, 0, 0, time.UTC),
		},
	}

	err = store.Store(eventUID, updatedInstances)
	if err != nil {
		t.Fatalf("Update Store failed: %v", err)
	}

	updated, err := store.Get(eventUID, time.Time{}, time.Time{})
	if err != nil {
		t.Fatalf("Get after update failed: %v", err)
	}

	if len(updated) != 1 {
		t.Errorf("Expected 1 instance after update, got %d", len(updated))
	}

	err = store.Delete(eventUID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err = store.Get(eventUID, time.Time{}, time.Time{})
	if err != storage.ErrNotFound {
		t.Errorf("Expected ErrNotFound after delete, got %v", err)
	}
}

func TestRecurrenceStoreIndexManagement(t *testing.T) {
	tmpDir := "/tmp/david-test-rec-index-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	eventUID := "test-event-index"
	instances := []storage.RecurrenceInstance{
		{
			EventUID:    eventUID,
			InstanceUID: "instance-1",
			DTStart:     time.Date(2026, 3, 25, 10, 0, 0, 0, time.UTC),
			DTEnd:       time.Date(2026, 3, 25, 11, 0, 0, 0, time.UTC),
		},
	}

	err = store.Store(eventUID, instances)
	if err != nil {
		t.Fatalf("Store failed: %v", err)
	}

	// Check primary index
	exists, err := db.Exists([]byte("recurrence:instance-1"))
	if err != nil {
		t.Fatalf("Exists check failed: %v", err)
	}
	if !exists {
		t.Error("Recurrence instance key should exist")
	}

	// Check by-event index
	exists, err = db.Exists([]byte("by:test-event-index:instance-1"))
	if err != nil {
		t.Fatalf("Exists check failed: %v", err)
	}
	if !exists {
		t.Error("By-event index key should exist")
	}
}

func TestRecurrenceStoreMultipleEvents(t *testing.T) {
	tmpDir := "/tmp/david-test-rec-multi-events-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	eventUID1 := "test-event-1"
	eventUID2 := "test-event-2"

	instances1 := []storage.RecurrenceInstance{
		{
			EventUID:    eventUID1,
			InstanceUID: "instance-1",
			DTStart:     time.Date(2026, 3, 25, 10, 0, 0, 0, time.UTC),
			DTEnd:       time.Date(2026, 3, 25, 11, 0, 0, 0, time.UTC),
		},
	}

	instances2 := []storage.RecurrenceInstance{
		{
			EventUID:    eventUID2,
			InstanceUID: "instance-2",
			DTStart:     time.Date(2026, 4, 1, 10, 0, 0, 0, time.UTC),
			DTEnd:       time.Date(2026, 4, 1, 11, 0, 0, 0, time.UTC),
		},
	}

	err = store.Store(eventUID1, instances1)
	if err != nil {
		t.Fatalf("Store event1 failed: %v", err)
	}

	err = store.Store(eventUID2, instances2)
	if err != nil {
		t.Fatalf("Store event2 failed: %v", err)
	}

	retrieved1, err := store.Get(eventUID1, time.Time{}, time.Time{})
	if err != nil {
		t.Fatalf("Get event1 failed: %v", err)
	}

	if len(retrieved1) != 1 {
		t.Errorf("Expected 1 instance for event1, got %d", len(retrieved1))
	}

	retrieved2, err := store.Get(eventUID2, time.Time{}, time.Time{})
	if err != nil {
		t.Fatalf("Get event2 failed: %v", err)
	}

	if len(retrieved2) != 1 {
		t.Errorf("Expected 1 instance for event2, got %d", len(retrieved2))
	}
}
