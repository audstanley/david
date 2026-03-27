package freebusy

import (
	"os"
	"testing"
	"time"

	"github.com/audstanley/david/app/storage"
)

func TestFreeBusyStoreStore(t *testing.T) {
	tmpDir := "/tmp/david-test-fb-store-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	fb := &storage.FreeBusyBlock{
		UID:    "fb-test-1",
		Start:  time.Date(2026, 3, 25, 10, 0, 0, 0, time.UTC),
		End:    time.Date(2026, 3, 25, 11, 0, 0, 0, time.UTC),
		Status: "BUSY",
	}

	err = store.Store(fb)
	if err != nil {
		t.Fatalf("Store failed: %v", err)
	}

	retrieved, err := store.Get("fb-test-1")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if retrieved.UID != "fb-test-1" {
		t.Errorf("Expected UID 'fb-test-1', got %q", retrieved.UID)
	}
	if retrieved.Status != "BUSY" {
		t.Errorf("Expected status 'BUSY', got %q", retrieved.Status)
	}
}

func TestFreeBusyStoreStoreDuplicate(t *testing.T) {
	tmpDir := "/tmp/david-test-fb-duplicate-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	fb := &storage.FreeBusyBlock{
		UID:    "fb-duplicate",
		Start:  time.Date(2026, 3, 25, 10, 0, 0, 0, time.UTC),
		End:    time.Date(2026, 3, 25, 11, 0, 0, 0, time.UTC),
		Status: "BUSY",
	}

	err = store.Store(fb)
	if err != nil {
		t.Fatalf("First Store failed: %v", err)
	}

	// Store again with same UID (should succeed, overwrites)
	fb2 := &storage.FreeBusyBlock{
		UID:    "fb-duplicate",
		Start:  time.Date(2026, 3, 26, 10, 0, 0, 0, time.UTC),
		End:    time.Date(2026, 3, 26, 11, 0, 0, 0, time.UTC),
		Status: "FREE",
	}

	err = store.Store(fb2)
	if err != nil {
		t.Fatalf("Second Store failed: %v", err)
	}

	retrieved, err := store.Get("fb-duplicate")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	// Should have the second block's data
	if retrieved.Status != "FREE" {
		t.Errorf("Expected status 'FREE', got %q", retrieved.Status)
	}
}

func TestFreeBusyStoreGetNotFound(t *testing.T) {
	tmpDir := "/tmp/david-test-fb-get-notfound-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	_, err = store.Get("non-existent-fb")
	if err != storage.ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

func TestFreeBusyStoreList(t *testing.T) {
	tmpDir := "/tmp/david-test-fb-list-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	// Store multiple blocks
	fb1 := &storage.FreeBusyBlock{
		UID:    "fb-list-1",
		Start:  time.Date(2026, 3, 25, 10, 0, 0, 0, time.UTC),
		End:    time.Date(2026, 3, 25, 11, 0, 0, 0, time.UTC),
		Status: "BUSY",
	}

	fb2 := &storage.FreeBusyBlock{
		UID:    "fb-list-2",
		Start:  time.Date(2026, 3, 26, 14, 0, 0, 0, time.UTC),
		End:    time.Date(2026, 3, 26, 15, 0, 0, 0, time.UTC),
		Status: "FREE",
	}

	fb3 := &storage.FreeBusyBlock{
		UID:    "fb-list-3",
		Start:  time.Date(2026, 3, 27, 9, 0, 0, 0, time.UTC),
		End:    time.Date(2026, 3, 27, 10, 0, 0, 0, time.UTC),
		Status: "BUSY-UNAVAILABLE",
	}

	store.Store(fb1)
	store.Store(fb2)
	store.Store(fb3)

	blocks, err := store.List()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(blocks) != 3 {
		t.Errorf("Expected 3 blocks, got %d", len(blocks))
	}
}

func TestFreeBusyStoreDelete(t *testing.T) {
	tmpDir := "/tmp/david-test-fb-delete-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	fb := &storage.FreeBusyBlock{
		UID:    "fb-delete-test",
		Start:  time.Date(2026, 3, 25, 10, 0, 0, 0, time.UTC),
		End:    time.Date(2026, 3, 25, 11, 0, 0, 0, time.UTC),
		Status: "BUSY",
	}

	err = store.Store(fb)
	if err != nil {
		t.Fatalf("Store failed: %v", err)
	}

	err = store.Delete("fb-delete-test")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err = store.Get("fb-delete-test")
	if err != storage.ErrNotFound {
		t.Errorf("Expected ErrNotFound after delete, got %v", err)
	}
}

