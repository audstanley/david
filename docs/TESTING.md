# Testing Guide

## Overview

david uses comprehensive testing with 85%+ coverage target. This guide explains how to run tests and contribute tests.

## Running Tests

### All Tests

```bash
# Run all tests
go test ./...

# Run with race detection (recommended)
go test ./... -race

# Run all tests with coverage
mage Coverage
```

### Specific Package

```bash
go test ./app/icalendar/... -v -cover
go test ./app/storage/... -v -cover
go test ./app/auth/... -v -cover
```

### Single Test

```bash
go test -v -run TestParser_ParseEvent ./app/icalendar/...
```

### Coverage Report

```bash
# Generate coverage HTML
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html

# Show coverage by file
go tool cover -func=coverage.out
```

## Test Structure

### Unit Tests

Unit tests are table-driven and test individual components in isolation.

**Location:** `*_test.go` files alongside source code

**Example:**
```go
func TestParser_ParseEvent(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    *Event
        wantErr bool
    }{
        {
            name:  "valid event",
            input: `BEGIN:VEVENT...`,
            want:  &Event{...},
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := ParseEvent(tt.input)
            // assertions...
        })
    }
}
```

### Integration Tests

Integration tests test components working together with real dependencies.

**Location:** `*_integration_test.go` files

**Example:**
```go
func TestEventCRUD(t *testing.T) {
    // Setup test database
    db := SetupTestDB(t)
    
    // Create event
    event := &Event{UID: "test-1", ...}
    err := db.Events.Create(event)
    assert.NoError(t, err)
    
    // Read event back
    got, err := db.Events.GetByID("test-1")
    assert.NoError(t, err)
    assert.Equal(t, event, got)
}
```

### Golden File Tests

Golden file tests compare parser output against expected results.

**Location:** `test/fixtures/` directory

**Usage:**
```go
func TestParser_EventFixture(t *testing.T) {
    // Read fixture
    input := MustReadFixture(t, "event_basic.ics")
    
    // Parse
    calendar, err := Parse(string(input))
    assert.NoError(t, err)
    
    // Generate back
    output, err := Generate(calendar)
    assert.NoError(t, err)
    
    // Compare (ignoring whitespace)
    AssertEqualIgnoringWhitespace(t, expectedOutput, output)
}
```

### Fuzz Tests

Fuzz tests find crashes and panics by feeding random input.

**Location:** `fuzz_test.go`

**Example:**
```go
func FuzzParser(f *testing.F) {
    // Add seed corpus
    f.Add("BEGIN:VCALENDAR\nVERSION:2.0\n...")
    
    f.Fuzz(func(t *testing.T, data string) {
        // This should not panic
        _, _ = Parse(data)
    })
}
```

## Test Fixtures

### ICS Fixtures

Located in `test/fixtures/ics/`:

- `event_basic.ics` - Simple event
- `event_recurring.ics` - Recurring event with RRULE
- `event_timezone.ics` - Event with timezone
- `event_alarm.ics` - Event with alarms
- `todo_basic.ics` - Basic todo
- `todo_completed.ics` - Completed todo
- `journal_basic.ics` - Basic journal
- `freebusy_basic.ics` - Free/busy block
- `timezone_america_new_york.ics` - Timezone definition
- `malformed.ics` - Invalid ICS for error testing

### Config Fixtures

Located in `test/fixtures/config/`:

- `test.yaml` - Test configuration

## Writing Tests

### Guidelines

1. **Use table-driven tests** for clarity and coverage
2. **Name tests clearly**: `TestFunction_TestCaseName`
3. **Test both success and error cases**
4. **Use testify for assertions** when appropriate
5. **Keep tests fast** - avoid slow I/O
6. **Use t.Cleanup()** for cleanup
7. **Document edge cases** with test comments

### Test Helpers

Use `test/helpers.go`:

```go
func TestExample(t *testing.T) {
    // Setup temp directory
    tmpDir := test.TempDir(t)
    
    // Load config
    config := config.LoadTestConfig(t)
    
    // Read fixture
    data := test.MustReadFixture(t, "event_basic.ics")
}
```

## Coverage Goals

| Component | Target |
|-----------|--------|
| iCalendar parser | 95% |
| Storage layer | 90% |
| Auth system | 90% |
| REST API | 85% |
| CLI commands | 80% |
| WebDAV | 75% |

## CI Testing

Tests run automatically on:

- Pull requests
- Push to main branch
- Scheduled daily runs

## Debugging Tests

### Verbose Output

```bash
go test -v ./...
```

### Race Detection

```bash
go test -race ./...
```

### Test Timeout

```bash
go test -timeout 5m ./...
```

### Skip Certain Tests

```bash
go test -run "^TestNotIntegration" ./...
```

## Test Data Cleanup

Tests automatically clean up:

- Temp directories
- Test databases
- Generated files

Using `t.Cleanup()`:

```go
func TestExample(t *testing.T) {
    tmpDir := t.TempDir()
    t.Cleanup(func() {
        // Additional cleanup
    })
}
```

## Contributing Tests

When adding new features:

1. Write tests first (TDD) or alongside
2. Ensure all new code is covered
3. Add golden file tests for parsing
4. Run `mage Coverage` to check coverage
5. Update documentation if needed
