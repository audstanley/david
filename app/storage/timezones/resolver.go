package timezones

import (
	"context"
	"fmt"
	"time"
)

// Resolver handles timezone resolution with priority chain
type Resolver struct {
	customStore *Store
	ianaLoader  *IANADownloader
	defaultTZ   string
}

// NewResolver creates a new timezone resolver
func NewResolver(customStore *Store, ianaLoader *IANADownloader, defaultTZ string) *Resolver {
	return &Resolver{
		customStore: customStore,
		ianaLoader:  ianaLoader,
		defaultTZ:   defaultTZ,
	}
}

// Resolve returns timezone rules for the given TZID
func (r *Resolver) Resolve(ctx context.Context, tzID string) (*TZRule, error) {
	// Priority 1: Check custom VTIMEZONE
	if rule, err := r.getCustomVTIMEZONE(ctx, tzID); err == nil {
		return rule, nil
	}

	// Priority 2: Check IANA tzdata
	if rule, err := r.getIANATimezone(tzID); err == nil {
		return rule, nil
	}

	// Priority 3: Fall back to UTC
	return &TZRule{
		Name:         "UTC",
		Offset:       0,
		Abbreviation: "UTC",
		IsDST:        false,
	}, nil
}

// getCustomVTIMEZONE retrieves custom timezone definition
func (r *Resolver) getCustomVTIMEZONE(ctx context.Context, tzID string) (*TZRule, error) {
	// Try to get custom VTIMEZONE from database
	// This would query the VTIMEZONE store
	// For now, return error to indicate no custom timezone
	return nil, fmt.Errorf("custom timezone not found")
}

// getIANATimezone retrieves IANA timezone data
func (r *Resolver) getIANATimezone(tzID string) (*TZRule, error) {
	// Get IANA entry
	entry, err := r.ianaLoader.GetEntry(tzID)
	if err != nil {
		return nil, fmt.Errorf("IANA timezone not found: %w", err)
	}

	// Convert to TZRule
	offset := int(entry.GmtOffset * 3600) // Convert hours to seconds

	return &TZRule{
		Name:         entry.Name,
		StartTime:    time.Now(),
		EndTime:      time.Now().Add(time.Hour * 24),
		Offset:       offset,
		Abbreviation: entry.Abbrev,
		IsDST:        entry.IsDST,
	}, nil
}

// ConvertLocalToUTC converts local time to UTC using timezone rules
func (r *Resolver) ConvertLocalToUTC(t time.Time, tzID string) (time.Time, error) {
	rule, err := r.Resolve(context.Background(), tzID)
	if err != nil {
		return time.Time{}, err
	}

	// Subtract offset to get UTC
	utc := t.Add(-time.Duration(rule.Offset) * time.Second)
	return utc.UTC(), nil
}

// ConvertUTCToLocal converts UTC to local time using timezone rules
func (r *Resolver) ConvertUTCToLocal(utc time.Time, tzID string) (time.Time, error) {
	rule, err := r.Resolve(context.Background(), tzID)
	if err != nil {
		return time.Time{}, err
	}

	// Add offset to get local time
	local := utc.Add(time.Duration(rule.Offset) * time.Second)
	return local, nil
}

// FormatInTimezone formats a time in the given timezone
func (r *Resolver) FormatInTimezone(t time.Time, tzID string, layout string) (string, error) {
	local, err := r.ConvertUTCToLocal(t, tzID)
	if err != nil {
		return "", err
	}
	return local.Format(layout), nil
}

// GetOffset returns the offset for a timezone at a specific time
func (r *Resolver) GetOffset(t time.Time, tzID string) (int, bool, error) {
	rule, err := r.Resolve(context.Background(), tzID)
	if err != nil {
		return 0, false, err
	}

	return rule.Offset, rule.IsDST, nil
}

// ListAvailable returns list of available timezone names
func (r *Resolver) ListAvailable() ([]string, error) {
	return r.ianaLoader.GetList()
}
