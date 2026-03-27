package timezones

import (
	"os"
	"testing"

	"github.com/audstanley/david/app/icalendar"
	"github.com/audstanley/david/app/storage"
)

func TestTimeZoneStoreStore(t *testing.T) {
	tmpDir := "/tmp/david-test-tz-store-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	tzID := "America/New_York"
	vtimezone := &storage.VTimeZone{
		TZID:       tzID,
		Components: []icalendar.Component{},
	}

	err = store.Store(tzID, vtimezone)
	if err != nil {
		t.Fatalf("Store failed: %v", err)
	}

	retrieved, err := store.Get(tzID)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if retrieved.TZID != tzID {
		t.Errorf("Expected TZID %q, got %q", tzID, retrieved.TZID)
	}
}

func TestTimeZoneStoreStoreDuplicate(t *testing.T) {
	tmpDir := "/tmp/david-test-tz-duplicate-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	tzID := "America/Chicago"
	vtimezone := &storage.VTimeZone{
		TZID:       tzID,
		Components: []icalendar.Component{},
	}

	err = store.Store(tzID, vtimezone)
	if err != nil {
		t.Fatalf("First Store failed: %v", err)
	}

	err = store.Store(tzID, vtimezone)
	if err != storage.ErrAlreadyExists {
		t.Errorf("Expected ErrAlreadyExists, got %v", err)
	}
}

func TestTimeZoneStoreStoreWithComponents(t *testing.T) {
	tmpDir := "/tmp/david-test-tz-components-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	tzID := "Europe/London"
	vtimezone := &storage.VTimeZone{
		TZID: tzID,
		Components: []icalendar.Component{
			{Name: "VTIMEZONE"},
			{Name: "STANDARD"},
			{Name: "DAYLIGHT"},
		},
	}

	err = store.Store(tzID, vtimezone)
	if err != nil {
		t.Fatalf("Store failed: %v", err)
	}

	retrieved, err := store.Get(tzID)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if len(retrieved.Components) != 3 {
		t.Errorf("Expected 3 components, got %d", len(retrieved.Components))
	}
}

func TestTimeZoneStoreGetNotFound(t *testing.T) {
	tmpDir := "/tmp/david-test-tz-get-notfound-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	_, err = store.Get("NonExistent/TZ")
	if err != storage.ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

func TestTimeZoneStoreUpdate(t *testing.T) {
	tmpDir := "/tmp/david-test-tz-update-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	tzID := "Test/TZ"
	vtimezone := &storage.VTimeZone{
		TZID:       tzID,
		Components: []icalendar.Component{{Name: "VTIMEZONE"}},
	}

	err = store.Store(tzID, vtimezone)
	if err != nil {
		t.Fatalf("Store failed: %v", err)
	}

	vtimezone.TZID = "Test/TZ-Updated"
	vtimezone.Components = append(vtimezone.Components, icalendar.Component{Name: "STANDARD"})

	err = store.Update(tzID, vtimezone)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	retrieved, err := store.Get(tzID)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if len(retrieved.Components) != 2 {
		t.Errorf("Expected 2 components after update, got %d", len(retrieved.Components))
	}
}

func TestTimeZoneStoreUpdateNonExistent(t *testing.T) {
	tmpDir := "/tmp/david-test-tz-update-notfound-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	vtimezone := &storage.VTimeZone{
		TZID:       "NonExistent/TZ",
		Components: []icalendar.Component{},
	}

	err = store.Update("NonExistent/TZ", vtimezone)
	if err != storage.ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

func TestTimeZoneStoreDelete(t *testing.T) {
	tmpDir := "/tmp/david-test-tz-delete-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	tzID := "Test/Delete-TZ"
	vtimezone := &storage.VTimeZone{
		TZID:       tzID,
		Components: []icalendar.Component{},
	}

	err = store.Store(tzID, vtimezone)
	if err != nil {
		t.Fatalf("Store failed: %v", err)
	}

	err = store.Delete(tzID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err = store.Get(tzID)
	if err != storage.ErrNotFound {
		t.Errorf("Expected ErrNotFound after delete, got %v", err)
	}
}

func TestTimeZoneStoreDeleteNonExistent(t *testing.T) {
	tmpDir := "/tmp/david-test-tz-delete-notfound-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	err = store.Delete("NonExistent/TZ")
	if err != storage.ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

func TestTimeZoneStoreCRUDCycle(t *testing.T) {
	tmpDir := "/tmp/david-test-tz-cycle-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	tzID := "Test/Cycle-TZ"
	vtimezone := &storage.VTimeZone{
		TZID:       tzID,
		Components: []icalendar.Component{{Name: "VTIMEZONE"}},
	}

	err = store.Store(tzID, vtimezone)
	if err != nil {
		t.Fatalf("Store failed: %v", err)
	}

	retrieved, err := store.Get(tzID)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if retrieved.TZID != tzID {
		t.Errorf("TZID mismatch after get")
	}

	vtimezone.Components = append(vtimezone.Components, icalendar.Component{Name: "STANDARD"})
	err = store.Update(tzID, vtimezone)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	updated, err := store.Get(tzID)
	if err != nil {
		t.Fatalf("Get failed after update: %v", err)
	}

	if len(updated.Components) != 2 {
		t.Errorf("Expected 2 components after update, got %d", len(updated.Components))
	}

	err = store.Delete(tzID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err = store.Get(tzID)
	if err != storage.ErrNotFound {
		t.Errorf("Expected ErrNotFound after delete, got %v", err)
	}
}

func TestTimeZoneStoreIndexManagement(t *testing.T) {
	tmpDir := "/tmp/david-test-tz-index-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	tzID := "Index/Test-TZ"
	vtimezone := &storage.VTimeZone{
		TZID:       tzID,
		Components: []icalendar.Component{},
	}

	err = store.Store(tzID, vtimezone)
	if err != nil {
		t.Fatalf("Store failed: %v", err)
	}

	exists, err := db.Exists([]byte("tz:" + tzID))
	if err != nil {
		t.Fatalf("Exists check failed: %v", err)
	}
	if !exists {
		t.Error("Timezone data key should exist")
	}

	exists, err = db.Exists([]byte("tzid:" + tzID))
	if err != nil {
		t.Fatalf("Exists check failed: %v", err)
	}
	if !exists {
		t.Error("Timezone ID index key should exist")
	}
}
