package webdav

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
)

// Handler wraps the standard WebDAV handler with CalDAV support
type Handler struct {
	fs             FileSystem
	jwtSecret      string
	enableAuth     bool
	corsMiddleware *cors.Cors
}

// Config holds WebDAV configuration
type Config struct {
	Host       string
	Port       string
	JWTSecret  string
	EnableAuth bool
	Production bool
}

// FileSystem is the interface for calendar filesystem operations
type FileSystem interface {
	MkCalendar(ctx context.Context, userID, calendarUID string, props map[string]string) error
	GetCalendar(ctx context.Context, userID, calendarUID string) (*Calendar, error)
	ListCalendars(ctx context.Context, userID string) ([]*Calendar, error)
	GetEvents(ctx context.Context, userID, calendarUID string, startTime, endTime *TimeRange) ([]*Event, error)
	GetEvent(ctx context.Context, userID, calendarUID, eventUID string) (*Event, error)
	CreateEvent(ctx context.Context, userID, calendarUID string, event *Event) error
	UpdateEvent(ctx context.Context, userID, calendarUID, eventUID string, event *Event) error
	DeleteEvent(ctx context.Context, userID, calendarUID, eventUID string) error
	GetPublicCalendar(ctx context.Context, hash string) (*Calendar, error)
	GetPublicEvents(ctx context.Context, hash string, startTime, endTime *TimeRange) ([]*Event, error)
}

// Calendar represents a calendar collection
type Calendar struct {
	UID         string
	OwnerID     string
	DisplayName string
	Description string
	Color       string
	IsPublic    bool
	PublicHash  string
	Timezone    string
}

// Event represents a calendar event
type Event struct {
	UID         string
	CalendarUID string
	DTStart     time.Time
	DTEnd       time.Time
	Summary     string
	Description string
	Location    string
	AllDay      bool
}

// TimeRange represents a date/time range for queries
type TimeRange struct {
	Start time.Time
	End   time.Time
}

// NewHandler creates a new WebDAV handler with CalDAV support
func NewHandler(fs FileSystem, cfg *Config) *Handler {
	h := &Handler{
		fs:         fs,
		jwtSecret:  cfg.JWTSecret,
		enableAuth: cfg.EnableAuth,
	}

	h.corsMiddleware = cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PROPFIND", "REPORT", "MKCOL", "COPY", "MOVE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type", "Depth", "Calendar-Timezone"},
		AllowCredentials: true,
		MaxAge:           3600,
	})

	return h
}

// ServeHTTP handles all WebDAV/CalDAV requests
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.corsMiddleware.Handler(http.HandlerFunc(h.handleRequest)).ServeHTTP(w, r)
}

func (h *Handler) handleRequest(w http.ResponseWriter, r *http.Request) {
	ctx, err := h.authenticate(r)
	if err != nil {
		logrus.WithError(err).Error("Authentication failed")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	r = r.WithContext(ctx)

	switch r.Method {
	case "OPTIONS":
		h.handleOptions(w, r)
	case "PROPFIND":
		h.handlePropfind(w, r)
	case "REPORT":
		h.handleReport(w, r)
	case "MKCOL", "MKCALENDAR":
		h.handleMkcalendar(w, r)
	case "GET":
		h.handleGet(w, r)
	case "PUT":
		h.handlePut(w, r)
	case "DELETE":
		h.handleDelete(w, r)
	case "COPY":
		h.handleCopy(w, r)
	case "MOVE":
		h.handleMove(w, r)
	case "PROPPATCH":
		h.handleProppatch(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) authenticate(r *http.Request) (context.Context, error) {
	if !h.enableAuth {
		return context.Background(), nil
	}

	authHeader := r.Header.Get("Authorization")

	// Try Bearer token
	if strings.HasPrefix(authHeader, "Bearer ") {
		_ = strings.TrimPrefix(authHeader, "Bearer ")
		// TODO: Verify JWT token using h.jwtSecret
		ctx := context.WithValue(r.Context(), "user_id", "user-from-token")
		return ctx, nil
	}

	// Try Basic auth
	username, _, ok := r.BasicAuth()
	if ok {
		// TODO: Verify credentials
		ctx := context.WithValue(r.Context(), "user_id", username)
		return ctx, nil
	}

	return nil, fmt.Errorf("no authentication provided")
}

func (h *Handler) handleOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Allow", "GET, POST, PUT, DELETE, PROPFIND, REPORT, MKCOL, COPY, MOVE")
	w.Header().Set("DAV", "1, 2, 3")
	w.Header().Set("MS-Author-Via", "DAV")
	h.corsMiddleware.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})).ServeHTTP(w, r)
}

