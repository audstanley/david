package todos

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/audstanley/david/app/storage"
)

func TestTodoCreate(t *testing.T) {
	tmpDir := "/tmp/david-test-todos-create-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	todo := &storage.Todo{
		UID:          "todo-test-create",
		CalendarUID:  "calendar-123",
		DTStamp:      time.Now().UTC(),
		Due:          time.Now().UTC().Add(24 * time.Hour),
		Summary:      "Test Todo",
		Description:  "This is a test todo",
		Status:       "NEEDS-ACTION",
		Priority:     5,
		Percent:      0,
		Sequence:     0,
		Created:      time.Time{},
		LastModified: time.Time{},
		Completed:    time.Time{},
	}

	err = store.Create(todo)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if todo.UID == "" {
		t.Error("UID should be set after creation")
	}

	if todo.Sequence != 1 {
		t.Errorf("Expected sequence 1, got %d", todo.Sequence)
	}
}

func TestTodoCreateDuplicateUID(t *testing.T) {
	tmpDir := "/tmp/david-test-todos-duplicate-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	todo1 := &storage.Todo{
		UID:          "todo-duplicate-test",
		CalendarUID:  "calendar-123",
		DTStamp:      time.Now().UTC(),
		Due:          time.Now().UTC().Add(24 * time.Hour),
		Summary:      "Test Todo 1",
		Description:  "This is a test todo",
		Status:       "NEEDS-ACTION",
		Priority:     5,
		Percent:      0,
		Sequence:     0,
		Created:      time.Time{},
		LastModified: time.Time{},
		Completed:    time.Time{},
	}

	err = store.Create(todo1)
	if err != nil {
		t.Fatalf("First Create failed: %v", err)
	}

	todo2 := &storage.Todo{
		UID:          "todo-duplicate-test",
		CalendarUID:  "calendar-456",
		DTStamp:      time.Now().UTC(),
		Due:          time.Now().UTC().Add(24 * time.Hour),
		Summary:      "Test Todo 2",
		Description:  "This is a test todo",
		Status:       "NEEDS-ACTION",
		Priority:     5,
		Percent:      0,
		Sequence:     0,
		Created:      time.Time{},
		LastModified: time.Time{},
		Completed:    time.Time{},
	}

	err = store.Create(todo2)
	if err != storage.ErrAlreadyExists {
		t.Errorf("Expected ErrAlreadyExists, got %v", err)
	}
}

func TestTodoCreateWithAllFields(t *testing.T) {
	tmpDir := "/tmp/david-test-todos-fields-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	todo := &storage.Todo{
		UID:          "todo-custom-fields",
		CalendarUID:  "calendar-123",
		DTStamp:      time.Now().UTC(),
		Due:          time.Now().UTC().Add(48 * time.Hour),
		Summary:      "Complete Todo",
		Description:  "Full description with details",
		Status:       "COMPLETED",
		Priority:     1,
		Percent:      100,
		Sequence:     0,
		Created:      time.Time{},
		LastModified: time.Time{},
		Completed:    time.Now().UTC(),
	}

	err = store.Create(todo)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	retrieved, err := store.GetByUID(todo.UID)
	if err != nil {
		t.Fatalf("GetByUID failed: %v", err)
	}

	if retrieved.Summary != todo.Summary {
		t.Errorf("Summary mismatch: got %q, want %q", retrieved.Summary, todo.Summary)
	}
	if retrieved.Description != todo.Description {
		t.Errorf("Description mismatch: got %q, want %q", retrieved.Description, todo.Description)
	}
	if retrieved.Status != todo.Status {
		t.Errorf("Status mismatch: got %q, want %q", retrieved.Status, todo.Status)
	}
	if retrieved.Priority != todo.Priority {
		t.Errorf("Priority mismatch: got %d, want %d", retrieved.Priority, todo.Priority)
	}
	if retrieved.Percent != todo.Percent {
		t.Errorf("Percent mismatch: got %d, want %d", retrieved.Percent, todo.Percent)
	}
	if !retrieved.Due.Equal(todo.Due) {
		t.Errorf("Due mismatch: got %v, want %v", retrieved.Due, todo.Due)
	}
}

