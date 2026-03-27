package test

import (
	"testing"

	"github.com/audstanley/david/app/storage/users"
)

// Test10_1_1: Test user storage with real database
func TestUserStorageIntegration(t *testing.T) {
	skipIfShort(t)

	db, tmpDir := CreateTempStorage(t)
	defer CleanupStorage(tmpDir)
	defer db.Close()

	userStore := users.New(db)

	t.Run("Create user", func(t *testing.T) {
		userID := GenerateTestUserID()
		err := userStore.Create(userID, "testuser", "hash123", "test@example.com", "Test User", "user")
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}
	})

	t.Run("Get user by ID", func(t *testing.T) {
		userID := GenerateTestUserID()
		err := userStore.Create(userID, "testuser2", "hash456", "test2@example.com", "Test User 2", "user")
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}

		user, err := userStore.GetByID(userID)
		if err != nil {
			t.Fatalf("Failed to get user: %v", err)
		}

		if user.Username != "testuser2" {
			t.Errorf("Expected username 'testuser2', got '%s'", user.Username)
		}
	})

	t.Run("Get user by username", func(t *testing.T) {
		userID := GenerateTestUserID()
		err := userStore.Create(userID, "uniqueuser", "hash789", "unique@example.com", "Unique User", "admin")
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}

		user, err := userStore.GetByUsername("uniqueuser")
		if err != nil {
			t.Fatalf("Failed to get user by username: %v", err)
		}

		if user.ID != userID {
			t.Errorf("Expected user ID '%s', got '%s'", userID, user.ID)
		}
	})

	t.Run("Update user", func(t *testing.T) {
		userID := GenerateTestUserID()
		err := userStore.Create(userID, "updatetest", "hash000", "update@example.com", "Original", "user")
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}

		updates := map[string]interface{}{
			"email":       "newemail@example.com",
			"displayName": "Updated Name",
			"role":        "admin",
		}
		err = userStore.Update(userID, updates)
		if err != nil {
			t.Fatalf("Failed to update user: %v", err)
		}

		user, err := userStore.GetByID(userID)
		if err != nil {
			t.Fatalf("Failed to get user after update: %v", err)
		}

		if user.Email != "newemail@example.com" {
			t.Errorf("Expected email 'newemail@example.com', got '%s'", user.Email)
		}
		if user.Role != "admin" {
			t.Errorf("Expected role 'admin', got '%s'", user.Role)
		}
	})

	t.Run("List users with pagination", func(t *testing.T) {
		// Create multiple users with unique IDs
		userIDs := make([]string, 10)
		for i := 0; i < 10; i++ {
			userID := GenerateTestUserID() + string(rune('0'+i))
			userIDs[i] = userID
			err := userStore.Create(userID, "pagetest"+string(rune('0'+i)), "hash", "pagetest"+string(rune('0'+i))+"@example.com", "Pagetest", "user")
			if err != nil {
				t.Fatalf("Failed to create user %d: %v", i, err)
			}
		}

		// Test pagination
		limit := 5
		offset := 0
		userList, err := userStore.List(limit, offset)
		if err != nil {
			t.Fatalf("Failed to list users: %v", err)
		}

		if len(userList) != limit {
			t.Errorf("Expected %d users, got %d", limit, len(userList))
		}

		// Test with offset
		offset = 5
		userList, err = userStore.List(limit, offset)
		if err != nil {
			t.Fatalf("Failed to list users with offset: %v", err)
		}

		if len(userList) != limit {
			t.Errorf("Expected %d users with offset, got %d", limit, len(userList))
		}
	})

	t.Run("Delete user", func(t *testing.T) {
		userID := GenerateTestUserID()
		err := userStore.Create(userID, "deletetest", "hash", "delete@example.com", "Delete Test", "user")
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}

		err = userStore.Delete(userID)
		if err != nil {
			t.Fatalf("Failed to delete user: %v", err)
		}

		_, err = userStore.GetByID(userID)
		if err == nil {
			t.Error("Expected error getting deleted user, got nil")
		}
	})

	t.Run("Duplicate username error", func(t *testing.T) {
		userID1 := GenerateTestUserID()
		err := userStore.Create(userID1, "duplicateuser", "hash1", "dup1@example.com", "Dup 1", "user")
		if err != nil {
			t.Fatalf("Failed to create first user: %v", err)
		}

		userID2 := GenerateTestUserID()
		err = userStore.Create(userID2, "duplicateuser", "hash2", "dup2@example.com", "Dup 2", "user")
		if err == nil {
			t.Error("Expected error for duplicate username, got nil")
		}
	})
}

// skipIfShort skips the test if -short flag is used
func skipIfShort(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
}
