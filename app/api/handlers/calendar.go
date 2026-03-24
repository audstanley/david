package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/audstanley/david/app/api/errors"
	"github.com/audstanley/david/app/api/models"
)

// CalendarHandler handles calendar endpoints
type CalendarHandler struct{}

// NewCalendarHandler creates a new calendar handler
func NewCalendarHandler() *CalendarHandler {
	return &CalendarHandler{}
}

// ListCalendarsHandler handles listing calendars
func (h *CalendarHandler) ListCalendarsHandler(w http.ResponseWriter, r *http.Request) {
	calendars := []map[string]interface{}{
		{
			"uid":          "cal-1",
			"display_name": "Work Calendar",
			"description":  "Work events",
			"color":        "#3788d8",
			"is_public":    false,
			"timezone":     "America/New_York",
			"created":      time.Now().Format(time.RFC3339),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(calendars)
}

// GetCalendarHandler handles getting a calendar by UID
func (h *CalendarHandler) GetCalendarHandler(w http.ResponseWriter, r *http.Request) {
	vars := map[string]string{"uid": "cal-1"}
	uid := vars["uid"] // TODO: Extract from mux.Vars(r)

	cal := map[string]interface{}{
		"uid":          uid,
		"display_name": "Work Calendar",
		"description":  "Work events",
		"color":        "#3788d8",
		"is_public":    false,
		"timezone":     "America/New_York",
		"created":      time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cal)
}

// CreateCalendarHandler handles creating a calendar
func (h *CalendarHandler) CreateCalendarHandler(w http.ResponseWriter, r *http.Request) {
	var req models.CreateCalendarRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.ValidationError("body", "invalid JSON").Write(w, http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		errors.ValidationError("request", err.Error()).Write(w, http.StatusBadRequest)
		return
	}

	cal := map[string]interface{}{
		"uid":          fmt.Sprintf("cal-%d", time.Now().Unix()),
		"display_name": req.DisplayName,
		"description":  req.Description,
		"color":        req.Color,
		"is_public":    req.IsPublic,
		"timezone":     req.Timezone,
		"created":      time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(cal)
}

// UpdateCalendarHandler handles updating a calendar
func (h *CalendarHandler) UpdateCalendarHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"uid":     "cal-1",
		"updated": true,
	})
}

// DeleteCalendarHandler handles deleting a calendar
func (h *CalendarHandler) DeleteCalendarHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
}

// ListSharesHandler handles listing calendar shares
func (h *CalendarHandler) ListSharesHandler(w http.ResponseWriter, r *http.Request) {
	shares := []map[string]interface{}{}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(shares)
}

// GrantShareHandler handles granting calendar access
func (h *CalendarHandler) GrantShareHandler(w http.ResponseWriter, r *http.Request) {
	var req models.ShareRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.ValidationError("body", "invalid JSON").Write(w, http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		errors.ValidationError("request", err.Error()).Write(w, http.StatusBadRequest)
		return
	}

	share := map[string]interface{}{
		"id":       time.Now().Unix(),
		"calendar": "cal-1",
		"user_id":  req.UserID,
		"role":     req.Role,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(share)
}

// RevokeShareHandler handles revoking calendar access
func (h *CalendarHandler) RevokeShareHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
}

// ExportCalendarHandler handles calendar export to ICS
func (h *CalendarHandler) ExportCalendarHandler(w http.ResponseWriter, r *http.Request) {
	ics := `BEGIN:VCALENDAR
VERSION:2.0
PRODID:-//david//test//EN
BEGIN:VEVENT
UID:test@example.com
DTSTAMP:20260324T120000Z
DTSTART:20260325T140000Z
DTEND:20260325T150000Z
SUMMARY:Test Event
END:VEVENT
END:VCALENDAR`

	w.Header().Set("Content-Type", "text/calendar; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=calendar.ics")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(ics))
}

// GetCalendarStatsHandler handles getting calendar statistics
func (h *CalendarHandler) GetCalendarStatsHandler(w http.ResponseWriter, r *http.Request) {
	stats := map[string]interface{}{
		"calendar_uid":  "cal-1",
		"event_count":   10,
		"todo_count":    5,
		"storage_bytes": 1024,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stats)
}
