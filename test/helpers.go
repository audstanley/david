// Package test provides helper utilities for testing
package test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// SetupTest creates a temporary directory for test data
func SetupTest(t *testing.T) string {
	tmpDir, err := os.MkdirTemp("", "david-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	t.Cleanup(func() {
		os.RemoveAll(tmpDir)
	})
	return tmpDir
}

// CreateTestDatabase creates a test LevelDB database
func CreateTestDatabase(t *testing.T, baseDir, dbName string) string {
	dbPath := filepath.Join(baseDir, "data")
	if err := os.MkdirAll(dbPath, 0755); err != nil {
		t.Fatalf("Failed to create database directory: %v", err)
	}
	return filepath.Join(dbPath, dbName)
}

// LoadTestConfig loads the test configuration
func LoadTestConfig(t *testing.T) string {
	configPath := filepath.Join("fixtures", "config", "test.yaml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Skip("Test config file not found, skipping")
	}
	return configPath
}

// GetFixturePath returns the path to a test fixture
func GetFixturePath(filename string) string {
	return filepath.Join("fixtures", filename)
}

// MustReadFixture reads a fixture file, failing the test if it fails
func MustReadFixture(t *testing.T, filename string) []byte {
	data, err := os.ReadFile(filepath.Join("fixtures", filename))
	if err != nil {
		t.Fatalf("Failed to read fixture %s: %v", filename, err)
	}
	return data
}

// AssertEqualIgnoringWhitespace compares two strings ignoring whitespace differences
func AssertEqualIgnoringWhitespace(t *testing.T, expected, actual string) {
	expectedNormalized := normalizeWhitespace(expected)
	actualNormalized := normalizeWhitespace(actual)
	if expectedNormalized != actualNormalized {
		t.Errorf("Expected:\n%q\nGot:\n%q", expectedNormalized, actualNormalized)
	}
}

func normalizeWhitespace(s string) string {
	// Remove all whitespace for comparison
	result := make([]rune, 0, len(s))
	for _, r := range s {
		if r != ' ' && r != '\t' && r != '\n' && r != '\r' {
			result = append(result, r)
		}
	}
	return string(result)
}

// CreateTempFile creates a temporary file with given content
func CreateTempFile(t *testing.T, name, content string) string {
	tmpFile, err := os.CreateTemp("", name+"-*")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer tmpFile.Close()

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	return tmpFile.Name()
}

// TempDir creates a temporary directory
func TempDir(t *testing.T) string {
	tmpDir, err := os.MkdirTemp("", "david-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	t.Cleanup(func() {
		os.RemoveAll(tmpDir)
	})
	return tmpDir
}

// FormatDuration formats a duration for display
func FormatDuration(d int64) string {
	if d < 1000 {
		return fmt.Sprintf("%dµs", d)
	} else if d < 1000000 {
		return fmt.Sprintf("%dms", d/1000)
	}
	return fmt.Sprintf("%ds", d/1000000)
}
