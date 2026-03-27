package storage

import (
	"os"
	"testing"
)

func TestBatchDelete(t *testing.T) {
	tmpDir := "/tmp/david-test-batch-" + t.Name()
	defer os.RemoveAll(tmpDir)

	db, err := NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer db.Close()

	// Test 1: Put then Delete
	key1 := []byte("test-key-1")
	val1 := []byte("test-value-1")

	if err := db.Batch([][2][]byte{{key1, val1}}); err != nil {
		t.Fatalf("failed to put: %v", err)
	}

	data, err := db.Get(key1)
	if err != nil {
		t.Fatalf("failed to get after put: %v", err)
	}
	if string(data) != "test-value-1" {
		t.Fatalf("expected 'test-value-1', got '%s'", string(data))
	}

	if err := db.Batch([][2][]byte{{key1, nil}}); err != nil {
		t.Fatalf("failed to delete: %v", err)
	}

	_, err = db.Get(key1)
	if !IsNotFound(err) {
		t.Fatalf("expected NOT_FOUND error, got: %v", err)
	}

	// Test 2: Multiple keys with mixed put and delete
	key2 := []byte("test-key-2")
	key3 := []byte("test-key-3")
	val2 := []byte("test-value-2")
	val3 := []byte("test-value-3")

	if err := db.Batch([][2][]byte{
		{key2, val2},
		{key3, val3},
	}); err != nil {
		t.Fatalf("failed to batch put: %v", err)
	}

	if err := db.Batch([][2][]byte{
		{key2, nil},
		{key3, []byte("updated-value-3")},
	}); err != nil {
		t.Fatalf("failed to batch delete/put: %v", err)
	}

	_, err = db.Get(key2)
	if !IsNotFound(err) {
		t.Fatalf("expected key2 to be deleted, got: %v", err)
	}

	data3, err := db.Get(key3)
	if err != nil {
		t.Fatalf("failed to get key3: %v", err)
	}
	if string(data3) != "updated-value-3" {
		t.Fatalf("expected 'updated-value-3', got '%s'", string(data3))
	}

	// Test 3: Delete non-existent key should succeed
	key4 := []byte("test-key-4")
	if err := db.Batch([][2][]byte{{key4, nil}}); err != nil {
		t.Fatalf("failed to delete non-existent key: %v", err)
	}
}
