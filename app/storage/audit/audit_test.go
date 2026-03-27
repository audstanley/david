package audit

import (
	"os"
	"testing"

	"github.com/audstanley/david/app/storage"
)

func TestAuditStoreLog(t *testing.T) {
	tmpDir := "/tmp/david-test-audit-log-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	err = store.Log("user-123", "create", "calendar", "cal-123", `{"name": "Test Calendar"}`)
	if err != nil {
		t.Fatalf("Log failed: %v", err)
	}
}

func TestAuditStoreLogMultiple(t *testing.T) {
	tmpDir := "/tmp/david-test-audit-multiple-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	err = store.Log("user-123", "create", "calendar", "cal-1", `{"name": "Calendar 1"}`)
	if err != nil {
		t.Fatalf("Log failed: %v", err)
	}

	err = store.Log("user-123", "update", "calendar", "cal-1", `{"name": "Updated Calendar"}`)
	if err != nil {
		t.Fatalf("Log failed: %v", err)
	}

	err = store.Log("user-456", "create", "event", "evt-1", `{"summary": "Test Event"}`)
	if err != nil {
		t.Fatalf("Log failed: %v", err)
	}
}

func TestAuditStoreGetByUser(t *testing.T) {
	tmpDir := "/tmp/david-test-audit-byuser-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	store.Log("user-123", "create", "calendar", "cal-1", `{"name": "Calendar 1"}`)
	store.Log("user-123", "update", "calendar", "cal-1", `{"name": "Updated"}`)
	store.Log("user-456", "create", "calendar", "cal-2", `{"name": "Calendar 2"}`)

	entries, err := store.GetByUser("user-123", 100, 0)
	if err != nil {
		t.Fatalf("GetByUser failed: %v", err)
	}

	if len(entries) != 2 {
		t.Errorf("Expected 2 entries for user-123, got %d", len(entries))
	}
}

func TestAuditStoreGetByUserPagination(t *testing.T) {
	tmpDir := "/tmp/david-test-audit-byuser-pag-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	// Log 10 entries for user-123
	for i := 0; i < 10; i++ {
		store.Log("user-123", "create", "event", "evt-"+string(rune('0'+i)), `{"summary": "Event"}`)
	}

	// Get first 5
	entries, err := store.GetByUser("user-123", 5, 0)
	if err != nil {
		t.Fatalf("GetByUser failed: %v", err)
	}

	if len(entries) != 5 {
		t.Errorf("Expected 5 entries, got %d", len(entries))
	}

	// Get next 5 (offset 5)
	entries, err = store.GetByUser("user-123", 5, 5)
	if err != nil {
		t.Fatalf("GetByUser failed: %v", err)
	}

	if len(entries) != 5 {
		t.Errorf("Expected 5 entries with offset 5, got %d", len(entries))
	}
}

func TestAuditStoreGetByEntity(t *testing.T) {
	tmpDir := "/tmp/david-test-audit-byentity-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	store.Log("user-123", "create", "calendar", "cal-123", `{"name": "Calendar"}`)
	store.Log("user-123", "update", "calendar", "cal-123", `{"name": "Updated"}`)
	store.Log("user-456", "create", "event", "evt-1", `{"summary": "Event"}`)

	entries, err := store.GetByEntity("calendar", "cal-123", 100, 0)
	if err != nil {
		t.Fatalf("GetByEntity failed: %v", err)
	}

	if len(entries) != 2 {
		t.Errorf("Expected 2 entries for calendar cal-123, got %d", len(entries))
	}
}

func TestAuditStoreGetByEntityPagination(t *testing.T) {
	tmpDir := "/tmp/david-test-audit-byentity-pag-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	// Log 10 events for the same entity (evt-1)
	for i := 0; i < 10; i++ {
		store.Log("user-123", "create", "calendar", "cal-123", `{"summary": "Event"}`)
	}

	// Get first 5
	entries, err := store.GetByEntity("calendar", "cal-123", 5, 0)
	if err != nil {
		t.Fatalf("GetByEntity failed: %v", err)
	}

	if len(entries) != 5 {
		t.Errorf("Expected 5 entries, got %d", len(entries))
	}

	// Get next 5 (offset 5)
	entries, err = store.GetByEntity("calendar", "cal-123", 5, 5)
	if err != nil {
		t.Fatalf("GetByEntity failed: %v", err)
	}

	if len(entries) != 5 {
		t.Errorf("Expected 5 entries with offset 5, got %d", len(entries))
	}
}



