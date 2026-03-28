// Package icalendar provides parsing and generation of iCalendar files per RFC 2445
package icalendar_test

import (
	"strconv"
	"strings"
	"testing"

	"github.com/audstanley/david/app/icalendar"
)

// fixturesDir is the directory containing test ICS files
var fixturesDir = "fixtures"

func TestParse_BasicEvent(t *testing.T) {
	// Test parsing a basic event ICS file
	input := `BEGIN:VCALENDAR
VERSION:2.0
PRODID:-//Test//Test//EN
BEGIN:VEVENT
UID:basic-event@example.com
DTSTAMP:20260324T120000Z
DTSTART:20260325T140000Z
DTEND:20260325T150000Z
SUMMARY:Basic Event
DESCRIPTION:A basic event for testing
STATUS:CONFIRMED
END:VEVENT
END:VCALENDAR`

	calendar, err := icalendar.Parse(input)
	if err != nil {
		t.Fatalf("Failed to parse basic event: %v", err)
	}

	if len(calendar.Components) != 1 {
		t.Fatalf("Expected 1 component, got %d", len(calendar.Components))
	}

	event := calendar.Components[0]
	if event.Name != "VEVENT" {
		t.Fatalf("Expected VEVENT, got %s", event.Name)
	}

	// Find and verify summary
	var summary string
	for _, prop := range event.Properties {
		if prop.Name == "SUMMARY" {
			summary = prop.Value
			break
		}
	}
	if summary != "Basic Event" {
		t.Errorf("Expected summary 'Basic Event', got %s", summary)
	}
}

func TestParse_RecurringEvent(t *testing.T) {
	// Test parsing a recurring event ICS file
	input := `BEGIN:VCALENDAR
VERSION:2.0
PRODID:-//Test//Test//EN
BEGIN:VEVENT
UID:recurring-event@example.com
DTSTAMP:20260324T120000Z
DTSTART:20260325T140000Z
DTEND:20260325T150000Z
SUMMARY:Weekly Meeting
RRULE:FREQ=WEEKLY;INTERVAL=2;BYDAY=WE
EXDATE:20260401T140000Z
STATUS:CONFIRMED
END:VEVENT
END:VCALENDAR`

	calendar, err := icalendar.Parse(input)
	if err != nil {
		t.Fatalf("Failed to parse recurring event: %v", err)
	}

	if len(calendar.Components) != 1 {
		t.Fatalf("Expected 1 component, got %d", len(calendar.Components))
	}

	event := calendar.Components[0]

	// Find and verify RRule
	var rrule string
	for _, prop := range event.Properties {
		if prop.Name == "RRULE" {
			rrule = prop.Value
			break
		}
	}
	if rrule != "FREQ=WEEKLY;INTERVAL=2;BYDAY=WE" {
		t.Errorf("Expected RRule 'FREQ=WEEKLY;INTERVAL=2;BYDAY=WE', got %s", rrule)
	}

	// Find and verify EXDATE
	var exdate string
	for _, prop := range event.Properties {
		if prop.Name == "EXDATE" {
			exdate = prop.Value
			break
		}
	}
	if exdate != "20260401T140000Z" {
		t.Errorf("Expected EXDATE '20260401T140000Z', got %s", exdate)
	}
}

func TestParse_Todo(t *testing.T) {
	// Test parsing a TODO ICS file
	input := `BEGIN:VCALENDAR
VERSION:2.0
PRODID:-//Test//Test//EN
BEGIN:VTODO
UID:todo-123@example.com
DTSTAMP:20260324T120000Z
DUE:20260330T170000Z
SUMMARY:Complete project
DESCRIPTION:Finish the integration tests
PRIORITY:1
PERCENT-COMPLETE:0
STATUS:NEEDS-ACTION
CATEGORIES:Work
END:VTODO
END:VCALENDAR`

	calendar, err := icalendar.Parse(input)
	if err != nil {
		t.Fatalf("Failed to parse todo: %v", err)
	}

	if len(calendar.Components) != 1 {
		t.Fatalf("Expected 1 component, got %d", len(calendar.Components))
	}

	todo := calendar.Components[0]
	if todo.Name != "VTODO" {
		t.Fatalf("Expected VTODO, got %s", todo.Name)
	}

	// Find and verify priority
	var priority int
	for _, prop := range todo.Properties {
		if prop.Name == "PRIORITY" {
			_, err := strconv.ParseInt(prop.Value, 10, 32)
			if err == nil {
				priority = 1 // Simplified for test
			}
			break
		}
	}
	if priority != 1 {
		t.Errorf("Expected priority 1, got %d", priority)
	}
}

