// Package icalendar provides parsing and generation of iCalendar files per RFC 2445
package icalendar

import (
	"strconv"
	"strings"
	"time"
)

// ParseEvent parses a VEVENT component from properties
func ParseEvent(props []Property, alarms []Component) (*Event, error) {
	event := &Event{}

	for _, prop := range props {
		switch prop.Name {
		case "UID":
			event.UID = prop.Value
		case "DTSTAMP":
			if t, err := ParseDateTime(prop.Value); err == nil {
				event.DTStamp = t
			}
		case "DTSTART":
			if t, err := ParseDateTime(prop.Value); err == nil {
				event.DTStart = t
			} else if t, err := ParseDate(prop.Value); err == nil {
				// For all-day events
				event.DTStart = t
			}
		case "DTEND":
			if t, err := ParseDateTime(prop.Value); err == nil {
				event.DTEnd = t
			} else if t, err := ParseDate(prop.Value); err == nil {
				event.DTEnd = t
			}
		case "DURATION":
			if dur, err := ParseDuration(prop.Value); err == nil {
				event.Duration = &dur
			}
		case "SUMMARY":
			event.Summary = prop.Value
		case "DESCRIPTION":
			event.Description = prop.Value
		case "LOCATION":
			event.Location = prop.Value
		case "ORGANIZER":
			event.Organizer = prop.Value
		case "ATTENDEE":
			event.Attendees = append(event.Attendees, parseAttendee(prop))
		case "STATUS":
			event.Status = prop.Value
		case "CLASS":
			event.Class = prop.Value
		case "PRIORITY":
			if prio, err := parseInt(prop.Value); err == nil {
				event.Priority = prio
			}
		case "SEQUENCE":
			if seq, err := parseInt(prop.Value); err == nil {
				event.Sequence = seq
			}
		case "CATEGORIES":
			event.Categories = splitComma(prop.Value)
		case "RESOURCES":
			event.Resources = splitComma(prop.Value)
		case "RRULE":
			if rrule, err := parseRRULE(prop.Value); err == nil {
				event.RRule = rrule
			}
		case "EXDATE":
			if dates, err := parseDateList(prop.Value); err == nil {
				event.ExDates = append(event.ExDates, dates...)
			}
		case "RDATE":
			if dates, err := parseDateList(prop.Value); err == nil {
				event.RDates = append(event.RDates, dates...)
			}
		case "RELATED-TO":
			event.RelatedTo = prop.Value
		case "GEO":
			if geo, err := parseGeo(prop.Value); err == nil {
				event.Geo = geo
			}
		case "URL":
			event.URL = prop.Value
		case "CREATED":
			if t, err := ParseDateTime(prop.Value); err == nil {
				event.Created = t
			}
		case "LAST-MODIFIED":
			if t, err := ParseDateTime(prop.Value); err == nil {
				event.LastModified = t
			}
		}
	}

	// Parse alarms
	for _, alarmComp := range alarms {
		if alarmComp.Name == "VALARM" {
			if alarm, err := parseAlarm(alarmComp.Properties); err == nil {
				event.Alarms = append(event.Alarms, alarm)
			}
		}
	}

	// Validate required properties
	if event.UID == "" {
		return nil, ErrInvalidUID
	}

	return event, nil
}

// parseAttendee parses an ATTENDEE property
func parseAttendee(prop Property) Attendee {
	attendee := Attendee{
		Email:    prop.Value,
		RSVP:     false,
		Role:     "REQ-PARTICIPANT",
		Cutype:   "INDIVIDUAL",
		PartStat: "NEEDS-ACTION",
	}

	// Parse parameters
	if rsvp, ok := prop.Parameters["RSVP"]; ok && len(rsvp) > 0 {
		attendee.RSVP = rsvp[0] == "TRUE"
	}
	if role, ok := prop.Parameters["ROLE"]; ok && len(role) > 0 {
		attendee.Role = role[0]
	}
	if partstat, ok := prop.Parameters["PARTSTAT"]; ok && len(partstat) > 0 {
		attendee.PartStat = partstat[0]
	}
	if cutype, ok := prop.Parameters["CUTYPE"]; ok && len(cutype) > 0 {
		attendee.Cutype = cutype[0]
	}
	if cn, ok := prop.Parameters["CN"]; ok && len(cn) > 0 {
		attendee.Name = cn[0]
	}

	return attendee
}

