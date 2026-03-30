package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/audstanley/david/app/api/errors"
	"github.com/audstanley/david/app/api/models"
	"github.com/audstanley/david/app/storage"
	"github.com/audstanley/david/app/storage/events"
)

// EventHandler handles event endpoints
type EventHandler struct {
	db         *storage.Storage
	eventStore *events.Store
}

// NewEventHandler creates a new event handler with storage dependencies
func NewEventHandler(db *storage.Storage) *EventHandler {
	eventStore := events.New(db)
	return &EventHandler{
		db:         db,
		eventStore: eventStore,
	}
}

// ListEventsHandler handles listing events
func (h *EventHandler) ListEventsHandler(w http.ResponseWriter, r *http.Request) {
	calendarUID := r.URL.Query().Get("calendar_uid")
	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")

	var startTime, endTime time.Time
	if startStr != "" {
		startTime, _ = time.Parse(time.RFC3339, startStr)
	}
	if endStr != "" {
		endTime, _ = time.Parse(time.RFC3339, endStr)
	}

	eventList, err := h.eventStore.List(calendarUID, startTime, endTime, 100, 0)
	if err != nil {
		errors.InternalError(err).Write(w, http.StatusInternalServerError)
		return
	}

	events := make([]map[string]interface{}, len(eventList))
	for i, event := range eventList {
		events[i] = map[string]interface{}{
			"uid":         event.UID,
			"summary":     event.Summary,
			"description": event.Description,
			"location":    event.Location,
			"dtstart":     event.DTStart.Format(time.RFC3339),
			"dtend":       event.DTEnd.Format(time.RFC3339),
			"status":      event.Status,
			"class":       event.Class,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(events)
}

// GetEventHandler handles getting an event by UID
func (h *EventHandler) GetEventHandler(w http.ResponseWriter, r *http.Request) {
	uid := r.URL.Query().Get("uid")
	if uid == "" {
		errors.ValidationError("uid", "required").Write(w, http.StatusBadRequest)
		return
	}

	event, err := h.eventStore.GetByUID(uid)
	if err != nil {
		errors.NotFound("event", uid).Write(w, http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"uid":         event.UID,
		"summary":     event.Summary,
		"description": event.Description,
		"location":    event.Location,
		"dtstart":     event.DTStart.Format(time.RFC3339),
		"dtend":       event.DTEnd.Format(time.RFC3339),
		"status":      event.Status,
		"class":       event.Class,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
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

	eventUID := fmt.Sprintf("event-%d", time.Now().Unix())
	calendarUID := "cal-1"
	event := &storage.Event{
		UID:         eventUID,
		CalendarUID: calendarUID,
		DTStamp:     time.Now(),
		DTStart:     time.Now(),
		DTEnd:       time.Now(),
		Summary:     req.Summary,
		Description: req.Description,
		Location:    req.Location,
		Status:      req.Status,
		Class:       req.Class,
		Sequence:    0,
	}

	err := h.eventStore.Create(event)
	if err != nil {
		errors.InternalError(err).Write(w, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"uid":         event.UID,
		"summary":     event.Summary,
		"description": event.Description,
		"location":    event.Location,
		"dtstart":     event.DTStart.Format(time.RFC3339),
		"dtend":       event.DTEnd.Format(time.RFC3339),
		"status":      event.Status,
		"class":       event.Class,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// UpdateEventHandler handles updating an event
func (h *EventHandler) UpdateEventHandler(w http.ResponseWriter, r *http.Request) {
	uid := r.URL.Query().Get("uid")
	if uid == "" {
		errors.ValidationError("uid", "required").Write(w, http.StatusBadRequest)
		return
	}

	var req models.UpdateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.ValidationError("body", "invalid JSON").Write(w, http.StatusBadRequest)
		return
	}

	event, err := h.eventStore.GetByUID(uid)
	if err != nil {
		errors.NotFound("event", uid).Write(w, http.StatusNotFound)
		return
	}

	if req.Summary != nil {
		event.Summary = *req.Summary
	}
	if req.Description != nil {
		event.Description = *req.Description
	}
	if req.Location != nil {
		event.Location = *req.Location
	}
	if req.DTStart != nil {
		event.DTStart = time.Now()
	}
	if req.DTEnd != nil {
		event.DTEnd = time.Now()
	}
	if req.Status != nil {
		event.Status = *req.Status
	}

	event.Sequence++
	err = h.eventStore.Update(uid, event, event.Sequence)
	if err != nil {
		errors.InternalError(err).Write(w, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"uid":     event.UID,
		"updated": true,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// DeleteEventHandler handles deleting an event
func (h *EventHandler) DeleteEventHandler(w http.ResponseWriter, r *http.Request) {
	uid := r.URL.Query().Get("uid")
	if uid == "" {
		errors.ValidationError("uid", "required").Write(w, http.StatusBadRequest)
		return
	}

	err := h.eventStore.Delete(uid)
	if err != nil {
		errors.InternalError(err).Write(w, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
}

// ExportEventHandler handles exporting a single event to ICS
func (h *EventHandler) ExportEventHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/calendar; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=event.ics")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("test ics"))
}
