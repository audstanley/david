// Package icalendar provides parsing and generation of iCalendar files per RFC 2445
package icalendar

import (
	"bufio"
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Calendar represents a VCALENDAR component
type Calendar struct {
	Version    string
	ProdID     string
	CalScale   string
	Method     string
	Properties []Property
	Components []Component
}

// Component represents any calendar component (VEVENT, VTODO, etc.)
type Component struct {
	Name       string
	Properties []Property
	Components []Component
}

// Property represents a calendar property
type Property struct {
	Name       string
	Parameters map[string][]string
	Value      string
	ValueType  string
}

// Parse parses an iCalendar string into a Calendar struct
func Parse(input string) (*Calendar, error) {
	lines, err := unfold(input)
	if err != nil {
		return nil, err
	}

	calendar := &Calendar{}
	currentComponent := &Component{}
	inCalendar := false

	scanner := bufio.NewScanner(bytes.NewReader([]byte(strings.Join(lines, "\n"))))
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		// Check for VCALENDAR boundaries
		if strings.HasPrefix(trimmed, "BEGIN:VCALENDAR") {
			inCalendar = true
			continue
		}
		if strings.HasPrefix(trimmed, "END:VCALENDAR") {
			if !inCalendar {
				return nil, ErrInvalidCalendar
			}
			if len(currentComponent.Properties) > 0 || len(currentComponent.Components) > 0 {
				calendar.Components = append(calendar.Components, *currentComponent)
			}
			break
		}

		// Check for component boundaries
		if strings.HasPrefix(trimmed, "BEGIN:") {
			componentName := strings.TrimPrefix(trimmed, "BEGIN:")
			if componentName == "VEVENT" || componentName == "VTODO" || componentName == "VJOURNAL" ||
				componentName == "VFREEBUSY" || componentName == "VTIMEZONE" || componentName == "VALARM" {
				if currentComponent.Name != "" {
					currentComponent.Components = append(currentComponent.Components, Component{
						Name:       currentComponent.Name,
						Properties: currentComponent.Properties,
						Components: currentComponent.Components,
					})
				}
				currentComponent = &Component{Name: componentName}
			}
			continue
		}
		if strings.HasPrefix(trimmed, "END:") {
			componentName := strings.TrimPrefix(trimmed, "END:")
			if componentName == currentComponent.Name {
				if inCalendar && (componentName == "VEVENT" || componentName == "VTODO" || componentName == "VJOURNAL" ||
					componentName == "VFREEBUSY" || componentName == "VTIMEZONE") {
					calendar.Components = append(calendar.Components, Component{
						Name:       currentComponent.Name,
						Properties: currentComponent.Properties,
						Components: currentComponent.Components,
					})
				}
				currentComponent = &Component{}
			}
			continue
		}

		// Parse property line
		if inCalendar {
			prop, err := parsePropertyLine(trimmed)
			if err != nil {
				return nil, fmt.Errorf("error parsing property line %q: %w", trimmed, err)
			}

			// Add to current component if inside one, otherwise to calendar
			if currentComponent.Name != "" {
				currentComponent.Properties = append(currentComponent.Properties, prop)
			} else {
				calendar.Properties = append(calendar.Properties, prop)
				setCalendarProperty(calendar, prop)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if !inCalendar {
		return nil, ErrInvalidCalendar
	}

	return calendar, nil
}

// unfold handles line folding per RFC 2445 Section 4.1
func unfold(input string) ([]string, error) {
	var lines []string
	var currentLine strings.Builder

	scanner := bufio.NewScanner(bytes.NewReader([]byte(input)))
	for scanner.Scan() {
		line := scanner.Text()

		// Check for folded line (continues with whitespace)
		if len(lines) > 0 && (strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t")) {
			// Remove leading whitespace and append to current line
			currentLine.WriteString(strings.TrimLeft(line, " \t"))
			continue
		}

		// Save previous line and start new one
		if currentLine.Len() > 0 {
			lines = append(lines, currentLine.String())
		}
		currentLine.Reset()
		currentLine.WriteString(line)
	}

	// Add last line
	if currentLine.Len() > 0 {
		lines = append(lines, currentLine.String())
	}

	return lines, scanner.Err()
}

// parsePropertyLine parses a single property line
func parsePropertyLine(line string) (Property, error) {
	// Find the first colon (property name separator)
	colonIdx := strings.Index(line, ":")
	if colonIdx == -1 {
		return Property{}, ErrMalformedICS
	}

	name := line[:colonIdx]
	value := line[colonIdx+1:]

	// Check for parameters (before the value colon)
	parameters := make(map[string][]string)
	if semicolonIdx := strings.Index(name, ";"); semicolonIdx != -1 {
		paramStr := name[semicolonIdx+1:]
		name = name[:semicolonIdx]

		// Parse parameters
		for _, paramPart := range splitByUnquotedComma(paramStr) {
			if eqIdx := strings.Index(paramPart, "="); eqIdx != -1 {
				paramName := strings.ToUpper(paramPart[:eqIdx])
				paramValue := paramPart[eqIdx+1:]

				// Handle quoted values
				if strings.HasPrefix(paramValue, "\"") && strings.HasSuffix(paramValue, "\"") {
					paramValue = paramValue[1 : len(paramValue)-1]
				}

				// Handle list values (comma-separated in unquoted context)
				if strings.Contains(paramValue, ",") && !strings.HasPrefix(paramValue, "\"") {
					parameters[paramName] = splitByUnquotedComma(paramValue)
				} else {
					parameters[paramName] = []string{paramValue}
				}
			}
		}
	}

	// Determine value type from VALUE parameter or infer
	valueType := inferValueType(name, value)
	if valParam, ok := parameters["VALUE"]; ok && len(valParam) > 0 {
		valueType = valParam[0]
	}

	return Property{
		Name:       strings.ToUpper(name),
		Parameters: parameters,
		Value:      value,
		ValueType:  valueType,
	}, nil
}

// splitByUnquotedComma splits a string by commas outside of quotes
func splitByUnquotedComma(s string) []string {
	var result []string
	var current strings.Builder
	inQuotes := false

	for _, ch := range s {
		if ch == '"' {
			inQuotes = !inQuotes
			current.WriteRune(ch)
		} else if ch == ',' && !inQuotes {
			result = append(result, current.String())
			current.Reset()
		} else {
			current.WriteRune(ch)
		}
	}

	if current.Len() > 0 {
		result = append(result, current.String())
	}

	return result
}

// inferValueType infers the value type from property name and value
func inferValueType(name, value string) string {
	// Common property names and their value types
	valueTypeMap := map[string]string{
		"UID":              "TEXT",
		"DTSTAMP":          "DATE-TIME",
		"DTSTART":          "DATE-TIME",
		"DTEND":            "DATE-TIME",
		"DUE":              "DATE-TIME",
		"COMPLETED":        "DATE-TIME",
		"CREATED":          "DATE-TIME",
		"LAST-MODIFIED":    "DATE-TIME",
		"ORGANIZER":        "CAL-ADDRESS",
		"ATTENDEE":         "CAL-ADDRESS",
		"URL":              "URI",
		"ATTACH":           "URI",
		"RRULE":            "RECUR",
		"DESCRIPTION":      "TEXT",
		"SUMMARY":          "TEXT",
		"LOCATION":         "TEXT",
		"STATUS":           "TEXT",
		"PRIORITY":         "INTEGER",
		"PERCENT-COMPLETE": "INTEGER",
		"CLASS":            "TEXT",
		"TRANSP":           "TEXT",
		"CATEGORIES":       "TEXT",
		"RESOURCES":        "TEXT",
		"COMMENT":          "TEXT",
		"CONTACT":          "TEXT",
		"GEO":              "FLOAT",
	}

	if vt, ok := valueTypeMap[name]; ok {
		return vt
	}
	return "TEXT"
}

// setCalendarProperty sets calendar-level properties
func setCalendarProperty(c *Calendar, prop Property) {
	switch prop.Name {
	case "VERSION":
		c.Version = prop.Value
	case "PRODID":
		c.ProdID = prop.Value
	case "CALSESCALE":
		c.CalScale = prop.Value
	case "METHOD":
		c.Method = prop.Value
	}
}

// Generate serializes a Calendar back to ICS format
func Generate(cal *Calendar) (string, error) {
	var buf strings.Builder

	buf.WriteString("BEGIN:VCALENDAR\r\n")

	// Required properties
	if cal.Version == "" {
		return "", ErrMissingRequiredProperty
	}
	if cal.ProdID == "" {
		return "", ErrMissingRequiredProperty
	}

	// Calendar properties
	buf.WriteString(fmt.Sprintf("VERSION:%s\r\n", cal.Version))
	buf.WriteString(fmt.Sprintf("PRODID:%s\r\n", cal.ProdID))
	if cal.CalScale != "" {
		buf.WriteString(fmt.Sprintf("CALSECALE:%s\r\n", cal.CalScale))
	}
	if cal.Method != "" {
		buf.WriteString(fmt.Sprintf("METHOD:%s\r\n", cal.Method))
	}

	// Components
	for _, comp := range cal.Components {
		compStr, err := generateComponent(comp)
		if err != nil {
			return "", err
		}
		buf.WriteString(compStr)
	}

	buf.WriteString("END:VCALENDAR\r\n")

	return buf.String(), nil
}

// generateComponent generates a component in ICS format
func generateComponent(comp Component) (string, error) {
	var buf strings.Builder

	buf.WriteString(fmt.Sprintf("BEGIN:%s\r\n", comp.Name))

	for _, prop := range comp.Properties {
		propStr, err := generateProperty(prop)
		if err != nil {
			return "", err
		}
		buf.WriteString(propStr)
	}

	for _, subComp := range comp.Components {
		subCompStr, err := generateComponent(subComp)
		if err != nil {
			return "", err
		}
		buf.WriteString(subCompStr)
	}

	buf.WriteString(fmt.Sprintf("END:%s\r\n", comp.Name))

	return buf.String(), nil
}

// generateProperty generates a property in ICS format
func generateProperty(prop Property) (string, error) {
	var buf strings.Builder

	// Property name
	buf.WriteString(prop.Name)

	// Parameters
	if len(prop.Parameters) > 0 {
		buf.WriteString(";")
		params := make([]string, 0, len(prop.Parameters))
		for name, values := range prop.Parameters {
			for _, val := range values {
				// Quote values containing special characters
				if strings.ContainsAny(val, ",;:") || strings.Contains(val, "\"") {
					val = "\"" + escapeQuotes(val) + "\""
				}
				params = append(params, fmt.Sprintf("%s=%s", name, val))
			}
		}
		buf.WriteString(strings.Join(params, ","))
	}

	// Value
	buf.WriteString(fmt.Sprintf(":%s\r\n", escapeValue(prop.Value)))

	return buf.String(), nil
}

// escapeQuotes escapes double quotes in a string
func escapeQuotes(s string) string {
	return strings.ReplaceAll(s, "\"", "\\\"")
}

// escapeValue escapes special characters in a value
func escapeValue(s string) string {
	// Escape backslashes first, then other special chars
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, ";", "\\;")
	s = strings.ReplaceAll(s, ",", "\\,")
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\r", "")
	return s
}

// parseInt parses an integer string
func parseInt(s string) (int, error) {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, ErrInvalidValue
	}
	return n, nil
}

// stringsJoin joins strings with a separator
func stringsJoin(slice []string, sep string) string {
	return strings.Join(slice, sep)
}

// ParseDateTime parses a DATE-TIME value
func ParseDateTime(s string) (time.Time, error) {
	// Format: YYYYMMDDTHHMMSS[Z|ZZZ]
	if len(s) < 15 {
		return time.Time{}, ErrInvalidValue
	}

	year, err := strconv.Atoi(s[0:4])
	if err != nil {
		return time.Time{}, ErrInvalidValue
	}
	month, err := strconv.Atoi(s[4:6])
	if err != nil {
		return time.Time{}, ErrInvalidValue
	}
	day, err := strconv.Atoi(s[6:8])
	if err != nil {
		return time.Time{}, ErrInvalidValue
	}
	hour, err := strconv.Atoi(s[9:11])
	if err != nil {
		return time.Time{}, ErrInvalidValue
	}
	minute, err := strconv.Atoi(s[11:13])
	if err != nil {
		return time.Time{}, ErrInvalidValue
	}
	second, err := strconv.Atoi(s[13:15])
	if err != nil {
		return time.Time{}, ErrInvalidValue
	}

	t := time.Date(year, time.Month(month), day, hour, minute, second, 0, time.UTC)

	// Check for timezone indicator
	if len(s) > 15 {
		if s[15] == 'Z' {
			// UTC
			return t, nil
		} else if s[15] == '+' || s[15] == '-' {
			// Offset timezone
			return parseOffsetTime(t, s[15:])
		}
	}

	return t, nil
}

// parseOffsetTime parses a time with timezone offset
func parseOffsetTime(t time.Time, offset string) (time.Time, error) {
	if len(offset) != 5 {
		return time.Time{}, ErrInvalidValue
	}

	sign := 1
	if offset[0] == '-' {
		sign = -1
	}

	hours, err := strconv.Atoi(offset[1:3])
	if err != nil {
		return time.Time{}, ErrInvalidValue
	}
	minutes, err := strconv.Atoi(offset[3:5])
	if err != nil {
		return time.Time{}, ErrInvalidValue
	}

	offsetSeconds := sign * (hours*3600 + minutes*60)
	loc := time.FixedZone("UTC+OFFSET", offsetSeconds)

	return t.In(loc), nil
}

// ParseDate parses a DATE value (YYYYMMDD)
func ParseDate(s string) (time.Time, error) {
	if len(s) < 8 {
		return time.Time{}, ErrInvalidValue
	}

	year, err := strconv.Atoi(s[0:4])
	if err != nil {
		return time.Time{}, ErrInvalidValue
	}
	month, err := strconv.Atoi(s[4:6])
	if err != nil {
		return time.Time{}, ErrInvalidValue
	}
	day, err := strconv.Atoi(s[6:8])
	if err != nil {
		return time.Time{}, ErrInvalidValue
	}

	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC), nil
}

