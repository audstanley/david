package models

import (
	"errors"
	"net/url"
	"strings"
	"time"
)

// AuthRequest represents a login request
type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Validate validates the auth request
func (r *AuthRequest) Validate() error {
	if r.Username == "" {
		return errors.New("username required")
	}
	if r.Password == "" {
		return errors.New("password required")
	}
	return nil
}

// RefreshRequest represents a refresh token request
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// Validate validates the refresh request
func (r *RefreshRequest) Validate() error {
	if r.RefreshToken == "" {
		return errors.New("refresh token required")
	}
	return nil
}

// CreateCalendarRequest represents a calendar creation request
type CreateCalendarRequest struct {
	DisplayName string `json:"display_name"`
	Description string `json:"description,omitempty"`
	Color       string `json:"color,omitempty"`
	Timezone    string `json:"timezone,omitempty"`
	IsPublic    bool   `json:"is_public"`
}

// Validate validates the calendar creation request
func (r *CreateCalendarRequest) Validate() error {
	if r.DisplayName == "" {
		return errors.New("display name required")
	}
	if len(r.DisplayName) > 200 {
		return errors.New("display name too long")
	}
	if r.Timezone != "" && !isValidTimezone(r.Timezone) {
		return errors.New("invalid timezone")
	}
	return nil
}

// UpdateCalendarRequest represents a calendar update request
type UpdateCalendarRequest struct {
	DisplayName string `json:"display_name,omitempty"`
	Description string `json:"description,omitempty"`
	Color       string `json:"color,omitempty"`
	Timezone    string `json:"timezone,omitempty"`
	IsPublic    *bool  `json:"is_public,omitempty"`
}

// Validate validates the calendar update request
func (r *UpdateCalendarRequest) Validate() error {
	if r.DisplayName != "" && len(r.DisplayName) > 200 {
		return errors.New("display name too long")
	}
	if r.Timezone != "" && !isValidTimezone(r.Timezone) {
		return errors.New("invalid timezone")
	}
	return nil
}

// CreateEventRequest represents an event creation request
type CreateEventRequest struct {
	UID         string   `json:"uid,omitempty"`
	Summary     string   `json:"summary"`
	Description string   `json:"description,omitempty"`
	Location    string   `json:"location,omitempty"`
	DTStart     string   `json:"dtstart"`
	DTEnd       string   `json:"dtend,omitempty"`
	DTStartTZ   string   `json:"dtstart_tz,omitempty"`
	DTEndTZ     string   `json:"dtend_tz,omitempty"`
	AllDay      bool     `json:"all_day,omitempty"`
	Attendees   []string `json:"attendees,omitempty"`
	Status      string   `json:"status,omitempty"`
	Class       string   `json:"class,omitempty"`
	Priority    int      `json:"priority,omitempty"`
	Categories  []string `json:"categories,omitempty"`
	RRule       string   `json:"rrule,omitempty"`
	Alarms      []Alarm  `json:"alarms,omitempty"`
}

// Alarm represents an event alarm
type Alarm struct {
	Action   string `json:"action"`
	Trigger  string `json:"trigger"`
	Duration string `json:"duration,omitempty"`
	Describe string `json:"description,omitempty"`
}

// Validate validates the event creation request
func (r *CreateEventRequest) Validate() error {
	if r.Summary == "" {
		return errors.New("summary required")
	}
	if r.DTStart == "" {
		return errors.New("dtstart required")
	}
	return nil
}

// UpdateEventRequest represents an event update request
type UpdateEventRequest struct {
	Summary     *string  `json:"summary,omitempty"`
	Description *string  `json:"description,omitempty"`
	Location    *string  `json:"location,omitempty"`
	DTStart     *string  `json:"dtstart,omitempty"`
	DTEnd       *string  `json:"dtend,omitempty"`
	Status      *string  `json:"status,omitempty"`
	Class       *string  `json:"class,omitempty"`
	Priority    *int     `json:"priority,omitempty"`
	Categories  []string `json:"categories,omitempty"`
}

// Validate validates the event update request
func (r *UpdateEventRequest) Validate() error {
	return nil
}

// CreateTodoRequest represents a todo creation request
type CreateTodoRequest struct {
	UID         string   `json:"uid,omitempty"`
	Summary     string   `json:"summary"`
	Description string   `json:"description,omitempty"`
	Due         string   `json:"due,omitempty"`
	Percent     int      `json:"percent_complete,omitempty"`
	Status      string   `json:"status,omitempty"`
	Priority    int      `json:"priority,omitempty"`
	Categories  []string `json:"categories,omitempty"`
}

// Validate validates the todo creation request
func (r *CreateTodoRequest) Validate() error {
	if r.Summary == "" {
		return errors.New("summary required")
	}
	return nil
}