func TestTodoGetByUID(t *testing.T) {
	tmpDir := "/tmp/david-test-todos-getbyuid-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	todo := &storage.Todo{
		UID:          "todo-getbyuid-test",
		CalendarUID:  "calendar-123",
		DTStamp:      time.Now().UTC(),
		Due:          time.Now().UTC().Add(24 * time.Hour),
		Summary:      "Test Todo",
		Description:  "This is a test todo",
		Status:       "NEEDS-ACTION",
		Priority:     5,
		Percent:      0,
		Sequence:     0,
		Created:      time.Time{},
		LastModified: time.Time{},
		Completed:    time.Time{},
	}

	err = store.Create(todo)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	retrieved, err := store.GetByUID(todo.UID)
	if err != nil {
		t.Fatalf("GetByUID failed: %v", err)
	}

	if retrieved.UID != todo.UID {
		t.Errorf("UID mismatch: got %q, want %q", retrieved.UID, todo.UID)
	}
	if retrieved.CalendarUID != todo.CalendarUID {
		t.Errorf("CalendarUID mismatch: got %q, want %q", retrieved.CalendarUID, todo.CalendarUID)
	}
	if retrieved.Summary != todo.Summary {
		t.Errorf("Summary mismatch: got %q, want %q", retrieved.Summary, todo.Summary)
	}
}

func TestTodoGetByUIDNonExistent(t *testing.T) {
	tmpDir := "/tmp/david-test-todos-nonexistent-" + t.Name()
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

func TestTodoUpdate(t *testing.T) {
	tmpDir := "/tmp/david-test-todos-update-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	todo := &storage.Todo{
		UID:          "todo-update-test",
		CalendarUID:  "calendar-123",
		DTStamp:      time.Now().UTC(),
		Due:          time.Now().UTC().Add(24 * time.Hour),
		Summary:      "Test Todo",
		Description:  "This is a test todo",
		Status:       "NEEDS-ACTION",
		Priority:     5,
		Percent:      0,
		Sequence:     0,
		Created:      time.Time{},
		LastModified: time.Time{},
		Completed:    time.Time{},
	}

	err = store.Create(todo)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	todo.Summary = "Updated Todo"
	todo.Status = "IN-PROCESS"
	todo.Percent = 50

	err = store.Update(todo.UID, todo, 0)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	retrieved, err := store.GetByUID(todo.UID)
	if err != nil {
		t.Fatalf("GetByUID failed: %v", err)
	}

	if retrieved.Summary != "Updated Todo" {
		t.Errorf("Summary not updated: got %q", retrieved.Summary)
	}
	if retrieved.Status != "IN-PROCESS" {
		t.Errorf("Status not updated: got %q", retrieved.Status)
	}
	if retrieved.Percent != 50 {
		t.Errorf("Percent not updated: got %d", retrieved.Percent)
	}
}

func TestTodoUpdateSequence(t *testing.T) {
	tmpDir := "/tmp/david-test-todos-sequence-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	todo := &storage.Todo{
		UID:          "todo-sequence-test",
		CalendarUID:  "calendar-123",
		DTStamp:      time.Now().UTC(),
		Due:          time.Now().UTC().Add(24 * time.Hour),
		Summary:      "Test Todo",
		Description:  "This is a test todo",
		Status:       "NEEDS-ACTION",
		Priority:     5,
		Percent:      0,
		Sequence:     0,
		Created:      time.Time{},
		LastModified: time.Time{},
		Completed:    time.Time{},
	}

	err = store.Create(todo)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	err = store.Update(todo.UID, todo, 5)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	retrieved, err := store.GetByUID(todo.UID)
	if err != nil {
		t.Fatalf("GetByUID failed: %v", err)
	}

	if retrieved.Sequence != 5 {
		t.Errorf("Expected sequence 5, got %d", retrieved.Sequence)
	}
}

func TestTodoUpdateNonExistent(t *testing.T) {
	tmpDir := "/tmp/david-test-todos-update-nonexistent-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	todo := &storage.Todo{
		UID:          "todo-update-test",
		CalendarUID:  "calendar-123",
		DTStamp:      time.Now().UTC(),
		Due:          time.Now().UTC().Add(24 * time.Hour),
		Summary:      "Test Todo",
		Description:  "This is a test todo",
		Status:       "NEEDS-ACTION",
		Priority:     5,
		Percent:      0,
		Sequence:     0,
		Created:      time.Time{},
		LastModified: time.Time{},
		Completed:    time.Time{},
	}

	err = store.Update("non-existent-uid", todo, 0)
	if err != storage.ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

