// Package icalendar provides parsing and generation of iCalendar files per RFC 2445
package icalendar

// Errors defined in this package
import "errors"

var (
	// ErrInvalidCalendar indicates an invalid VCALENDAR structure
	ErrInvalidCalendar = errors.New("invalid VCALENDAR structure")

	// ErrMissingRequiredProperty indicates a required property is missing
	ErrMissingRequiredProperty = errors.New("missing required property")

	// ErrInvalidValue indicates an invalid property value format
	ErrInvalidValue = errors.New("invalid property value")

	// ErrInvalidUID indicates an invalid or missing UID
	ErrInvalidUID = errors.New("invalid or missing UID")

	// ErrInvalidDTStart indicates an invalid DTSTART
	ErrInvalidDTStart = errors.New("invalid DTSTART")

	// ErrInvalidRRULE indicates an invalid recurrence rule
	ErrInvalidRRULE = errors.New("invalid RRULE")

	// ErrInvalidTimezone indicates an invalid timezone definition
	ErrInvalidTimezone = errors.New("invalid timezone definition")

	// ErrMalformedICS indicates a malformed ICS file
	ErrMalformedICS = errors.New("malformed ICS file")
)