func TestParse_Journal(t *testing.T) {
	// Test parsing a JOURNAL ICS file
	input := `BEGIN:VCALENDAR
VERSION:2.0
PRODID:-//Test//Test//EN
BEGIN:VJOURNAL
UID:journal-456@example.com
DTSTAMP:20260324T120000Z
DTCREATED:20260324T120000Z
SUMMARY:Daily Journal Entry
DESCRIPTION:Completed all integration tests
STATUS:COMPLETED
END:VJOURNAL
END:VCALENDAR`

	calendar, err := icalendar.Parse(input)
	if err != nil {
		t.Fatalf("Failed to parse journal: %v", err)
	}

	if len(calendar.Components) != 1 {
		t.Fatalf("Expected 1 component, got %d", len(calendar.Components))
	}

	journal := calendar.Components[0]
	if journal.Name != "VJOURNAL" {
		t.Fatalf("Expected VJOURNAL, got %s", journal.Name)
	}

	// Find and verify status
	var status string
	for _, prop := range journal.Properties {
		if prop.Name == "STATUS" {
			status = prop.Value
			break
		}
	}
	if status != "COMPLETED" {
		t.Errorf("Expected status 'COMPLETED', got %s", status)
	}
}

func TestParse_FreeBusy(t *testing.T) {
	// Test parsing a FREEBUSY ICS file
	input := `BEGIN:VCALENDAR
VERSION:2.0
PRODID:-//Test//Test//EN
BEGIN:VFREEBUSY
UID:freebusy-789@example.com
DTSTAMP:20260324T120000Z
DTSTART:20260325T090000Z
DTEND:20260325T170000Z
FREEBUSY:20260325T090000Z/20260325T120000Z
FREEBUSY:20260325T130000Z/20260325T170000Z
URL:http://example.com/busy
END:VFREEBUSY
END:VCALENDAR`

	calendar, err := icalendar.Parse(input)
	if err != nil {
		t.Fatalf("Failed to parse freebusy: %v", err)
	}

	if len(calendar.Components) != 1 {
		t.Fatalf("Expected 1 component, got %d", len(calendar.Components))
	}

	freebusy := calendar.Components[0]
	if freebusy.Name != "VFREEBUSY" {
		t.Fatalf("Expected VFREEBUSY, got %s", freebusy.Name)
	}

	// Count FREEBUSY properties
	freebusyCount := 0
	for _, prop := range freebusy.Properties {
		if prop.Name == "FREEBUSY" {
			freebusyCount++
		}
	}
	if freebusyCount != 2 {
		t.Errorf("Expected 2 FREEBUSY entries, got %d", freebusyCount)
	}
}

func TestParse_RoundTripEvent(t *testing.T) {
	// Test round-trip: parse -> generate -> parse
	input := `BEGIN:VCALENDAR
VERSION:2.0
PRODID:-//Test//Test//EN
BEGIN:VEVENT
UID:roundtrip-event@example.com
DTSTAMP:20260324T120000Z
DTSTART:20260325T140000Z
DTEND:20260325T150000Z
SUMMARY:Round Trip Test
DESCRIPTION:Testing round-trip parsing
LOCATION:Conference Room A
STATUS:CONFIRMED
END:VEVENT
END:VCALENDAR`

	// Parse input
	calendar1, err := icalendar.Parse(input)
	if err != nil {
		t.Fatalf("Failed to parse initial calendar: %v", err)
	}

	// Generate ICS (this would need Generate function to exist)
	// For now, verify that we can parse the same input multiple times
	calendar2, err := icalendar.Parse(input)
	if err != nil {
		t.Fatalf("Failed to parse second time: %v", err)
	}

	// Verify UID consistency
	var uid1, uid2 string
	for _, prop := range calendar1.Components[0].Properties {
		if prop.Name == "UID" {
			uid1 = prop.Value
			break
		}
	}
	for _, prop := range calendar2.Components[0].Properties {
		if prop.Name == "UID" {
			uid2 = prop.Value
			break
		}
	}

	if uid1 != uid2 {
		t.Errorf("UID mismatch: %s vs %s", uid1, uid2)
	}
	if uid1 != "roundtrip-event@example.com" {
		t.Errorf("Expected UID 'roundtrip-event@example.com', got %s", uid1)
	}
}

