package users

import (
	"os"
	"testing"

	"github.com/audstanley/david/app/storage"
)

func TestUserStoreCreate(t *testing.T) {
	tmpDir := "/tmp/david-test-users-create-" + t.Name()
	defer os.RemoveAll(tmpDir)

	usersDB, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer usersDB.Close()

	store := New(usersDB)

	// Test 1: Create user with valid data
	userID := "test-user-create-valid"
	if err := store.Create(userID, "validuser", "passwordhash", "valid@example.com", "Valid User", "user"); err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	// Verify all indexes were created
	if exists, _ := usersDB.Exists([]byte(keyPrefixUser + userID)); !exists {
		t.Error("user data key not created")
	}
	if exists, _ := usersDB.Exists([]byte(keyPrefixUserByUID + userID)); !exists {
		t.Error("user by UID index not created")
	}
	if exists, _ := usersDB.Exists([]byte(keyPrefixUserByName + "validuser")); !exists {
		t.Error("user by name index not created")
	}

	// Verify user can be retrieved
	user, err := store.GetByID(userID)
	if err != nil {
		t.Fatalf("failed to get user: %v", err)
	}
	if user.Username != "validuser" {
		t.Errorf("expected username 'validuser', got '%s'", user.Username)
	}
	if user.Email != "valid@example.com" {
		t.Errorf("expected email 'valid@example.com', got '%s'", user.Email)
	}
	if user.DisplayName != "Valid User" {
		t.Errorf("expected display name 'Valid User', got '%s'", user.DisplayName)
	}
	if user.Role != "user" {
		t.Errorf("expected role 'user', got '%s'", user.Role)
	}
}

func TestUserStoreCreateDuplicateUsername(t *testing.T) {
	tmpDir := "/tmp/david-test-users-duplicate-" + t.Name()
	defer os.RemoveAll(tmpDir)

	usersDB, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer usersDB.Close()

	store := New(usersDB)

	// Create first user
	userID1 := "test-user-dup-1"
	if err := store.Create(userID1, "duplicateuser", "passwordhash1", "user1@example.com", "User 1", "user"); err != nil {
		t.Fatalf("failed to create first user: %v", err)
	}

	// Try to create second user with same username
	userID2 := "test-user-dup-2"
	err = store.Create(userID2, "duplicateuser", "passwordhash2", "user2@example.com", "User 2", "manager")
	if err == nil {
		t.Fatal("expected error when creating user with duplicate username, got nil")
	}
	if !storage.IsAlreadyExists(err) {
		t.Errorf("expected ALREADY_EXISTS error, got: %v", err)
	}

	// Verify only first user exists
	users, err := store.List(100, 0)
	if err != nil {
		t.Fatalf("failed to list users: %v", err)
	}
	if len(users) != 1 {
		t.Errorf("expected 1 user, got %d", len(users))
	}
}

func TestUserStoreGetByIDNotFound(t *testing.T) {
	tmpDir := "/tmp/david-test-users-get-notfound-" + t.Name()
	defer os.RemoveAll(tmpDir)

	usersDB, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer usersDB.Close()

	store := New(usersDB)

	// Try to get non-existent user
	_, err = store.GetByID("non-existent-user-id")
	if err == nil {
		t.Fatal("expected error when getting non-existent user, got nil")
	}
	if !storage.IsNotFound(err) {
		t.Errorf("expected NOT_FOUND error, got: %v", err)
	}
}

func TestUserStoreGetByUsername(t *testing.T) {
	tmpDir := "/tmp/david-test-users-get-byname-" + t.Name()
	defer os.RemoveAll(tmpDir)

	usersDB, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer usersDB.Close()

	store := New(usersDB)

	// Create user
	userID := "test-user-get-byname"
	expectedUsername := "searchableuser"
	if err := store.Create(userID, expectedUsername, "passwordhash", "search@example.com", "Search User", "admin"); err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	// Get user by username
	user, err := store.GetByUsername(expectedUsername)
	if err != nil {
		t.Fatalf("failed to get user by username: %v", err)
	}

	if user.ID != userID {
		t.Errorf("expected user ID '%s', got '%s'", userID, user.ID)
	}
	if user.Username != expectedUsername {
		t.Errorf("expected username '%s', got '%s'", expectedUsername, user.Username)
	}
}

