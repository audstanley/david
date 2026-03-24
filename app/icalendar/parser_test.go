// Package icalendar provides parsing and generation of iCalendar files per RFC 2445
package icalendar

import (
	"strings"
	"testing"
	"time"
)

func TestParse_BasicCalendar(t *testing.T) {
	input := `BEGIN:VCALENDAR
VERSION:2.0
PRODID:-//Test//Test//EN
BEGIN:VEVENT
UID:test-event@example.com
DTSTAMP:20260324T120000Z
DTSTART:20260325T140000Z
DTEND:20260325T150000Z
SUMMARY:Test Event
STATUS:CONFIRMED
END:VEVENT
END:VCALENDAR`

	calendar, err := Parse(input)
	if err != nil {
		t.Fatalf("Failed to parse calendar: %v", err)
	}

	if calendar.Version != "2.0" {
		t.Errorf("Expected version 2.0, got %s", calendar.Version)
	}
	if calendar.ProdID != "-//Test//Test//EN" {
		t.Errorf("Expected prodid -//Test//Test//EN, got %s", calendar.ProdID)
	}
	if len(calendar.Components) != 1 {
		t.Errorf("Expected 1 component, got %d", len(calendar.Components))
	}

	event := calendar.Components[0]
	if event.Name != "VEVENT" {
		t.Errorf("Expected component name VEVENT, got %s", event.Name)
	}

	// Find UID property
	var uid string
	for _, prop := range event.Properties {
		if prop.Name == "UID" {
			uid = prop.Value
			break
		}
	}
	if uid != "test-event@example.com" {
		t.Errorf("Expected UID test-event@example.com, got %s", uid)
	}
}

func TestParse_DateTime(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantYear int
		wantMon  time.Month
		wantDay  int
		wantHour int
		wantMin  int
		wantSec  int
	}{
		{
			name:     "UTC time",
			input:    "20260325T140000Z",
			wantYear: 2026,
			wantMon:  3,
			wantDay:  25,
			wantHour: 14,
			wantMin:  0,
			wantSec:  0,
		},
		{
			name:     "Positive offset",
			input:    "20260325T140000+0530",
			wantYear: 2026,
			wantMon:  3,
			wantDay:  25,
			wantHour: 14,
			wantMin:  0,
			wantSec:  0,
		},
		{
			name:     "Negative offset",
			input:    "20260325T140000-0800",
			wantYear: 2026,
			wantMon:  3,
			wantDay:  25,
			wantHour: 14,
			wantMin:  0,
			wantSec:  0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tt, err := ParseDateTime(tc.input)
			if err != nil {
				t.Fatalf("Failed to parse datetime: %v", err)
			}
			year, month, day := tt.Date()

			if year != tc.wantYear {
				t.Errorf("Expected year %d, got %d", tc.wantYear, year)
			}
			if month != tc.wantMon {
				t.Errorf("Expected month %d, got %d", tc.wantMon, month)
			}
			if day != tc.wantDay {
				t.Errorf("Expected day %d, got %d", tc.wantDay, day)
			}

			_, offset := tt.Zone()
			expectedOffset := 0
			if tc.name == "Positive offset" {
				expectedOffset = 5*3600 + 30*60
			} else if tc.name == "Negative offset" {
				expectedOffset = -(8 * 3600)
			}
			if offset != expectedOffset && tc.name != "UTC time" {
				t.Errorf("Expected offset %d, got %d", expectedOffset, offset)
			}
		})
	}
}

func TestParse_Date(t *testing.T) {
	input := "20260325"
	tt, err := ParseDate(input)
	if err != nil {
		t.Fatalf("Failed to parse date: %v", err)
	}

	if tt.Year() != 2026 {
		t.Errorf("Expected year 2026, got %d", tt.Year())
	}
	if tt.Month() != 3 {
		t.Errorf("Expected month 3, got %d", tt.Month())
	}
	if tt.Day() != 25 {
		t.Errorf("Expected day 25, got %d", tt.Day())
	}
}

func TestParse_Duration(t *testing.T) {
	// Test days (basic functionality)
	dur, err := ParseDuration("P2D")
	if err != nil {
		t.Fatalf("Failed to parse duration: %v", err)
	}
	expected := 48 * time.Hour
	if dur != expected {
		t.Errorf("Expected %v, got %v", expected, dur)
	}
}

func TestGenerate_BasicEvent(t *testing.T) {
	calendar := &Calendar{
		Version: "2.0",
		ProdID:  "-//Test//Test//EN",
		Components: []Component{
			{
				Name: "VEVENT",
				Properties: []Property{
					{Name: "UID", Value: "test-event@example.com"},
					{Name: "DTSTAMP", Value: "20260324T120000Z"},
					{Name: "DTSTART", Value: "20260325T140000Z"},
					{Name: "DTEND", Value: "20260325T150000Z"},
					{Name: "SUMMARY", Value: "Test Event"},
				},
			},
		},
	}

	output, err := Generate(calendar)
	if err != nil {
		t.Fatalf("Failed to generate: %v", err)
	}

	if !strings.HasPrefix(output, "BEGIN:VCALENDAR") {
		t.Errorf("Output should start with BEGIN:VCALENDAR")
	}
	if !strings.HasSuffix(output, "END:VCALENDAR\r\n") {
		t.Errorf("Output should end with END:VCALENDAR")
	}
}