// ParseDuration parses a DURATION value (PTnHnMnS or PnD)
func ParseDuration(s string) (time.Duration, error) {
	if !strings.HasPrefix(s, "P") {
		return 0, ErrInvalidValue
	}

	s = strings.TrimPrefix(s, "P")

	var hours, minutes, seconds, days int

	timePart := false
	for i := 0; i < len(s); i++ {
		ch := s[i]
		if ch == 'T' {
			timePart = true
			continue
		}
		if ch == 'H' {
			numStr := s[:i]
			if numStr != "" {
				hours, _ = strconv.Atoi(numStr)
			}
			s = s[i+1:]
			i = -1
		} else if ch == 'M' {
			numStr := s[:i]
			if numStr != "" {
				minutes, _ = strconv.Atoi(numStr)
			}
			s = s[i+1:]
			i = -1
		} else if ch == 'S' {
			numStr := s[:i]
			if numStr != "" {
				seconds, _ = strconv.Atoi(numStr)
			}
		} else if ch >= '0' && ch <= '9' {
			if !timePart {
				days, _ = strconv.Atoi(string(ch))
			}
		}
	}

	return time.Duration(days)*24*time.Hour +
		time.Duration(hours)*time.Hour +
		time.Duration(minutes)*time.Minute +
		time.Duration(seconds)*time.Second, nil
}