func TestUserStoreGetByUsernameNotFound(t *testing.T) {
	tmpDir := "/tmp/david-test-users-get-byname-notfound-" + t.Name()
	defer os.RemoveAll(tmpDir)

	usersDB, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer usersDB.Close()

	store := New(usersDB)

	// Try to get user with non-existent username
	_, err = store.GetByUsername("non-existent-username")
	if err == nil {
		t.Fatal("expected error when getting non-existent username, got nil")
	}
	if !storage.IsNotFound(err) {
		t.Errorf("expected NOT_FOUND error, got: %v", err)
	}
}

func TestUserStoreUpdateFields(t *testing.T) {
	tmpDir := "/tmp/david-test-users-update-fields-" + t.Name()
	defer os.RemoveAll(tmpDir)

	usersDB, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer usersDB.Close()

	store := New(usersDB)

	// Create user
	userID := "test-user-update"
	if err := store.Create(userID, "originaluser", "oldhash", "original@example.com", "Original Name", "user"); err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	// Update multiple fields
	updates := map[string]interface{}{
		"username":    "newuser",
		"email":       "new@example.com",
		"displayName": "New Display Name",
		"role":        "admin",
	}

	if err := store.Update(userID, updates); err != nil {
		t.Fatalf("failed to update user: %v", err)
	}

	// Verify updates
	user, err := store.GetByID(userID)
	if err != nil {
		t.Fatalf("failed to get user after update: %v", err)
	}

	if user.Username != "newuser" {
		t.Errorf("expected username 'newuser', got '%s'", user.Username)
	}
	if user.Email != "new@example.com" {
		t.Errorf("expected email 'new@example.com', got '%s'", user.Email)
	}
	if user.DisplayName != "New Display Name" {
		t.Errorf("expected display name 'New Display Name', got '%s'", user.DisplayName)
	}
	if user.Role != "admin" {
		t.Errorf("expected role 'admin', got '%s'", user.Role)
	}

	// Verify old name index was deleted and new one was created
	if exists, _ := usersDB.Exists([]byte(keyPrefixUserByName + "originaluser")); exists {
		t.Error("old name index still exists after update")
	}
	if exists, _ := usersDB.Exists([]byte(keyPrefixUserByName + "newuser")); !exists {
		t.Error("new name index not created after update")
	}
}

func TestUserStoreUpdateOnlyPassword(t *testing.T) {
	tmpDir := "/tmp/david-test-users-update-password-" + t.Name()
	defer os.RemoveAll(tmpDir)

	usersDB, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer usersDB.Close()

	store := New(usersDB)

	// Create user
	userID := "test-user-update-password"
	if err := store.Create(userID, "passworduser", "oldhash", "pass@example.com", "Password User", "manager"); err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	// Update only password
	updates := map[string]interface{}{
		"passwordHash": "newhash",
	}

	if err := store.Update(userID, updates); err != nil {
		t.Fatalf("failed to update user password: %v", err)
	}

	// Verify password was updated, other fields unchanged
	user, err := store.GetByID(userID)
	if err != nil {
		t.Fatalf("failed to get user: %v", err)
	}

	if user.PasswordHash != "newhash" {
		t.Errorf("expected password hash 'newhash', got '%s'", user.PasswordHash)
	}
	if user.Username != "passworduser" {
		t.Errorf("username changed unexpectedly: %s", user.Username)
	}
	if user.Email != "pass@example.com" {
		t.Errorf("email changed unexpectedly: %s", user.Email)
	}
}

func TestUserStoreListPagination(t *testing.T) {
	tmpDir := "/tmp/david-test-users-list-" + t.Name()
	defer os.RemoveAll(tmpDir)

	usersDB, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer usersDB.Close()

	store := New(usersDB)

	// Create multiple users
	expectedUsers := 15
	for i := 0; i < expectedUsers; i++ {
		userID := "test-user-list-" + string(rune('0'+i))
		username := "user" + string(rune('0'+i))
		if err := store.Create(userID, username, "hash", "user"+string(rune('0'+i))+"@example.com", "User "+string(rune('0'+i)), "user"); err != nil {
			t.Fatalf("failed to create user %d: %v", i, err)
		}
	}

	// List with limit 5, offset 0
	users, err := store.List(5, 0)
	if err != nil {
		t.Fatalf("failed to list users: %v", err)
	}
	if len(users) != 5 {
		t.Errorf("expected 5 users with limit=5, offset=0, got %d", len(users))
	}

	// List with limit 5, offset 5
	users, err = store.List(5, 5)
	if err != nil {
		t.Fatalf("failed to list users: %v", err)
	}
	if len(users) != 5 {
		t.Errorf("expected 5 users with limit=5, offset=5, got %d", len(users))
	}

	// List with limit 10, offset 10
	users, err = store.List(10, 10)
	if err != nil {
		t.Fatalf("failed to list users: %v", err)
	}
	if len(users) != 5 {
		t.Errorf("expected 5 users with limit=10, offset=10 (last 5), got %d", len(users))
	}

	// List all users (limit > total)
	users, err = store.List(100, 0)
	if err != nil {
		t.Fatalf("failed to list all users: %v", err)
	}
	if len(users) != expectedUsers {
		t.Errorf("expected %d users, got %d", expectedUsers, len(users))
	}
}