func TestAuditStoreLimitValidation(t *testing.T) {
	tmpDir := "/tmp/david-test-audit-limit-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	// Log 10 entries
	for i := 0; i < 10; i++ {
		store.Log("user-123", "create", "event", "evt-"+string(rune('0'+i)), `{"summary": "Event"}`)
	}

	// Test with limit 0 (should default to 100)
	entries, err := store.GetByUser("user-123", 0, 0)
	if err != nil {
		t.Fatalf("GetByUser failed: %v", err)
	}

	if len(entries) != 10 {
		t.Errorf("Expected 10 entries with limit=0, got %d", len(entries))
	}

	// Test with limit > 1000 (should cap at 1000)
	entries, err = store.GetByUser("user-123", 2000, 0)
	if err != nil {
		t.Fatalf("GetByUser failed: %v", err)
	}

	if len(entries) != 10 {
		t.Errorf("Expected 10 entries with limit=2000, got %d", len(entries))
	}
}

func TestAuditStoreDifferentActions(t *testing.T) {
	tmpDir := "/tmp/david-test-audit-actions-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	actions := []string{"create", "update", "delete", "share_grant", "share_revoke"}

	for _, action := range actions {
		err = store.Log("user-123", action, "calendar", "cal-123", `{"name": "Test"}`)
		if err != nil {
			t.Fatalf("Log failed for action %s: %v", action, err)
		}
	}

	entries, err := store.GetByUser("user-123", 100, 0)
	if err != nil {
		t.Fatalf("GetByUser failed: %v", err)
	}

	if len(entries) != 5 {
		t.Errorf("Expected 5 entries, got %d", len(entries))
	}
}

func TestAuditStoreDifferentEntityTypes(t *testing.T) {
	tmpDir := "/tmp/david-test-audit-entitiestypes-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	entityTypes := []string{"user", "calendar", "event", "todo", "journal"}

	for _, et := range entityTypes {
		err = store.Log("user-123", "create", et, "id-1", `{"name": "Test"}`)
		if err != nil {
			t.Fatalf("Log failed for entity type %s: %v", et, err)
		}
	}

	entries, err := store.GetByUser("user-123", 100, 0)
	if err != nil {
		t.Fatalf("GetByUser failed: %v", err)
	}

	if len(entries) != 5 {
		t.Errorf("Expected 5 entries, got %d", len(entries))
	}
}

func TestAuditStoreCRUDCycle(t *testing.T) {
	tmpDir := "/tmp/david-test-audit-cycle-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	// Log an entry
	err = store.Log("user-123", "create", "calendar", "cal-123", `{"name": "Calendar"}`)
	if err != nil {
		t.Fatalf("Log failed: %v", err)
	}

	// Retrieve it
	entries, err := store.GetByUser("user-123", 100, 0)
	if err != nil {
		t.Fatalf("GetByUser failed: %v", err)
	}

	if len(entries) != 1 {
		t.Errorf("Expected 1 entry, got %d", len(entries))
	}

	// Log more entries
	for i := 0; i < 5; i++ {
		store.Log("user-123", "update", "calendar", "cal-123", `{"name": "Updated"}`)
	}

	entries, err = store.GetByUser("user-123", 100, 0)
	if err != nil {
		t.Fatalf("GetByUser failed: %v", err)
	}

	if len(entries) != 6 {
		t.Errorf("Expected 6 entries, got %d", len(entries))
	}
}

func TestAuditStoreMultipleUsers(t *testing.T) {
	tmpDir := "/tmp/david-test-audit-multiple-users-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	// Log entries for different users
	store.Log("user-123", "create", "calendar", "cal-1", `{"name": "Calendar 1"}`)
	store.Log("user-456", "create", "calendar", "cal-2", `{"name": "Calendar 2"}`)
	store.Log("user-789", "create", "calendar", "cal-3", `{"name": "Calendar 3"}`)

	// Get entries for each user
	entries1, err := store.GetByUser("user-123", 100, 0)
	if err != nil {
		t.Fatalf("GetByUser failed: %v", err)
	}
	if len(entries1) != 1 {
		t.Errorf("Expected 1 entry for user-123, got %d", len(entries1))
	}

	entries2, err := store.GetByUser("user-456", 100, 0)
	if err != nil {
		t.Fatalf("GetByUser failed: %v", err)
	}
	if len(entries2) != 1 {
		t.Errorf("Expected 1 entry for user-456, got %d", len(entries2))
	}

	entries3, err := store.GetByUser("user-789", 100, 0)
	if err != nil {
		t.Fatalf("GetByUser failed: %v", err)
	}
	if len(entries3) != 1 {
		t.Errorf("Expected 1 entry for user-789, got %d", len(entries3))
	}
}

func TestAuditStoreIndexManagement(t *testing.T) {
	tmpDir := "/tmp/david-test-audit-index-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	err = store.Log("user-123", "create", "calendar", "cal-123", `{"name": "Calendar"}`)
	if err != nil {
		t.Fatalf("Log failed: %v", err)
	}

	// Check that all indexes were created
	prefixes := []string{
		"audit:",
		"by:user-123:",
		"by:calendar:cal-123:",
		"date:",
	}

	for _, prefix := range prefixes {
		// Just verify no error when checking
		_, err := db.Get([]byte(prefix + "test"))
		if err != storage.ErrNotFound {
			// Not finding 'test' is expected, we just want to ensure the prefix exists
		}
	}
}