// parseRRULE parses a recurrence rule
func parseRRULE(s string) (*RRule, error) {
	rrule := &RRule{}
	pairs := strings.Split(s, ";")

	for _, pair := range pairs {
		if eqIdx := strings.Index(pair, "="); eqIdx != -1 {
			key := strings.ToUpper(pair[:eqIdx])
			value := pair[eqIdx+1:]

			switch key {
			case "FREQ":
				rrule.Freq = value
			case "INTERVAL":
				if n, err := parseInt(value); err == nil {
					rrule.Interval = n
				}
			case "BYSECOND":
				for _, v := range strings.Split(value, ",") {
					if n, err := parseInt(v); err == nil {
						rrule.BySecond = append(rrule.BySecond, n)
					}
				}
			case "BYMINUTE":
				for _, v := range strings.Split(value, ",") {
					if n, err := parseInt(v); err == nil {
						rrule.ByMinute = append(rrule.ByMinute, n)
					}
				}
			case "BYHOUR":
				for _, v := range strings.Split(value, ",") {
					if n, err := parseInt(v); err == nil {
						rrule.ByHour = append(rrule.ByHour, n)
					}
				}
			case "BYDAY":
				rrule.ByDay = strings.Split(value, ",")
			case "BYMONTHDAY":
				for _, v := range strings.Split(value, ",") {
					if n, err := parseInt(v); err == nil {
						rrule.ByMonthDay = append(rrule.ByMonthDay, n)
					}
				}
			case "BYYEARDAY":
				for _, v := range strings.Split(value, ",") {
					if n, err := parseInt(v); err == nil {
						rrule.ByYearDay = append(rrule.ByYearDay, n)
					}
				}
			case "BYWEEKNO":
				for _, v := range strings.Split(value, ",") {
					if n, err := parseInt(v); err == nil {
						rrule.ByWeekNo = append(rrule.ByWeekNo, n)
					}
				}
			case "BYMONTH":
				for _, v := range strings.Split(value, ",") {
					if n, err := parseInt(v); err == nil {
						rrule.ByMonth = append(rrule.ByMonth, n)
					}
				}
			case "BYSETPOS":
				for _, v := range strings.Split(value, ",") {
					if n, err := parseInt(v); err == nil {
						rrule.BySetPos = append(rrule.BySetPos, n)
					}
				}
			case "COUNT":
				if n, err := parseInt(value); err == nil {
					rrule.Count = n
				}
			case "UNTIL":
				if t, err := ParseDateTime(value); err == nil {
					rrule.Until = t
				}
			}
		}
	}

	return rrule, nil
}

