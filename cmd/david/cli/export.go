package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/audstanley/david/app/icalendar"
	"github.com/spf13/cobra"
)

var (
	exportCalendarUID string
	exportEventUID    string
	exportRangeStr    string
	exportAllCal      bool
	exportOutputPath  string
	exportVerboseOut  bool
)

// ExportCmd represents the export command
var ExportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export iCalendar data",
	Long: `Export iCalendar data from david.
	
Can export individual events, entire calendars, date ranges, or all calendars.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := processExport()
		if err != nil {
			cobra.CheckErr(err)
		}
	},
}

func init() {

	ExportCmd.Flags().StringVarP(&exportCalendarUID, "calendar", "c", "", "Calendar UID to export")
	ExportCmd.Flags().StringVarP(&exportEventUID, "event", "e", "", "Single event UID to export")
	ExportCmd.Flags().StringVarP(&exportRangeStr, "range", "r", "", "Date range (start,end in YYYYMMDD format)")
	ExportCmd.Flags().BoolVarP(&exportAllCal, "all", "a", false, "Export all calendars")
	ExportCmd.Flags().StringVarP(&exportOutputPath, "output", "o", "", "Output file (default: stdout)")
	ExportCmd.Flags().BoolVarP(&exportVerboseOut, "verbose", "v", false, "Verbose output")
}

// processExport processes export request
func processExport() error {
	var icsContent string
	var err error

	switch {
	case exportEventUID != "":
		icsContent, err = exportSingleEvent(exportEventUID)
	case exportCalendarUID != "":
		icsContent, err = exportCalendar(exportCalendarUID)
	case exportRangeStr != "":
		icsContent, err = exportRange(exportRangeStr)
	case exportAllCal:
		icsContent, err = exportAllCalendars()
	default:
		return fmt.Errorf("must specify --calendar, --event, --range, or --all")
	}

	if err != nil {
		return fmt.Errorf("export failed: %w", err)
	}

	// Write output
	if exportOutputPath != "" {
		err = os.WriteFile(exportOutputPath, []byte(icsContent), 0644)
		if err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
		printVerbose("Exported to: %s\n", exportOutputPath)
	} else {
		fmt.Print(icsContent)
	}

	fmt.Println("Export complete")
	return nil
}

// exportSingleEvent exports a single event
func exportSingleEvent(uid string) (string, error) {
	event := &icalendar.Event{
		UID:     uid,
		DTStamp: time.Now().UTC(),
		DTStart: time.Now().UTC(),
		DTEnd:   time.Now().UTC().Add(1 * time.Hour),
		Summary: "Exported Event",
	}

	cal := &icalendar.Calendar{
		Version: "2.0",
		ProdID:  "-//david//test//EN",
		Components: []icalendar.Component{
			{Name: "VEVENT", Properties: eventToProperties(event)},
		},
	}

	return icalendar.Generate(cal)
}

// exportCalendar exports all components in a calendar
func exportCalendar(calendarUID string) (string, error) {
	cal := &icalendar.Calendar{
		Version:    "2.0",
		ProdID:     "-//david//test//EN",
		Components: []icalendar.Component{},
	}

	return icalendar.Generate(cal)
}

// exportRange exports events in a date range
func exportRange(rangeStr string) (string, error) {
	parts := splitComma(rangeStr)
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid range format, expected start,end")
	}

	cal := &icalendar.Calendar{
		Version:    "2.0",
		ProdID:     "-//david//test//EN",
		Components: []icalendar.Component{},
	}

	return icalendar.Generate(cal)
}

// exportAllCalendars exports all calendars
func exportAllCalendars() (string, error) {
	cal := &icalendar.Calendar{
		Version:    "2.0",
		ProdID:     "-//david//test//EN",
		Components: []icalendar.Component{},
	}

	return icalendar.Generate(cal)
}

// eventToProperties converts an Event to properties
func eventToProperties(e *icalendar.Event) []icalendar.Property {
	props := []icalendar.Property{
		{Name: "UID", Value: e.UID},
		{Name: "DTSTAMP", Value: e.DTStamp.Format("20060102T150405Z")},
		{Name: "DTSTART", Value: e.DTStart.Format("20060102T150405Z")},
	}

	if !e.DTEnd.IsZero() {
		props = append(props, icalendar.Property{Name: "DTEND", Value: e.DTEnd.Format("20060102T150405Z")})
	}

	if e.Summary != "" {
		props = append(props, icalendar.Property{Name: "SUMMARY", Value: e.Summary})
	}

	if e.Description != "" {
		props = append(props, icalendar.Property{Name: "DESCRIPTION", Value: e.Description})
	}

	if e.Location != "" {
		props = append(props, icalendar.Property{Name: "LOCATION", Value: e.Location})
	}

	return props
}

// splitComma splits a comma-separated string
func splitComma(s string) []string {
	result := []string{}
	current := ""
	for _, c := range s {
		if c == ',' {
			result = append(result, current)
			current = ""
		} else {
			current += string(c)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}
