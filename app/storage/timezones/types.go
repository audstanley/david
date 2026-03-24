package timezones

import (
	"errors"
	"time"
)

// Errors
var (
	ErrManualDownloadRequired = errors.New("manual IANA tzdata download required - see documentation")
	ErrBundledNotAvailable    = errors.New("bundled tzdata not available - enable bundled download option")
	ErrNotFound               = errors.New("timezone not found")
)

// TZDataEntry represents IANA timezone data entry
type TZDataEntry struct {
	Name          string  `json:"id"`
	CountryCode   string  `json:"countryCode"`
	CountryName   string  `json:"countryName"`
	GmtOffset     float64 `json:"gmtOffset"`
	GmtOffsetName string  `json:"gmtOffsetName"`
	Abbrev        string  `json:"abbreviation"`
	TzName        string  `json:"tzName"`
	IsDST         bool    `json:"isDst"`
}

// TZRule represents a timezone rule (offset transition)
type TZRule struct {
	Name         string
	StartTime    time.Time
	EndTime      time.Time
	Offset       int // seconds from UTC
	Abbreviation string
	IsDST        bool
}

// TimezoneConfig holds timezone configuration
type TimezoneConfig struct {
	// UseFreeAPI enables free API (option A)
	UseFreeAPI bool
	// APIKey for free API service
	APIKey string
	// UseBundled enables bundled tzdata (option B)
	UseBundled bool
	// UseManual enables manual download (option C)
	UseManual bool
	// CacheDuration how long to cache downloaded data
	CacheDuration time.Duration
}

// DefaultTimezoneConfig returns default timezone settings
func DefaultTimezoneConfig() TimezoneConfig {
	return TimezoneConfig{
		UseFreeAPI:    true,
		UseBundled:    true,
		UseManual:     true,
		CacheDuration: time.Hour * 24,
	}
}