func TestUserStoreListEmpty(t *testing.T) {
	tmpDir := "/tmp/david-test-users-list-empty-" + t.Name()
	defer os.RemoveAll(tmpDir)

	usersDB, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer usersDB.Close()

	store := New(usersDB)

	// List from empty database
	users, err := store.List(100, 0)
	if err != nil {
		t.Fatalf("failed to list users from empty database: %v", err)
	}
	if len(users) != 0 {
		t.Errorf("expected 0 users from empty database, got %d", len(users))
	}
}

func TestUserStoreListLimitValidation(t *testing.T) {
	tmpDir := "/tmp/david-test-users-list-limit-" + t.Name()
	defer os.RemoveAll(tmpDir)

	usersDB, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer usersDB.Close()

	store := New(usersDB)

	// Create 10 users
	for i := 0; i < 10; i++ {
		userID := "test-user-limit-" + string(rune('0'+i))
		if err := store.Create(userID, "user"+string(rune('0'+i)), "hash", "u"+string(rune('0'+i))+"@example.com", "User", "user"); err != nil {
			t.Fatalf("failed to create user: %v", err)
		}
	}

	// Test with limit 0 (should default to 100)
	users, err := store.List(0, 0)
	if err != nil {
		t.Fatalf("failed to list users: %v", err)
	}
	if len(users) != 10 {
		t.Errorf("expected 10 users with limit=0 (default), got %d", len(users))
	}

	// Test with limit > 1000 (should cap at 1000)
	users, err = store.List(2000, 0)
	if err != nil {
		t.Fatalf("failed to list users: %v", err)
	}
	if len(users) != 10 {
		t.Errorf("expected 10 users with limit=2000, got %d", len(users))
	}
}

func TestUserStoreUpdateNonExistent(t *testing.T) {
	tmpDir := "/tmp/david-test-users-update-notfound-" + t.Name()
	defer os.RemoveAll(tmpDir)

	usersDB, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer usersDB.Close()

	store := New(usersDB)

	// Try to update non-existent user
	updates := map[string]interface{}{
		"username": "newuser",
	}
	err = store.Update("non-existent-user", updates)
	if err == nil {
		t.Fatal("expected error when updating non-existent user, got nil")
	}
	if !storage.IsNotFound(err) {
		t.Errorf("expected NOT_FOUND error, got: %v", err)
	}
}

func TestUserStoreDeleteNonExistent(t *testing.T) {
	tmpDir := "/tmp/david-test-users-delete-notfound-" + t.Name()
	defer os.RemoveAll(tmpDir)

	usersDB, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer usersDB.Close()

	store := New(usersDB)

	// Try to delete non-existent user
	err = store.Delete("non-existent-user")
	if err == nil {
		t.Fatal("expected error when deleting non-existent user, got nil")
	}
	if !storage.IsNotFound(err) {
		t.Errorf("expected NOT_FOUND error, got: %v", err)
	}
}

func TestUserStoreVerifyPassword(t *testing.T) {
	tmpDir := "/tmp/david-test-users-verify-" + t.Name()
	defer os.RemoveAll(tmpDir)

	usersDB, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	defer usersDB.Close()

	store := New(usersDB)

	// Test HashPassword function
	hash, err := store.HashPassword("testpassword123")
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	// Verify correct password
	match, err := store.VerifyPassword(hash, "testpassword123")
	if err != nil {
		t.Fatalf("failed to verify password: %v", err)
	}
	if !match {
		t.Error("password verification failed for correct password")
	}

	// Verify wrong password
	match, err = store.VerifyPassword(hash, "wrongpassword")
	if err != nil {
		t.Logf("Expected error for wrong password: %v", err)
	}
	if match {
		t.Error("password verification succeeded for wrong password")
	}
}
