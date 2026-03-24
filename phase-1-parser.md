# Phase 1: Core iCalendar Parser & Generator

## Overview
This phase implements the core iCalendar (RFC 2445) parser and generator. This is the foundation for all iCalendar functionality - without it, we cannot parse imported files or generate exported files.

## Goals
- Implement full RFC 2445 compliance for iCalendar parsing
- Implement bidirectional serialization (Go structs ↔ ICS files)
- Support all value types defined in RFC 2445
- Support all property parameters
- Implement RRULE parser with validation
- Achieve 95%+ unit test coverage
- Implement fuzz testing for parser robustness

## Tasks

### 1. Core Parser Structure
- [ ] Create `app/icalendar/parser.go` - Main parser entry point
- [ ] Implement content line parsing (name, parameters, value)
- [ ] Implement line folding/unfolding (RFC 2445 Section 4.1)
- [ ] Implement character set handling (UTF-8 default)
- [ ] Implement escaped character handling (backslash escaping)
- [ ] Create `app/icalendar/errors.go` - Custom error types
- [ ] Implement error handling with wrap/unwrap

### 2. Value Type Parsers
- [ ] Create `app/icalendar/value_types.go` - Value type definitions
- [ ] Implement TEXT value parser
- [ ] Implement URI value parser
- [ ] Implement BOOLEAN value parser (TRUE/FALSE)
- [ ] Implement INTEGER value parser
- [ ] Implement FLOAT value parser
- [ ] Implement DATE value parser (YYYYMMDD)
- [ ] Implement TIME value parser (HHMMSS)
- [ ] Implement DATE-TIME value parser (YYYYMMDDTHHMMSS[Z|ZZZ])
- [ ] Implement DURATION value parser (PTnHnMnS, PnD, etc.)
- [ ] Implement PERIOD value parser (start/end or start/duration)
- [ ] Implement CAL-ADDRESS value parser (mailto URI)
- [ ] Implement BINARY value parser (BASE64 encoding)
- [ ] Implement UID value parser (unique identifier format)

### 3. Property Parameter Parsers
- [ ] Create `app/icalendar/parameter.go` - Parameter definitions
- [ ] Implement ALTREP parameter parser
- [ ] Implement CN parameter parser (Common Name)
- [ ] Implement CUTYPE parameter parser (Calendar User Type)
- [ ] Implement DELEGATED-FROM parameter parser
- [ ] Implement DELEGATED-TO parameter parser
- [ ] Implement DIR parameter parser (Directory Reference)
- [ ] Implement ENCODING parameter parser (8BIT, BASE64)
- [ ] Implement FMTTYPE parameter parser (Format Type)
- [ ] Implement FBTYPE parameter parser (Free/Busy Type)
- [ ] Implement LANGUAGE parameter parser
- [ ] Implement MEMBER parameter parser (Group Membership)
- [ ] Implement PARTSTAT parameter parser (Participation Status)
- [ ] Implement RANGE parameter parser (Recurrence Range)
- [ ] Implement RELTYPE parameter parser (Relationship Type)
- [ ] Implement ROLE parameter parser (Participation Role)
- [ ] Implement RSVP parameter parser (RSVP Expectation)
- [ ] Implement SENT-BY parameter parser
- [ ] Implement TZID parameter parser (Time Zone ID)
- [ ] Implement VALUE parameter parser (Value Type hint)
- [ ] Implement X-param parser (Vendor extensions)

### 4. Component Types
- [ ] Create `app/icalendar/types.go` - Main type definitions
- [ ] Define `Calendar` struct (VCALENDAR container)
- [ ] Define `Event` struct (VEVENT component)
- [ ] Define `Todo` struct (VTODO component)
- [ ] Define `Journal` struct (VJOURNAL component)
- [ ] Define `FreeBusy` struct (VFREEBUSY component)
- [ ] Define `TimeZone` struct (VTIMEZONE component)
- [ ] Define `TimeZoneDaylight` struct (STANDARD subcomponent)
- [ ] Define `TimeZoneStandard` struct (DAYLIGHT subcomponent)
- [ ] Define `Alarm` struct (VALARM component)
- [ ] Define `Attendee` struct
- [ ] Define `RRule` struct (Recurrence Rule)
- [ ] Define `Property` struct (generic property)

### 5. Component-Specific Files
- [ ] Create `app/icalendar/component/vcalendar.go`
  - [ ] Parse VERSION property (required, must be 2.0)
  - [ ] Parse PRODID property (required)
  - [ ] Parse CALSESCALE property (optional)
  - [ ] Parse METHOD property (optional)