func TestFreeBusyStoreDeleteNonExistent(t *testing.T) {
	tmpDir := "/tmp/david-test-fb-delete-notfound-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	err = store.Delete("non-existent-fb")
	if err != storage.ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

func TestFreeBusyStoreCRUDCycle(t *testing.T) {
	tmpDir := "/tmp/david-test-fb-cycle-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	fb := &storage.FreeBusyBlock{
		UID:    "fb-cycle-test",
		Start:  time.Date(2026, 3, 25, 10, 0, 0, 0, time.UTC),
		End:    time.Date(2026, 3, 25, 11, 0, 0, 0, time.UTC),
		Status: "BUSY",
	}

	// Create
	err = store.Store(fb)
	if err != nil {
		t.Fatalf("Store failed: %v", err)
	}

	// Read
	retrieved, err := store.Get("fb-cycle-test")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if retrieved.Status != "BUSY" {
		t.Errorf("Expected status 'BUSY', got %q", retrieved.Status)
	}

	// List
	blocks, err := store.List()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(blocks) != 1 {
		t.Errorf("Expected 1 block in list, got %d", len(blocks))
	}

	// Delete
	err = store.Delete("fb-cycle-test")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify deleted
	_, err = store.Get("fb-cycle-test")
	if err != storage.ErrNotFound {
		t.Errorf("Expected ErrNotFound after delete, got %v", err)
	}
}

func TestFreeBusyStoreAllStatuses(t *testing.T) {
	tmpDir := "/tmp/david-test-fb-statuses-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	statuses := []string{"FREE", "BUSY", "BUSY-UNAVAILABLE"}

	for i, status := range statuses {
		fb := &storage.FreeBusyBlock{
			UID:    "fb-status-" + string(rune('0'+i)),
			Start:  time.Date(2026, 3, 25, 10+i, 0, 0, 0, time.UTC),
			End:    time.Date(2026, 3, 25, 11+i, 0, 0, 0, time.UTC),
			Status: status,
		}

		err = store.Store(fb)
		if err != nil {
			t.Fatalf("Store failed for status %s: %v", status, err)
		}
	}

	blocks, err := store.List()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(blocks) != 3 {
		t.Errorf("Expected 3 blocks, got %d", len(blocks))
	}

	// Verify each status
	found := make(map[string]bool)
	for _, block := range blocks {
		found[block.Status] = true
	}

	for _, status := range statuses {
		if !found[status] {
			t.Errorf("Expected to find status %q", status)
		}
	}
}

func TestFreeBusyStoreMultipleBlocks(t *testing.T) {
	tmpDir := "/tmp/david-test-fb-multiple-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	// Store 10 blocks
	for i := 0; i < 10; i++ {
		fb := &storage.FreeBusyBlock{
			UID:    "fb-multi-" + string(rune('0'+i)),
			Start:  time.Date(2026, 3, 25, 10+i, 0, 0, 0, time.UTC),
			End:    time.Date(2026, 3, 25, 11+i, 0, 0, 0, time.UTC),
			Status: "BUSY",
		}

		err = store.Store(fb)
		if err != nil {
			t.Fatalf("Store failed: %v", err)
		}
	}

	blocks, err := store.List()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(blocks) != 10 {
		t.Errorf("Expected 10 blocks, got %d", len(blocks))
	}
}

func TestFreeBusyStoreIndexManagement(t *testing.T) {
	tmpDir := "/tmp/david-test-fb-index-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	fb := &storage.FreeBusyBlock{
		UID:    "fb-index-test",
		Start:  time.Date(2026, 3, 25, 10, 0, 0, 0, time.UTC),
		End:    time.Date(2026, 3, 25, 11, 0, 0, 0, time.UTC),
		Status: "BUSY",
	}

	err = store.Store(fb)
	if err != nil {
		t.Fatalf("Store failed: %v", err)
	}

	// Check that the key exists
	exists, err := db.Exists([]byte("freebusy:fb-index-test"))
	if err != nil {
		t.Fatalf("Exists check failed: %v", err)
	}
	if !exists {
		t.Error("Free/busy key should exist")
	}
}