// UpdateTodoRequest represents a todo update request
type UpdateTodoRequest struct {
	Summary     *string `json:"summary,omitempty"`
	Description *string `json:"description,omitempty"`
	Due         *string `json:"due,omitempty"`
	Percent     *int    `json:"percent_complete,omitempty"`
	Status      *string `json:"status,omitempty"`
	Priority    *int    `json:"priority,omitempty"`
}

// Validate validates the todo update request
func (r *UpdateTodoRequest) Validate() error {
	return nil
}

// CreateJournalRequest represents a journal creation request
type CreateJournalRequest struct {
	UID         string `json:"uid,omitempty"`
	Summary     string `json:"summary"`
	Description string `json:"description,omitempty"`
	DTCreated   string `json:"dtcreated,omitempty"`
	Status      string `json:"status,omitempty"`
}

// Validate validates the journal creation request
func (r *CreateJournalRequest) Validate() error {
	if r.Summary == "" {
		return errors.New("summary required")
	}
	return nil
}

// UpdateJournalRequest represents a journal update request
type UpdateJournalRequest struct {
	Summary     *string `json:"summary,omitempty"`
	Description *string `json:"description,omitempty"`
	Status      *string `json:"status,omitempty"`
}

// Validate validates the journal update request
func (r *UpdateJournalRequest) Validate() error {
	return nil
}

// ShareRequest represents a share creation request
type ShareRequest struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
}

// Validate validates the share request
func (r *ShareRequest) Validate() error {
	if r.UserID == "" {
		return errors.New("user ID required")
	}
	if r.Role != "read" && r.Role != "write" && r.Role != "admin" && r.Role != "everyone" {
		return errors.New("invalid role")
	}
	return nil
}

// QueryParams represents query parameters for list endpoints
type QueryParams struct {
	Page     int       `url:"page"`
	PerPage  int       `url:"per_page"`
	From     time.Time `url:"from"`
	To       time.Time `url:"to"`
	Status   string    `url:"status"`
	Calendar string    `url:"calendar"`
}

// Validate validates the query params
func (q *QueryParams) Validate() error {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.PerPage < 1 || q.PerPage > 100 {
		q.PerPage = 50
	}
	return nil
}

// GetQueryParams extracts query params from URL
func GetQueryParams(u *url.URL) *QueryParams {
	q := &QueryParams{
		Page:    1,
		PerPage: 50,
	}

	values := u.Query()
	if p := values.Get("page"); p != "" {
		q.Page = parseInt(p)
	}
	if pp := values.Get("per_page"); pp != "" {
		q.PerPage = parseInt(pp)
	}
	if f := values.Get("from"); f != "" {
		if t, err := time.Parse("2006-01-02T15:04:05Z", f); err == nil {
			q.From = t
		}
	}
	if t := values.Get("to"); t != "" {
		if t, err := time.Parse("2006-01-02T15:04:05Z", t); err == nil {
			q.To = t
		}
	}
	q.Status = values.Get("status")
	q.Calendar = values.Get("calendar")

	return q
}

// CreateAPIKeyRequest represents an API key creation request
type CreateAPIKeyRequest struct {
	Name        string   `json:"name"`
	Permissions []string `json:"permissions,omitempty"`
	ExpiryDays  int      `json:"expiry_days,omitempty"`
}

// Validate validates the API key creation request
func (r *CreateAPIKeyRequest) Validate() error {
	if r.Name == "" {
		return errors.New("name required")
	}
	if len(r.Name) > 100 {
		return errors.New("name too long")
	}
	if r.ExpiryDays < 1 || r.ExpiryDays > 365 {
		return errors.New("expiry must be between 1 and 365 days")
	}
	return nil
}

// Helper functions

func parseInt(s string) int {
	var n int
	if _, err := time.ParseDuration(s); err == nil {
		return 0 // invalid for int
	}
	// Simple parsing, assume it's a number
	for _, c := range s {
		if c >= '0' && c <= '9' {
			n = n*10 + int(c-'0')
		}
	}
	return n
}

func isValidTimezone(tz string) bool {
	if tz == "" {
		return true
	}
	// Common timezones
	commonTimezones := map[string]bool{
		"UTC":                 true,
		"America/New_York":    true,
		"America/Los_Angeles": true,
		"America/Chicago":     true,
		"Europe/London":       true,
		"Europe/Paris":        true,
		"Europe/Berlin":       true,
		"Asia/Tokyo":          true,
		"Asia/Shanghai":       true,
		"Australia/Sydney":    true,
	}
	return commonTimezones[strings.TrimSpace(tz)] || strings.Contains(tz, "/")
}