- [ ] Create `app/icalendar/component/vevent.go`
  - [ ] Parse all VEVENT properties
  - [ ] Implement UID validation (required)
  - [ ] Implement DTSTAMP validation (required)
  - [ ] Handle DTSTART (required)
  - [ ] Handle DTEND or DURATION (mutually exclusive, one required if all-day event)
  - [ ] Parse SUMMARY (required for some contexts)
  - [ ] Parse ORGANIZER (required in iTIP)
  - [ ] Parse ATTENDEE (multiple allowed)
  - [ ] Parse RRULE (recurrence rule)
  - [ ] Parse EXDATE, RDATE (exception/date ranges)
  - [ ] Parse VALARM subcomponents
- [ ] Create `app/icalendar/component/vtodo.go`
  - [ ] Parse all VTODO properties
  - [ ] Implement UID validation (required)
  - [ ] Parse DTSTAMP (required)
  - [ ] Parse DTSTART (optional)
  - [ ] Parse DUE or DURATION (mutually exclusive)
  - [ ] Parse COMPLETED (optional)
  - [ ] Parse STATUS (optional)
  - [ ] Parse PERCENT-COMPLETE
  - [ ] Parse PRIORITY
- [ ] Create `app/icalendar/component/vjournal.go`
  - [ ] Parse all VJOURNAL properties
  - [ ] Implement UID validation (required)
  - [ ] Parse DTSTAMP (required)
  - [ ] Parse DTSTART (optional)
  - [ ] Parse STATUS (optional)
- [ ] Create `app/icalendar/component/vfreebusy.go`
  - [ ] Parse all VFREEBUSY properties
  - [ ] Implement UID validation (required)
  - [ ] Parse DTSTAMP (required)
  - [ ] Parse DTSTART, DTEND (required)
  - [ ] Parse FREE-BUSY (multiple allowed)
- [ ] Create `app/icalendar/component/vtimezone.go`
  - [ ] Parse VTIMEZONE properties
  - [ ] Parse TZID (required)
  - [ ] Parse LAST-MODIFIED (optional)
  - [ ] Parse STANDARD subcomponent
  - [ ] Parse DAYLIGHT subcomponent
  - [ ] Parse TZOFFSETFROM, TZOFFSETTO
  - [ ] Parse TZNAME
  - [ ] Parse RDATE for timezone rules
- [ ] Create `app/icalendar/component/valarm.go`
  - [ ] Parse VALARM properties
  - [ ] Parse TRIGGER (required)
  - [ ] Parse ACTION (required)
  - [ ] Parse DESCRIPTION
  - [ ] Parse DURATION
  - [ ] Parse ATTACH

### 6. RRULE Parser
- [ ] Create `app/icalendar/rrule.go` - RRULE-specific parsing
- [ ] Implement FREQ parser (SECONDLY, MINUTELY, HOURLY, DAILY, WEEKLY, MONTHLY, YEARLY)
- [ ] Implement INTERVAL parser
- [ ] Implement BYSECOND parser
- [ ] Implement BYMINUTE parser
- [ ] Implement BYHOUR parser
- [ ] Implement BYDAY parser (MO, TU, WE, TH, FR, SA, SU, +nMO, etc.)
- [ ] Implement BYMONTHDAY parser
- [ ] Implement BYYEARDAY parser
- [ ] Implement BYWEEKNO parser
- [ ] Implement BYMONTH parser
- [ ] Implement BYSETPOS parser
- [ ] Implement COUNT parser
- [ ] Implement UNTIL parser
- [ ] Implement EXDATE parser (exception dates)
- [ ] Implement validation of RRULE constraints
- [ ] Implement error handling for invalid RRULE

### 7. Generator (Serialization)
- [ ] Create `app/icalendar/generator.go` - Main generator
- [ ] Implement content line generation
- [ ] Implement line folding (lines > 75 octets)
- [ ] Implement property serialization
- [ ] Implement parameter serialization
- [ ] Implement value type serialization
- [ ] Implement component serialization
- [ ] Implement calendar serialization (VCALENDAR wrapper)
- [ ] Implement all value type generators:
  - [ ] TEXT generator (escaping)
  - [ ] URI generator
  - [ ] BOOLEAN generator
  - [ ] INTEGER generator
  - [ ] FLOAT generator
  - [ ] DATE generator
  - [ ] TIME generator
  - [ ] DATE-TIME generator (with timezone handling)
  - [ ] DURATION generator
  - [ ] PERIOD generator
  - [ ] BINARY generator (BASE64)
- [ ] Implement RRULE generator

### 8. Validation
- [ ] Create `app/icalendar/validation.go` - Validation logic
- [ ] Implement required property validation per component
- [ ] Implement conditional property validation
- [ ] Implement property value format validation
- [ ] Implement RRULE validation
- [ ] Implement timezone validation
- [ ] Implement duplicate property detection
- [ ] Implement component nesting validation (VALARM inside components)
- [ ] Create validation errors with context