// parseDateList parses a comma-separated date list
func parseDateList(s string) ([]time.Time, error) {
	var dates []time.Time
	for _, dateStr := range strings.Split(s, ",") {
		// Try DATE-TIME first, then DATE
		if t, err := ParseDateTime(dateStr); err == nil {
			dates = append(dates, t)
		} else if t, err := ParseDate(dateStr); err == nil {
			dates = append(dates, t)
		}
	}
	return dates, nil
}

// parseGeo parses GEO value (lat;lon)
func parseGeo(s string) (*Geo, error) {
	parts := strings.Split(s, ";")
	if len(parts) != 2 {
		return nil, ErrInvalidValue
	}

	lat, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return nil, ErrInvalidValue
	}
	lon, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return nil, ErrInvalidValue
	}

	return &Geo{Lat: lat, Lon: lon}, nil
}

// formatDateList formats dates for output
func formatDateList(dates []time.Time) string {
	var parts []string
	for _, d := range dates {
		parts = append(parts, formatDateTime(d))
	}
	return strings.Join(parts, ",")
}

// formatDateTime formats a time for ICS output
func formatDateTime(t time.Time) string {
	return t.UTC().Format("20060102T150405Z")
}

// formatDuration formats a duration for ICS output
func formatDuration(d time.Duration) string {
	if d < 0 {
		d = -d
	}

	hours := int(d / time.Hour)
	minutes := int((d % time.Hour) / time.Minute)
	seconds := int((d % time.Minute) / time.Second)

	var buf strings.Builder
	buf.WriteString("P")

	if hours > 0 || minutes > 0 || seconds > 0 {
		if hours > 0 {
			buf.WriteString(strconv.Itoa(hours))
			buf.WriteString("H")
		}
		if minutes > 0 {
			buf.WriteString(strconv.Itoa(minutes))
			buf.WriteString("M")
		}
		if seconds > 0 {
			buf.WriteString(strconv.Itoa(seconds))
			buf.WriteString("S")
		}
	} else {
		buf.WriteString("0S")
	}

	return buf.String()
}

