package timezones

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/audstanley/david/app/storage"
)

const (
	keyPrefixIANA      = "iana:"
	keyPrefixIANAEntry = "entry:"
	keyPrefixIANACache = "cache:"
)

// IANADownloader handles IANA tzdata downloads
type IANADownloader struct {
	store          *Store
	apiKey         string
	cacheDuration  time.Duration
	lastUpdateTime time.Time
}

// NewIANADownloader creates a new IANA downloader
func NewIANADownloader(db *storage.Storage) *IANADownloader {
	return &IANADownloader{
		store:         New(db),
		cacheDuration: time.Hour * 24, // Check once per day
	}
}

// DownloadOptions configures download behavior
type DownloadOptions struct {
	// UseFreeAPI enables free API (option A)
	UseFreeAPI bool
	// APIKey for free API service (timezonedb.com)
	APIKey string
	// UseBundled enables bundled tzdata (option B)
	UseBundled bool
	// UseManual enables manual download (option C)
	UseManual bool
}

// Download attempts to get IANA tzdata using configured options
func (d *IANADownloader) Download(ctx context.Context, opts DownloadOptions) error {
	// Check if we have valid cache
	if !d.isCacheExpired() {
		return d.loadFromCache()
	}

	// Try options in priority order
	if opts.UseFreeAPI && opts.APIKey != "" {
		if err := d.downloadFromAPI(ctx, opts.APIKey); err == nil {
			d.saveToCache()
			d.lastUpdateTime = time.Now()
			return nil
		}
	}

	if opts.UseBundled {
		if err := d.downloadFromBundled(ctx); err == nil {
			d.saveToCache()
			d.lastUpdateTime = time.Now()
			return nil
		}
	}

	if opts.UseManual {
		return ErrManualDownloadRequired
	}

	// Default: try free API without key
	if !opts.UseFreeAPI || opts.APIKey == "" {
		if err := d.downloadFromAPI(ctx, ""); err == nil {
			d.saveToCache()
			d.lastUpdateTime = time.Now()
			return nil
		}
	}

	return fmt.Errorf("all download methods failed")
}

// isCacheExpired checks if cached data is expired
func (d *IANADownloader) isCacheExpired() bool {
	if d.lastUpdateTime.IsZero() {
		return true
	}
	return time.Since(d.lastUpdateTime) > d.cacheDuration
}

// loadFromCache loads timezone data from cache
func (d *IANADownloader) loadFromCache() error {
	// Check if cache exists in database
	cacheKey := []byte(keyPrefixIANACache + "last_update")
	data, err := d.store.db.Get(cacheKey)
	if err != nil {
		return fmt.Errorf("no cache found")
	}

	var cache struct {
		LastUpdate string `json:"last_update"`
	}
	if err := json.Unmarshal(data, &cache); err != nil {
		return fmt.Errorf("failed to parse cache: %w", err)
	}

	if cache.LastUpdate != "" {
		d.lastUpdateTime, _ = time.Parse(time.RFC3339, cache.LastUpdate)
	}

	return nil
}

// saveToCache saves timezone data to cache
func (d *IANADownloader) saveToCache() error {
	cache := struct {
		LastUpdate string `json:"last_update"`
	}{
		LastUpdate: time.Now().UTC().Format(time.RFC3339),
	}

	data, err := json.Marshal(cache)
	if err != nil {
		return fmt.Errorf("failed to marshal cache: %w", err)
	}

	return d.store.db.Put([]byte(keyPrefixIANACache+"last_update"), data)
}

// downloadFromAPI downloads from free API (option A)
func (d *IANADownloader) downloadFromAPI(ctx context.Context, apiKey string) error {
	// Free API without key - limited data
	// For full data, users need API key from timezonedb.com
	url := "https://timezonedb.com/api/timezones"

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	// Parse and store timezone data
	var zones []TZDataEntry
	if err := json.Unmarshal(body, &zones); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	for _, zone := range zones {
		// TODO: Store in database
		// d.store.AddIANAEntry(ctx, zone.Name, zone)
		_ = zone
	}

	return nil
}

// downloadFromBundled loads bundled tzdata (option B)
func (d *IANADownloader) downloadFromBundled(ctx context.Context) error {
	// This would load tzdata bundled with the binary
	// For now, return error to indicate this needs implementation
	// In production, you'd embed tzdata files and parse them
	return ErrBundledNotAvailable
}

// downloadFromManual requires manual download (option C)
func (d *IANADownloader) downloadFromManual(ctx context.Context) error {
	return ErrManualDownloadRequired
}

// GetList returns list of available IANA timezone names
func (d *IANADownloader) GetList() ([]string, error) {
	var names []string

	err := d.store.db.Iterate([]byte(keyPrefixIANAEntry), func(key, value []byte) error {
		names = append(names, string(value))
		return nil
	})

	return names, err
}

// GetEntry retrieves timezone data by name
func (d *IANADownloader) GetEntry(name string) (*TZDataEntry, error) {
	// Return sample entry for now - full implementation would query database
	return &TZDataEntry{
		Name:        name,
		CountryCode: "US",
		TzName:      name,
	}, nil
}