### 9. Test Fixtures
- [ ] Create `test/fixtures/ics/event_basic.ics` - Basic event
- [ ] Create `test/fixtures/ics/event_recurring.ics` - Recurring event
- [ ] Create `test/fixtures/ics/event_timezone.ics` - Event with timezone
- [ ] Create `test/fixtures/ics/event_alarm.ics` - Event with alarms
- [ ] Create `test/fixtures/ics/todo_basic.ics` - Basic todo
- [ ] Create `test/fixtures/ics/todo_completed.ics` - Completed todo
- [ ] Create `test/fixtures/ics/journal_basic.ics` - Basic journal
- [ ] Create `test/fixtures/ics/freebusy_basic.ics` - Basic free/busy
- [ ] Create `test/fixtures/ics/timezone_america_new_york.ics` - Timezone definition
- [ ] Create `test/fixtures/ics/malformed.ics` - Invalid ICS for error testing
- [ ] Create `test/fixtures/ics/unicode.ics` - Unicode content testing
- [ ] Create `test/fixtures/ics/long_description.ics` - Line folding testing

### 10. Unit Tests
- [ ] Create `app/icalendar/parser_test.go` - Table-driven tests
  - [ ] Test content line parsing
  - [ ] Test line folding/unfolding
  - [ ] Test parameter parsing
  - [ ] Test value type parsing
- [ ] Create `app/icalendar/value_types_test.go`
  - [ ] Test each value type parser
  - [ ] Test each value type generator
- [ ] Create `app/icalendar/generator_test.go`
  - [ ] Test round-trip (parse → generate → parse)
  - [ ] Test line folding
  - [ ] Test escaping
- [ ] Create `app/icalendar/rrule_test.go`
  - [ ] Test RRULE parsing
  - [ ] Test RRULE validation
  - [ ] Test RRULE generation
- [ ] Create `app/icalendar/validation_test.go`
  - [ ] Test required property validation
  - [ ] Test conditional validation
  - [ ] Test duplicate detection
- [ ] Create `app/icalendar/component/vevent_test.go`
- [ ] Create `app/icalendar/component/vtodo_test.go`
- [ ] Create `app/icalendar/component/vjournal_test.go`
- [ ] Create `app/icalendar/component/vfreebusy_test.go`
- [ ] Create `app/icalendar/component/vtimezone_test.go`
- [ ] Create `app/icalendar/component/valarm_test.go`

### 11. Golden File Tests
- [ ] Set up golden file testing infrastructure
- [ ] Create golden files for all test fixtures
- [ ] Implement comparison logic (ignore whitespace differences)
- [ ] Add golden file tests to test suite
- [ ] Document golden file update process

### 12. Fuzz Testing
- [ ] Create `app/icalendar/fuzz_test.go`
- [ ] Implement fuzzer for content line parsing
- [ ] Implement fuzzer for property parsing
- [ ] Implement fuzzer for value type parsing
- [ ] Implement fuzzer for component parsing
- [ ] Configure fuzz test timeouts
- [ ] Add fuzz tests to CI

## Data Models

### Core Types
```go
// Calendar represents a VCALENDAR component
type Calendar struct {
    Version    string
    ProdID     string
    CalScale   string  // GREGORIAN or X-WSCALE
    Method     string  // REQUEST, REPLY, etc.
    Properties []Property
    Components []Component
}

// Component represents any calendar component
type Component struct {
    Name       string  // VEVENT, VTODO, etc.
    Properties []Property
    Components []Component  // Nested components (VALARM)
}

// Property represents a calendar property
type Property struct {
    Name       string
    Parameters map[string][]string
    Value      string
    ValueType  string  // TEXT, URI, DATE-TIME, etc.
}

// RRule represents a recurrence rule
type RRule struct {
    Freq         string
    Interval     int
    BySecond     []int
    ByMinute     []int
    ByHour       []int
    ByDay        []string
    ByMonthDay   []int
    ByYearDay    []int
    ByWeekNo     []int
    ByMonth      []int
    BySetPos     []int
    Count        int
    Until        time.Time
    ExcludeDates []time.Time
}
```

## Success Criteria
- [ ] All RFC 2445 value types parse correctly
- [ ] All RFC 2445 property parameters parse correctly
- [ ] All component types parse correctly
- [ ] RRULE parsing handles all valid rules
- [ ] Round-trip parsing/generation preserves data
- [ ] Line folding/unfolding works correctly
- [ ] Validation catches invalid iCalendar data
- [ ] 95%+ unit test coverage
- [ ] Fuzz tests run successfully (no crashes)
- [ ] Golden file tests pass
- [ ] All tests pass with `go test -race ./app/icalendar/...`

## Notes
- RFC 2445 is complex - refer to the spec frequently
- Be strict on parsing (fail gracefully on invalid data)
- Be lenient on generation (follow RFC but handle edge cases)
- Document any non-standard extensions you support
- Keep parser and generator separate for clarity

## Next Phase
Move to **Phase 2: LevelDB Storage Layer** when all items above are complete and coverage is verified.
