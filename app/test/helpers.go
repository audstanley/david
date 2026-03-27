// Package test provides helper functions and fixtures for integration tests
package test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/audstanley/david/app"
	"github.com/audstanley/david/app/config"
	"github.com/audstanley/david/app/storage"
	"github.com/audstanley/david/app/storage/calendars"
	"github.com/audstanley/david/app/storage/users"
)

// CreateTempStorage creates a temporary storage instance for testing
// The caller is responsible for cleaning up the directory
func CreateTempStorage(t *testing.T) (*storage.Storage, string) {
	tmpDir := t.TempDir()

	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	return db, tmpDir
}

// CreateTestConfig creates a test configuration
func CreateTestConfig() *config.Config {
	return &config.Config{
		Address: "localhost",
		Port:    "8080",
		Dir:     "/tmp/david-test",
		Users:   make(map[string]config.UserConfig),
		Log: config.LogConfig{
			Debug:      true,
			Production: false,
		},
	}
}

// GenerateTestUserID generates a unique test user ID
func GenerateTestUserID() string {
	return fmt.Sprintf("test-user-%d", os.Getpid())
}

// GenerateTestCalendarUID generates a unique test calendar UID
func GenerateTestCalendarUID() string {
	return fmt.Sprintf("test-calendar-%d", os.Getpid())
}

// GenerateTestEventUID generates a unique test event UID
func GenerateTestEventUID() string {
	return fmt.Sprintf("test-event-%d", os.Getpid())
}

// GenerateTestTodoUID generates a unique test todo UID
func GenerateTestTodoUID() string {
	return fmt.Sprintf("test-todo-%d", os.Getpid())
}

// GenerateTestJournalUID generates a unique test journal UID
func GenerateTestJournalUID() string {
	return fmt.Sprintf("test-journal-%d", os.Getpid())
}

// LoadICalFile loads an iCalendar file from the fixtures directory
func LoadICalFile(filename string) ([]byte, error) {
	path := filepath.Join("fixtures", filename)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read fixture file %s: %w", filename, err)
	}
	return data, nil
}

// CreateTestUser creates a test user with the given username using user store
func CreateTestUser(db *storage.Storage, username, displayName string) (*storage.User, error) {
	userStore := users.New(db)

	password := "testpassword123"
	hash := app.GenHash([]byte(password))

	userID := GenerateTestUserID()
	err := userStore.Create(userID, username, hash, fmt.Sprintf("%s@test.com", username), displayName, "user")
	if err != nil {
		return nil, fmt.Errorf("failed to create test user: %w", err)
	}

	user, err := userStore.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve test user: %w", err)
	}

	return user, nil
}

// CleanupStorage removes a test storage directory
func CleanupStorage(dir string) error {
	return os.RemoveAll(dir)
}

// MustGetenv returns the value of an environment variable or panics if not set
func MustGetenv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic(fmt.Sprintf("Environment variable %s is required", key))
	}
	return val
}

// SkipIfCI skips the test if running in CI
func SkipIfCI(t *testing.T) {
	if os.Getenv("CI") == "true" {
		t.Skip("Skipping test in CI environment")
	}
}

// WithTempDir runs a function with a temporary directory
func WithTempDir(t *testing.T, fn func(dir string) error) {
	tmpDir := t.TempDir()
	if err := fn(tmpDir); err != nil {
		t.Fatalf("Function failed: %v", err)
	}
}

// BatchCreateUsers creates multiple test users in a batch
func BatchCreateUsers(db *storage.Storage, userConfigs []struct {
	username    string
	displayName string
	role        string
}) ([]*storage.User, error) {
	userStore := users.New(db)
	var userObjects []*storage.User

	for _, config := range userConfigs {
		userID := GenerateTestUserID()
		password := "testpassword123"
		hash := app.GenHash([]byte(password))

		err := userStore.Create(userID, config.username, hash, fmt.Sprintf("%s@test.com", config.username), config.displayName, config.role)
		if err != nil {
			return nil, fmt.Errorf("failed to create user %s: %w", config.username, err)
		}

		user, err := userStore.GetByID(userID)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve user %s: %w", config.username, err)
		}

		userObjects = append(userObjects, user)
	}

	return userObjects, nil
}

// BatchCreateCalendars creates multiple test calendars in a batch
func BatchCreateCalendars(db *storage.Storage, ownerID string, calendarConfigs []struct {
	name   string
	public bool
}) ([]*storage.Calendar, error) {
	calendarStore := calendars.New(db)
	var calendars []*storage.Calendar

	for _, config := range calendarConfigs {
		calendarUID := GenerateTestCalendarUID()
		err := calendarStore.Create(calendarUID, ownerID, config.name, "Test calendar", "#ff0000", "America/New_York", config.public)
		if err != nil {
			return nil, fmt.Errorf("failed to create calendar %s: %w", config.name, err)
		}

		calendar, err := calendarStore.GetByUID(calendarUID)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve calendar %s: %w", config.name, err)
		}

		calendars = append(calendars, calendar)
	}

	return calendars, nil
}

// MustCreateTestUser creates a test user or panics
func MustCreateTestUser(t *testing.T, db *storage.Storage, username, displayName string) *storage.User {
	user, err := CreateTestUser(db, username, displayName)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	return user
}

// MustCreateTestCalendar creates a test calendar or panics
func MustCreateTestCalendar(t *testing.T, db *storage.Storage, ownerID, name, description string) *storage.Calendar {
	calendarUID := GenerateTestCalendarUID()
	calendarStore := calendars.New(db)

	err := calendarStore.Create(calendarUID, ownerID, name, description, "#ff0000", "America/New_York", false)
	if err != nil {
		t.Fatalf("Failed to create test calendar: %v", err)
	}

	calendar, err := calendarStore.GetByUID(calendarUID)
	if err != nil {
		t.Fatalf("Failed to retrieve test calendar: %v", err)
	}
	return calendar
}