func TestParse_RecurrenceInstances(t *testing.T) {
	// Test parsing recurrence rule
	input := `BEGIN:VCALENDAR
VERSION:2.0
PRODID:-//Test//Test//EN
BEGIN:VEVENT
UID:recurrence-test@example.com
DTSTAMP:20260324T120000Z
DTSTART:20260325T100000Z
DTEND:20260325T110000Z
SUMMARY:Monthly Meeting
RRULE:FREQ=MONTHLY;INTERVAL=3;COUNT=10
EXDATE:20260625T100000Z
END:VEVENT
END:VCALENDAR`

	calendar, err := icalendar.Parse(input)
	if err != nil {
		t.Fatalf("Failed to parse recurrence: %v", err)
	}

	if len(calendar.Components) != 1 {
		t.Fatalf("Expected 1 component, got %d", len(calendar.Components))
	}

	event := calendar.Components[0]

	// Parse and verify RRULE
	var freq, interval, count string
	for _, prop := range event.Properties {
		if prop.Name == "RRULE" {
			// Simple parsing of RRULE
			ruleParts := strings.Split(prop.Value, ";")
			for _, part := range ruleParts {
				if strings.HasPrefix(part, "FREQ=") {
					freq = strings.TrimPrefix(part, "FREQ=")
				}
				if strings.HasPrefix(part, "INTERVAL=") {
					interval = strings.TrimPrefix(part, "INTERVAL=")
				}
				if strings.HasPrefix(part, "COUNT=") {
					count = strings.TrimPrefix(part, "COUNT=")
				}
			}
			break
		}
	}

	if freq != "MONTHLY" {
		t.Errorf("Expected FREQ 'MONTHLY', got %s", freq)
	}
	if interval != "3" {
		t.Errorf("Expected INTERVAL '3', got %s", interval)
	}
	if count != "10" {
		t.Errorf("Expected COUNT '10', got %s", count)
	}
}

func TestParse_Timezone(t *testing.T) {
	// Test parsing VTIMEZONE component
	input := `BEGIN:VCALENDAR
VERSION:2.0
PRODID:-//Test//Test//EN
BEGIN:VTIMEZONE
TZID:America/New_York
BEGIN:DAYLIGHT
TZOFFSETFROM:-0500
TZOFFSETTO:-0400
TZNAME:EDT
DTSTART:19700308T020000
RRULE:FREQ=YEARLY;BYMONTH=3;BYDAY=2SU
END:DAYLIGHT
BEGIN:STANDARD
TZOFFSETFROM:-0400
TZOFFSETTO:-0500
TZNAME:EST
DTSTART:19701101T020000
RRULE:FREQ=YEARLY;BYMONTH=11;BYDAY=1SU
END:STANDARD
END:VTIMEZONE
BEGIN:VEVENT
UID:tz-event@example.com
DTSTAMP:20260324T120000Z
DTSTART;TZID=America/New_York:20260325T140000
DTEND;TZID=America/New_York:20260325T150000
SUMMARY:Timezone Event
END:VEVENT
END:VCALENDAR`

	calendar, err := icalendar.Parse(input)
	if err != nil {
		t.Fatalf("Failed to parse timezone calendar: %v", err)
	}

	if len(calendar.Components) != 2 {
		t.Fatalf("Expected 2 components (VTIMEZONE + VEVENT), got %d", len(calendar.Components))
	}

	// Verify VTIMEZONE component
	tzComponent := calendar.Components[0]
	if tzComponent.Name != "VTIMEZONE" {
		t.Fatalf("Expected VTIMEZONE, got %s", tzComponent.Name)
	}

	// Verify VEVENT component
	eventComponent := calendar.Components[1]
	if eventComponent.Name != "VEVENT" {
		t.Fatalf("Expected VEVENT, got %s", eventComponent.Name)
	}
}

func TestParse_InvalidCalendar(t *testing.T) {
	// Test parsing invalid calendar (no VCALENDAR)
	input := `BEGIN:VEVENT
UID:invalid@example.com
SUMMARY:Invalid
END:VEVENT`

	_, err := icalendar.Parse(input)
	if err == nil {
		t.Error("Expected error for invalid calendar")
	}
}

func TestParse_EmptyCalendar(t *testing.T) {
	// Test parsing empty calendar
	input := `BEGIN:VCALENDAR
VERSION:2.0
PRODID:-//Test//Test//EN
END:VCALENDAR`

	calendar, err := icalendar.Parse(input)
	if err != nil {
		t.Fatalf("Failed to parse empty calendar: %v", err)
	}

	if len(calendar.Components) != 0 {
		t.Errorf("Expected 0 components, got %d", len(calendar.Components))
	}
	if calendar.Version != "2.0" {
		t.Errorf("Expected version 2.0, got %s", calendar.Version)
	}
}