// generateRRule generates RRULE string
func generateRRule(rrule *RRule) string {
	var buf strings.Builder
	buf.WriteString("FREQ=" + rrule.Freq)

	if rrule.Interval > 1 {
		buf.WriteString(";INTERVAL=" + strconv.Itoa(rrule.Interval))
	}

	if len(rrule.BySecond) > 0 {
		buf.WriteString(";BYSECOND=" + stringsJoin(intsToStrings(rrule.BySecond), ","))
	}
	if len(rrule.ByMinute) > 0 {
		buf.WriteString(";BYMINUTE=" + stringsJoin(intsToStrings(rrule.ByMinute), ","))
	}
	if len(rrule.ByHour) > 0 {
		buf.WriteString(";BYHOUR=" + stringsJoin(intsToStrings(rrule.ByHour), ","))
	}
	if len(rrule.ByDay) > 0 {
		buf.WriteString(";BYDAY=" + stringsJoin(rrule.ByDay, ","))
	}
	if len(rrule.ByMonthDay) > 0 {
		buf.WriteString(";BYMONTHDAY=" + stringsJoin(intsToStrings(rrule.ByMonthDay), ","))
	}
	if len(rrule.ByYearDay) > 0 {
		buf.WriteString(";BYYEARDAY=" + stringsJoin(intsToStrings(rrule.ByYearDay), ","))
	}
	if len(rrule.ByWeekNo) > 0 {
		buf.WriteString(";BYWEEKNO=" + stringsJoin(intsToStrings(rrule.ByWeekNo), ","))
	}
	if len(rrule.ByMonth) > 0 {
		buf.WriteString(";BYMONTH=" + stringsJoin(intsToStrings(rrule.ByMonth), ","))
	}
	if len(rrule.BySetPos) > 0 {
		buf.WriteString(";BYSETPOS=" + stringsJoin(intsToStrings(rrule.BySetPos), ","))
	}
	if rrule.Count > 0 {
		buf.WriteString(";COUNT=" + strconv.Itoa(rrule.Count))
	}
	if !rrule.Until.IsZero() {
		buf.WriteString(";UNTIL=" + formatDateTime(rrule.Until))
	}

	return buf.String()
}

func intsToStrings(ints []int) []string {
	result := make([]string, len(ints))
	for i, n := range ints {
		result[i] = strconv.Itoa(n)
	}
	return result
}