// parseAlarm parses a VALARM component
func parseAlarm(props []Property) (Alarm, error) {
	alarm := Alarm{
		Action: "DISPLAY",
	}

	for _, prop := range props {
		switch prop.Name {
		case "ACTION":
			alarm.Action = prop.Value
		case "TRIGGER":
			alarm.Trigger = prop.Value
			if prop.Value[0] == '-' {
				// Negative trigger (before)
				if dur, err := ParseDuration(prop.Value[1:]); err == nil {
					alarm.Duration = -dur
				}
			} else {
				if dur, err := ParseDuration(prop.Value); err == nil {
					alarm.Duration = dur
				}
			}
		case "DURATION":
			if dur, err := ParseDuration(prop.Value); err == nil {
				alarm.Duration = dur
			}
		case "DESCRIPTION":
			alarm.Description = prop.Value
		case "ATTACH":
			alarm.Attach = prop.Value
		case "SUMMARY":
			alarm.Summary = prop.Value
		}
	}

	return alarm, nil
}

// generateEvent generates a VEVENT component
func generateEvent(e *Event) string {
	var buf string

	buf += "BEGIN:VEVENT\r\n"

	buf += "UID:" + e.UID + "\r\n"
	buf += "DTSTAMP:" + formatDateTime(time.Now()) + "\r\n"
	buf += "DTSTART:" + formatDateTime(e.DTStart) + "\r\n"

	if !e.DTEnd.IsZero() {
		buf += "DTEND:" + formatDateTime(e.DTEnd) + "\r\n"
	}
	if e.Duration != nil {
		buf += "DURATION:" + formatDuration(*e.Duration) + "\r\n"
	}

	if e.Summary != "" {
		buf += "SUMMARY:" + escapeValue(e.Summary) + "\r\n"
	}
	if e.Description != "" {
		buf += "DESCRIPTION:" + escapeValue(e.Description) + "\r\n"
	}
	if e.Location != "" {
		buf += "LOCATION:" + escapeValue(e.Location) + "\r\n"
	}
	if e.Organizer != "" {
		buf += "ORGANIZER:" + e.Organizer + "\r\n"
	}
	if e.Status != "" {
		buf += "STATUS:" + e.Status + "\r\n"
	}
	if e.Class != "" {
		buf += "CLASS:" + e.Class + "\r\n"
	}
	if e.Priority > 0 {
		buf += "PRIORITY:" + itoa(e.Priority) + "\r\n"
	}
	if e.Sequence > 0 {
		buf += "SEQUENCE:" + itoa(e.Sequence) + "\r\n"
	}
	if len(e.Categories) > 0 {
		buf += "CATEGORIES:" + stringsJoin(e.Categories, ",") + "\r\n"
	}
	if len(e.Resources) > 0 {
		buf += "RESOURCES:" + stringsJoin(e.Resources, ",") + "\r\n"
	}
	if e.RRule != nil {
		buf += "RRULE:" + generateRRule(e.RRule) + "\r\n"
	}
	if len(e.ExDates) > 0 {
		buf += "EXDATE:" + formatDateList(e.ExDates) + "\r\n"
	}
	if len(e.RDates) > 0 {
		buf += "RDATE:" + formatDateList(e.RDates) + "\r\n"
	}
	if e.URL != "" {
		buf += "URL:" + e.URL + "\r\n"
	}
	if e.Created.IsZero() {
		buf += "CREATED:" + formatDateTime(time.Now()) + "\r\n"
	} else {
		buf += "CREATED:" + formatDateTime(e.Created) + "\r\n"
	}
	if !e.LastModified.IsZero() {
		buf += "LAST-MODIFIED:" + formatDateTime(e.LastModified) + "\r\n"
	}

	// Alarms
	for _, alarm := range e.Alarms {
		buf += generateAlarm(&alarm)
	}

	buf += "END:VEVENT\r\n"

	return buf
}

// generateAlarm generates a VALARM component
func generateAlarm(a *Alarm) string {
	var buf string

	buf += "BEGIN:VALARM\r\n"
	buf += "ACTION:" + a.Action + "\r\n"

	if a.Duration < 0 {
		buf += "TRIGGER:-" + formatDuration(-a.Duration) + "\r\n"
	} else {
		buf += "TRIGGER:" + formatDuration(a.Duration) + "\r\n"
	}

	if a.Description != "" {
		buf += "DESCRIPTION:" + escapeValue(a.Description) + "\r\n"
	}
	if a.Summary != "" {
		buf += "SUMMARY:" + escapeValue(a.Summary) + "\r\n"
	}

	buf += "END:VALARM\r\n"

	return buf
}

// Helper functions

func itoa(n int) string {
	return strconv.Itoa(n)
}

func splitComma(s string) []string {
	return strings.Split(s, ",")
}