func TestTodoDelete(t *testing.T) {
	tmpDir := "/tmp/david-test-todos-delete-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	todo := &storage.Todo{
		UID:          "todo-delete-test",
		CalendarUID:  "calendar-123",
		DTStamp:      time.Now().UTC(),
		Due:          time.Now().UTC().Add(24 * time.Hour),
		Summary:      "Test Todo",
		Description:  "This is a test todo",
		Status:       "NEEDS-ACTION",
		Priority:     5,
		Percent:      0,
		Sequence:     0,
		Created:      time.Time{},
		LastModified: time.Time{},
		Completed:    time.Time{},
	}

	err = store.Create(todo)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	err = store.Delete(todo.UID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err = store.GetByUID(todo.UID)
	if err != storage.ErrNotFound {
		t.Errorf("Expected ErrNotFound after delete, got %v", err)
	}
}

func TestTodoDeleteNonExistent(t *testing.T) {
	tmpDir := "/tmp/david-test-todos-delete-nonexistent-" + t.Name()
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

func TestTodoList(t *testing.T) {
	tmpDir := "/tmp/david-test-todos-list-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	calendar1 := "calendar-123"
	calendar2 := "calendar-456"

	todo1 := &storage.Todo{
		UID:          "todo-list-1",
		CalendarUID:  calendar1,
		DTStamp:      time.Now().UTC(),
		Due:          time.Now().UTC().Add(24 * time.Hour),
		Summary:      "Todo 1",
		Description:  "Description 1",
		Status:       "NEEDS-ACTION",
		Priority:     5,
		Percent:      0,
		Sequence:     0,
		Created:      time.Time{},
		LastModified: time.Time{},
		Completed:    time.Time{},
	}

	todo2 := &storage.Todo{
		UID:          "todo-list-2",
		CalendarUID:  calendar1,
		DTStamp:      time.Now().UTC(),
		Due:          time.Now().UTC().Add(24 * time.Hour),
		Summary:      "Todo 2",
		Description:  "Description 2",
		Status:       "NEEDS-ACTION",
		Priority:     5,
		Percent:      0,
		Sequence:     0,
		Created:      time.Time{},
		LastModified: time.Time{},
		Completed:    time.Time{},
	}

	todo3 := &storage.Todo{
		UID:          "todo-list-3",
		CalendarUID:  calendar2,
		DTStamp:      time.Now().UTC(),
		Due:          time.Now().UTC().Add(24 * time.Hour),
		Summary:      "Todo 3",
		Description:  "Description 3",
		Status:       "NEEDS-ACTION",
		Priority:     5,
		Percent:      0,
		Sequence:     0,
		Created:      time.Time{},
		LastModified: time.Time{},
		Completed:    time.Time{},
	}

	err = store.Create(todo1)
	if err != nil {
		t.Fatalf("Create todo1 failed: %v", err)
	}

	err = store.Create(todo2)
	if err != nil {
		t.Fatalf("Create todo2 failed: %v", err)
	}

	err = store.Create(todo3)
	if err != nil {
		t.Fatalf("Create todo3 failed: %v", err)
	}

	todos, err := store.List(calendar1, 100, 0)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(todos) != 2 {
		t.Errorf("Expected 2 todos for calendar1, got %d", len(todos))
	}

	todos, err = store.List(calendar2, 100, 0)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(todos) != 1 {
		t.Errorf("Expected 1 todo for calendar2, got %d", len(todos))
	}
}

func TestTodoListPagination(t *testing.T) {
	tmpDir := "/tmp/david-test-todos-pagination-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	calendarUID := "calendar-123"

	var todoUIDs []string
	for i := 0; i < 5; i++ {
		todo := &storage.Todo{
			UID:          fmt.Sprintf("todo-pagination-%d", i),
			CalendarUID:  calendarUID,
			DTStamp:      time.Now().UTC(),
			Due:          time.Now().UTC().Add(24 * time.Hour),
			Summary:      fmt.Sprintf("Todo %d", i),
			Description:  "Description",
			Status:       "NEEDS-ACTION",
			Priority:     5,
			Percent:      0,
			Sequence:     0,
			Created:      time.Time{},
			LastModified: time.Time{},
			Completed:    time.Time{},
		}
		err := store.Create(todo)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}
		todoUIDs = append(todoUIDs, todo.UID)
	}

	todos, err := store.List(calendarUID, 2, 0)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(todos) != 2 {
		t.Errorf("Expected 2 todos, got %d", len(todos))
	}

	todos, err = store.List(calendarUID, 2, 2)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(todos) != 2 {
		t.Errorf("Expected 2 todos, got %d", len(todos))
	}

	todos, err = store.List(calendarUID, 2, 4)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(todos) != 1 {
		t.Errorf("Expected 1 todo, got %d", len(todos))
	}

	todos, err = store.List(calendarUID, 2, 5)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(todos) != 0 {
		t.Errorf("Expected 0 todos, got %d", len(todos))
	}
}