func (h *Handler) handlePropfind(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value("user_id").(string)

	path := r.URL.Path
	var calendarUID string

	if strings.HasPrefix(path, "/webdav/") {
		parts := strings.Split(strings.TrimPrefix(path, "/webdav/"), "/")
		if len(parts) >= 2 {
			calendarUID = parts[1]
		}
	}

	calendar := &Calendar{
		UID:         calendarUID,
		OwnerID:     userID,
		DisplayName: calendarUID,
		Description: "Default calendar",
		Color:       "#FFFFFF",
		Timezone:    "UTC",
	}

	response := fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<D:multistatus xmlns:D="DAV:" xmlns:C="urn:ietf:params:xml:ns:caldav">
  <D:response>
    <D:href>%s</D:href>
    <D:propstat>
      <D:prop>
        <D:resourcetype><C:calendar-collection/></D:resourcetype>
        <C:calendar-description>%s</C:calendar-description>
        <C:calendar-color>%s</C:calendar-color>
        <C:calendar-timezone>%s</C:calendar-timezone>
        <D:getcontenttype>text/calendar</D:getcontenttype>
      </D:prop>
      <D:status>HTTP/1.1 200 OK</D:status>
    </D:propstat>
  </D:response>
</D:multistatus>`,
		r.URL.Path,
		calendar.Description,
		calendar.Color,
		calendar.Timezone,
	)

	w.Header().Set("Content-Type", "text/xml; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

func (h *Handler) handleReport(w http.ResponseWriter, r *http.Request) {
	// Read request body
	body, _ := io.ReadAll(r.Body)
	_ = body // TODO: Parse CalDAV query

	response := `<?xml version="1.0" encoding="utf-8"?>
<D:multistatus xmlns:D="DAV:" xmlns:C="urn:ietf:params:xml:ns:caldav">
  <D:response>
    <D:href/>
    <D:propstat>
      <D:prop/>
      <D:status>HTTP/1.1 207 No Content</D:status>
    </D:propstat>
  </D:response>
</D:multistatus>`

	w.Header().Set("Content-Type", "text/xml; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

func (h *Handler) handleMkcalendar(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	_ = ctx.Value("user_id").(string)

	calendarUID := h.extractCalendarUID(r.URL.Path)

	_ = map[string]string{
		"displayName": calendarUID,
		"description": "",
		"color":       "#FFFFFF",
		"timezone":    "UTC",
	}

	// TODO: Actually call fs.MkCalendar

	w.Header().Set("Content-Type", "text/xml; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`<?xml version="1.0" encoding="utf-8"?>
<D:multistatus xmlns:D="DAV:" xmlns:C="urn:ietf:params:xml:ns:caldav">
  <D:response>
    <D:href>` + r.URL.Path + `</D:href>
    <D:propstat>
      <D:prop>
        <D:resourcetype><C:calendar-collection/></D:resourcetype>
      </D:prop>
      <D:status>HTTP/1.1 201 Created</D:status>
    </D:propstat>
  </D:response>
</D:multistatus>`))
}

func (h *Handler) handleGet(w http.ResponseWriter, r *http.Request) {
	calendarUID, eventUID := h.parseCalendarPath(r.URL.Path)

	event := &Event{
		UID:         eventUID,
		CalendarUID: calendarUID,
		Summary:     "Sample Event",
		DTStart:     time.Now(),
		DTEnd:       time.Now().Add(time.Hour),
	}
	_ = event

	icsContent := h.generateICS(event)

	w.Header().Set("Content-Type", "text/calendar; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename="+event.UID+".ics")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(icsContent))
}

func (h *Handler) handlePut(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "text/calendar" && !strings.HasPrefix(contentType, "text/calendar") {
		http.Error(w, "Unsupported Media Type", http.StatusUnsupportedMediaType)
		return
	}

	buf := make([]byte, 4096)
	n, err := r.Body.Read(buf)
	if err != nil && err.Error() != "EOF" {
		logrus.WithError(err).Error("Failed to read ICS content")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	_ = buf[:n]
	_ = n // Use n to avoid unused warning

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`<?xml version="1.0" encoding="utf-8"?>
<D:response xmlns:D="DAV:">
  <D:href>` + r.URL.Path + `</D:href>
  <D:propstat>
    <D:prop/>
    <D:status>HTTP/1.1 201 Created</D:status>
  </D:propstat>
</D:response>`))
}

func (h *Handler) handleDelete(w http.ResponseWriter, r *http.Request) {
	calendarUID, eventUID := h.parseCalendarPath(r.URL.Path)
	_ = calendarUID
	_ = eventUID

	// TODO: Actually delete event

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) handleCopy(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *Handler) handleMove(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *Handler) handleProppatch(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *Handler) extractCalendarUID(path string) string {
	parts := strings.Split(strings.TrimPrefix(path, "/webdav/"), "/")
	if len(parts) >= 2 {
		return parts[1]
	}
	return ""
}

func (h *Handler) parseCalendarPath(path string) (calendarUID, itemUID string) {
	parts := strings.Split(strings.TrimPrefix(path, "/webdav/"), "/")
	if len(parts) >= 2 {
		calendarUID = parts[1]
	}
	if len(parts) >= 3 {
		itemUID = strings.TrimSuffix(parts[2], ".ics")
	}
	return
}

func (h *Handler) generateICS(event *Event) string {
	return fmt.Sprintf(`BEGIN:VCALENDAR
VERSION:2.0
PRODID:-//david//en
BEGIN:VEVENT
UID:%s
DTSTAMP:%s
DTSTART:%s
DTEND:%s
SUMMARY:%s
DESCRIPTION:%s
LOCATION:%s
END:VEVENT
END:VCALENDAR`,
		event.UID,
		time.Now().UTC().Format("20060102T150405Z"),
		event.DTStart.Format("20060102T150405Z"),
		event.DTEnd.Format("20060102T150405Z"),
		event.Summary,
		event.Description,
		event.Location,
	)
}
