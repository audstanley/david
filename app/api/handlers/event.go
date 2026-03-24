package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/audstanley/david/app/api/errors"
	"github.com/audstanley/david/app/api/models"
)

// EventHandler handles event endpoints
type EventHandler struct{}

// NewEventHandler creates a new event handler
func NewEventHandler() *EventHandler {
	return &EventHandler{}
}

// ListEventsHandler handles listing events
func (h *EventHandler) ListEventsHandler(w http.ResponseWriter, r *http.Request) {
	events := []map[string]interface{}{
		{
			"uid":         "event-1",
			"summary":     "Team Meeting",
			"description": "Weekly team sync",
			"location":    "Conference Room A",
			"dtstart":     "2026-03-25T14:00:00Z",
			"dtend":       "2026-03-25T15:00:00Z",
			"status":      "CONFIRMED",
			"class":       "PUBLIC",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(events)
}

// GetEventHandler handles getting an event by UID
func (h *EventHandler) GetEventHandler(w http.ResponseWriter, r *http.Request) {
	event := map[string]interface{}{
		"uid":         "event-1",
		"summary":     "Team Meeting",
		"description": "Weekly team sync",
		"location":    "Conference Room A",
		"dtstart":     "2026-03-25T14:00:00Z",
		"dtend":       "2026-03-25T15:00:00Z",
		"status":      "CONFIRMED",
		"class":       "PUBLIC",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(event)
}

// CreateEventHandler handles creating an event
func (h *EventHandler) CreateEventHandler(w http.ResponseWriter, r *http.Request) {
	var req models.CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.ValidationError("body", "invalid JSON").Write(w, http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		errors.ValidationError("request", err.Error()).Write(w, http.StatusBadRequest)
		return
	}

	event := map[string]interface{}{
		"uid":         "event-" + time.Now().Format("20060102150405"),
		"summary":     req.Summary,
		"description": req.Description,
		"location":    req.Location,
		"dtstart":     req.DTStart,
		"dtend":       req.DTEnd,
		"status":      req.Status,
		"class":       req.Class,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(event)
}

// UpdateEventHandler handles updating an event
func (h *EventHandler) UpdateEventHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"uid":     "event-1",
		"updated": true,
	})
}

// DeleteEventHandler handles deleting an event
func (h *EventHandler) DeleteEventHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
}

// ExportEventHandler handles exporting a single event to ICS
func (h *EventHandler) ExportEventHandler(w http.ResponseWriter, r *http.Request) {
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
	w.Header().Set("Content-Disposition", "attachment; filename=event.ics")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(ics))
}