func TestTodoListDefaultLimit(t *testing.T) {
	tmpDir := "/tmp/david-test-todos-defaultlimit-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	calendarUID := "calendar-123"

	for i := 0; i < 3; i++ {
		todo := &storage.Todo{
			UID:          fmt.Sprintf("todo-defaultlimit-%d", i),
			CalendarUID:  calendarUID,
			DTStamp:      time.Now().UTC(),
			Due:          time.Now().UTC().Add(24 * time.Hour),
			Summary:      fmt.Sprintf("Todo %d", i),
			Description:  "Description",
			Status:       "NEEDS-ACTION",
			Priority:     5,
			Percent:      0,
			Sequence:     0,
			Created:      time.Time{},
			LastModified: time.Time{},
			Completed:    time.Time{},
		}
		err := store.Create(todo)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}
	}

	todos, err := store.List(calendarUID, 0, 0)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(todos) != 3 {
		t.Errorf("Expected 3 todos with default limit, got %d", len(todos))
	}
}

func TestTodoListEmpty(t *testing.T) {
	tmpDir := "/tmp/david-test-todos-empty-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	todos, err := store.List("non-existent-calendar", 100, 0)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(todos) != 0 {
		t.Errorf("Expected 0 todos, got %d", len(todos))
	}
}

func TestTodoCreateGetUpdateDeleteCycle(t *testing.T) {
	tmpDir := "/tmp/david-test-todos-cycle-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	todo := &storage.Todo{
		UID:          "todo-cycle-test",
		CalendarUID:  "calendar-123",
		DTStamp:      time.Now().UTC(),
		Due:          time.Now().UTC().Add(24 * time.Hour),
		Summary:      "Test Todo",
		Description:  "This is a test todo",
		Status:       "NEEDS-ACTION",
		Priority:     5,
		Percent:      0,
		Sequence:     0,
		Created:      time.Time{},
		LastModified: time.Time{},
		Completed:    time.Time{},
	}

	err = store.Create(todo)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	retrieved, err := store.GetByUID(todo.UID)
	if err != nil {
		t.Fatalf("GetByUID failed: %v", err)
	}
	if retrieved.Summary != "Test Todo" {
		t.Error("Summary mismatch after retrieval")
	}

	retrieved.Summary = "Modified Todo"
	err = store.Update(todo.UID, retrieved, 1)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	updated, err := store.GetByUID(todo.UID)
	if err != nil {
		t.Fatalf("GetByUID failed: %v", err)
	}
	if updated.Summary != "Modified Todo" {
		t.Errorf("Summary not updated: got %q", updated.Summary)
	}

	err = store.Delete(todo.UID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err = store.GetByUID(todo.UID)
	if err != storage.ErrNotFound {
		t.Errorf("Expected ErrNotFound after delete, got %v", err)
	}
}

func TestTodoListOrderConsistency(t *testing.T) {
	tmpDir := "/tmp/david-test-todos-consistency-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	store := New(db)

	calendarUID := "calendar-123"

	for i := 0; i < 10; i++ {
		todo := &storage.Todo{
			UID:          fmt.Sprintf("todo-consistency-%d", i),
			CalendarUID:  calendarUID,
			DTStamp:      time.Now().UTC(),
			Due:          time.Now().UTC().Add(24 * time.Hour),
			Summary:      fmt.Sprintf("Todo %d", i),
			Description:  "Description",
			Status:       "NEEDS-ACTION",
			Priority:     5,
			Percent:      0,
			Sequence:     0,
			Created:      time.Time{},
			LastModified: time.Time{},
			Completed:    time.Time{},
		}
		err := store.Create(todo)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}
	}

	var firstList []*storage.Todo
	for iteration := 0; iteration < 3; iteration++ {
		todos, err := store.List(calendarUID, 100, 0)
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}

		if firstList == nil {
			firstList = todos
		} else {
			if len(todos) != len(firstList) {
				t.Errorf("Iteration %d: length mismatch, got %d, want %d", iteration, len(todos), len(firstList))
			}
		}
	}
}