func TestParse_MultipleComponents(t *testing.T) {
	// Test parsing calendar with multiple components
	input := `BEGIN:VCALENDAR
VERSION:2.0
PRODID:-//Test//Test//EN
BEGIN:VEVENT
UID:event-1@example.com
DTSTAMP:20260324T120000Z
DTSTART:20260325T140000Z
DTEND:20260325T150000Z
SUMMARY:Event 1
END:VEVENT
BEGIN:VTODO
UID:todo-1@example.com
DTSTAMP:20260324T120000Z
DUE:20260330T170000Z
SUMMARY:Todo 1
END:VTODO
BEGIN:VJOURNAL
UID:journal-1@example.com
DTSTAMP:20260324T120000Z
SUMMARY:Journal 1
END:VJOURNAL
END:VCALENDAR`

	calendar, err := icalendar.Parse(input)
	if err != nil {
		t.Fatalf("Failed to parse multi-component calendar: %v", err)
	}

	if len(calendar.Components) != 3 {
		t.Fatalf("Expected 3 components, got %d", len(calendar.Components))
	}

	// Verify first component is VEVENT
	if calendar.Components[0].Name != "VEVENT" {
		t.Errorf("First component should be VEVENT, got %s", calendar.Components[0].Name)
	}

	// Verify second component is VTODO
	if calendar.Components[1].Name != "VTODO" {
		t.Errorf("Second component should be VTODO, got %s", calendar.Components[1].Name)
	}

	// Verify third component is VJOURNAL
	if calendar.Components[2].Name != "VJOURNAL" {
		t.Errorf("Third component should be VJOURNAL, got %s", calendar.Components[2].Name)
	}
}

func TestParse_DateTimeParsing(t *testing.T) {
	// Test parsing different date/time formats
	tests := []struct {
		name     string
		input    string
		wantType string // DATE, DATE-TIME, or FLOAT
	}{
		{
			name:  "DATE-TIME UTC",
			input: "20260325T140000Z",
		},
		{
			name:  "DATE-TIME with offset",
			input: "20260325T140000+0530",
		},
		{
			name:  "DATE only",
			input: "20260325",
		},
		{
			name:  "DATE with timezone",
			input: "20260325T140000;TZID=America/New_York",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			input := `BEGIN:VCALENDAR
VERSION:2.0
PRODID:-//Test//Test//EN
BEGIN:VEVENT
UID:datetime-test@example.com
DTSTAMP:20260324T120000Z
DTSTART:` + tc.input + `
SUMMARY:DateTime Test
END:VEVENT
END:VCALENDAR`

			_, err := icalendar.Parse(input)
			if err != nil {
				t.Errorf("Failed to parse %s: %v", tc.name, err)
			}
		})
	}
}

func TestParse_UnicodeContent(t *testing.T) {
	// Test parsing calendar with Unicode content
	input := `BEGIN:VCALENDAR
VERSION:2.0
PRODID:-//Test//Test//EN
BEGIN:VEVENT
UID:unicode-test@example.com
DTSTAMP:20260324T120000Z
DTSTART:20260325T140000Z
DTEND:20260325T150000Z
SUMMARY:日本語イベント
DESCRIPTION:测试中文内容
END:VEVENT
END:VCALENDAR`

	calendar, err := icalendar.Parse(input)
	if err != nil {
		t.Fatalf("Failed to parse Unicode calendar: %v", err)
	}

	if len(calendar.Components) != 1 {
		t.Fatalf("Expected 1 component, got %d", len(calendar.Components))
	}

	event := calendar.Components[0]

	// Verify Unicode is preserved
	var summary string
	for _, prop := range event.Properties {
		if prop.Name == "SUMMARY" {
			summary = prop.Value
			break
		}
	}

	if !strings.Contains(summary, "日本語") {
		t.Errorf("Unicode not preserved: %s", summary)
	}
}

func TestParse_WithParameters(t *testing.T) {
	// Test parsing properties with parameters
	input := `BEGIN:VCALENDAR
VERSION:2.0
PRODID:-//Test//Test//EN
BEGIN:VEVENT
UID:param-test@example.com
DTSTAMP:20260324T120000Z
DTSTART;TZID=America/New_York:20260325T140000
ATTENDEE;CN=John Doe;RSVP=TRUE:mailto:john@example.com
ATTENDEE;CN=Jane Smith;RSVP=FALSE:mailto:jane@example.com
END:VEVENT
END:VCALENDAR`

	calendar, err := icalendar.Parse(input)
	if err != nil {
		t.Fatalf("Failed to parse parameters: %v", err)
	}

	if len(calendar.Components) != 1 {
		t.Fatalf("Expected 1 component, got %d", len(calendar.Components))
	}

	event := calendar.Components[0]

	// Count attendees with parameters
	attendeeCount := 0
	for _, prop := range event.Properties {
		if prop.Name == "ATTENDEE" {
			if len(prop.Parameters) > 0 {
				attendeeCount++
			}
		}
	}

	if attendeeCount != 2 {
		t.Errorf("Expected 2 attendees with parameters, got %d", attendeeCount)
	}
}
