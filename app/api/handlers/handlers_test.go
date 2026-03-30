package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/audstanley/david/app/api/errors"
	"github.com/audstanley/david/app/storage"
	"github.com/audstanley/david/app/storage/calendars"
	"github.com/audstanley/david/app/storage/events"
	"github.com/audstanley/david/app/storage/users"
)

// createMockStorage creates a minimal storage for unit tests
func createMockStorage(t *testing.T) *storage.Storage {
	tmpDir := "/tmp/david-test-" + fmt.Sprintf("%d", time.Now().UnixNano()) + "-" + t.Name()
	db, err := storage.NewStorage(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}
	return db
}

// createTestCalendar creates a test calendar in the given storage
func createTestCalendar(t *testing.T, db *storage.Storage) {
	calendarStore := calendars.New(db)
	err := calendarStore.Create("cal-1", "user-1", "Test Calendar", "Test Description", "#ff0000", "UTC", false)
	if err != nil {
		t.Fatalf("Failed to create test calendar: %v", err)
	}
}

// createTestEvent creates a test event in the given storage
func createTestEvent(t *testing.T, db *storage.Storage) {
	eventStore := events.New(db)
	dtStart, _ := time.Parse(time.RFC3339, "2026-03-25T14:00:00Z")
	dtEnd, _ := time.Parse(time.RFC3339, "2026-03-25T15:00:00Z")
	event := &storage.Event{
		UID:         "event-1",
		CalendarUID: "cal-1",
		Summary:     "Test Event",
		Description: "Test Description",
		Location:    "Test Location",
		DTStart:     dtStart,
		DTEnd:       dtEnd,
	}
	err := eventStore.Create(event)
	if err != nil {
		t.Fatalf("Failed to create test event: %v", err)
	}
}

// Auth Handler Tests

func TestAuthHandlerLoginSuccess(t *testing.T) {
	db := createMockStorage(t)
	h := NewAuthHandler("test-secret", db)

	// Create test user
	userStore := users.New(db)
	hash, err := userStore.HashPassword("password")
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}
	err = userStore.Create("user-1", "admin", hash, "admin@test.com", "Admin User", "admin")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	body := bytes.NewBufferString(`{"username":"admin","password":"password"}`)
	req := httptest.NewRequest("POST", "/login", body)
	w := httptest.NewRecorder()

	h.LoginHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if resp["token_type"] != "Bearer" {
		t.Errorf("Expected token_type 'Bearer', got %v", resp["token_type"])
	}
}

func TestAuthHandlerLoginInvalidCredentials(t *testing.T) {
	db := createMockStorage(t)
	h := NewAuthHandler("test-secret", db)

	body := bytes.NewBufferString(`{"username":"admin","password":"wrongpassword"}`)
	req := httptest.NewRequest("POST", "/login", body)
	w := httptest.NewRecorder()

	h.LoginHandler(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestAuthHandlerLoginInvalidJSON(t *testing.T) {
	db := createMockStorage(t)
	h := NewAuthHandler("test-secret", db)

	body := bytes.NewBufferString(`invalid json`)
	req := httptest.NewRequest("POST", "/login", body)
	w := httptest.NewRecorder()

	h.LoginHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestAuthHandlerLoginMissingFields(t *testing.T) {
	db := createMockStorage(t)
	h := NewAuthHandler("test-secret", db)

	body := bytes.NewBufferString(`{"username":"admin"}`)
	req := httptest.NewRequest("POST", "/login", body)
	w := httptest.NewRecorder()

	h.LoginHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestAuthHandlerRefreshSuccess(t *testing.T) {
	db := createMockStorage(t)
	h := NewAuthHandler("test-secret", db)

	// Create test user
	userStore := users.New(db)
	hash, _ := userStore.HashPassword("password")
	err := userStore.Create("user-1", "testuser", hash, "test@test.com", "Test User", "user")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Login to get valid refresh token
	loginBody := bytes.NewBufferString(`{"username":"testuser","password":"password"}`)
	loginReq := httptest.NewRequest("POST", "/login", loginBody)
	loginW := httptest.NewRecorder()
	h.LoginHandler(loginW, loginReq)

	if loginW.Code != http.StatusOK {
		t.Fatalf("Login failed: %d", loginW.Code)
	}

	var loginResp map[string]interface{}
	if err := json.NewDecoder(loginW.Body).Decode(&loginResp); err != nil {
		t.Fatalf("Failed to decode login response: %v", err)
	}

	refreshToken, ok := loginResp["refresh_token"].(string)
	if !ok {
		t.Fatal("No refresh token in response")
	}

	body := bytes.NewBufferString(fmt.Sprintf(`{"refresh_token":"%s"}`, refreshToken))
	req := httptest.NewRequest("POST", "/refresh", body)
	w := httptest.NewRecorder()

	h.RefreshHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if _, ok := resp["access_token"]; !ok {
		t.Error("Response should contain access_token")
	}
}

func TestAuthHandlerRefreshInvalidToken(t *testing.T) {
	db := createMockStorage(t)
	h := NewAuthHandler("test-secret", db)

	body := bytes.NewBufferString(`{"refresh_token":""}`)
	req := httptest.NewRequest("POST", "/refresh", body)
	w := httptest.NewRecorder()

	h.RefreshHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestAuthHandlerLogoutSuccess(t *testing.T) {
	db := createMockStorage(t)
	h := NewAuthHandler("test-secret", db)

	req := httptest.NewRequest("POST", "/logout", nil)
	w := httptest.NewRecorder()

	h.LogoutHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestAuthHandlerVerifyValidToken(t *testing.T) {
	db := createMockStorage(t)
	h := NewAuthHandler("test-secret", db)

	// Create test user
	userStore := users.New(db)
	hash, _ := userStore.HashPassword("password")
	err := userStore.Create("user-1", "testuser", hash, "test@test.com", "Test User", "user")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	req := httptest.NewRequest("GET", "/verify", nil)
	req.Header.Set("Authorization", "Bearer mock_user-1_user_20260324120000")
	w := httptest.NewRecorder()

	h.VerifyHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	valid, ok := resp["valid"].(bool)
	if !ok || !valid {
		t.Error("Token should be valid")
	}
}

func TestAuthHandlerVerifyInvalidToken(t *testing.T) {
	db := createMockStorage(t)
	h := NewAuthHandler("test-secret", db)

	req := httptest.NewRequest("GET", "/verify", nil)
	req.Header.Set("Authorization", "Bearer invalid_token")
	w := httptest.NewRecorder()

	h.VerifyHandler(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestAuthHandlerVerifyNoToken(t *testing.T) {
	db := createMockStorage(t)
	h := NewAuthHandler("test-secret", db)

	req := httptest.NewRequest("GET", "/verify", nil)
	w := httptest.NewRecorder()

	h.VerifyHandler(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestAuthHandlerCreateAPIKey(t *testing.T) {
	db := createMockStorage(t)
	h := NewAuthHandler("test-secret", db)

	body := bytes.NewBufferString(`{"name":"Test Key","permissions":["calendar.read"],"expiry_days":30}`)
	req := httptest.NewRequest("POST", "/api-keys", body)
	w := httptest.NewRecorder()

	h.CreateAPIKeyHandler(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if resp["name"] != "Test Key" {
		t.Errorf("Expected name 'Test Key', got %v", resp["name"])
	}
}

func TestAuthHandlerCreateAPIKeyInvalidName(t *testing.T) {
	db := createMockStorage(t)
	h := NewAuthHandler("test-secret", db)

	body := bytes.NewBufferString(`{"name":"","permissions":["calendar.read"],"expiry_days":30}`)
	req := httptest.NewRequest("POST", "/api-keys", body)
	w := httptest.NewRecorder()

	h.CreateAPIKeyHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestAuthHandlerCreateAPIKeyInvalidExpiry(t *testing.T) {
	db := createMockStorage(t)
	h := NewAuthHandler("test-secret", db)

	body := bytes.NewBufferString(`{"name":"Test Key","permissions":["calendar.read"],"expiry_days":0}`)
	req := httptest.NewRequest("POST", "/api-keys", body)
	w := httptest.NewRecorder()

	h.CreateAPIKeyHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestAuthHandlerListAPIKeys(t *testing.T) {
	db := createMockStorage(t)
	h := NewAuthHandler("test-secret", db)

	req := httptest.NewRequest("GET", "/api-keys", nil)
	w := httptest.NewRecorder()

	h.ListAPIKeyHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestAuthHandlerRevokeAPIKey(t *testing.T) {
	db := createMockStorage(t)
	h := NewAuthHandler("test-secret", db)

	req := httptest.NewRequest("DELETE", "/api-keys/1", nil)
	w := httptest.NewRecorder()

	h.RevokeAPIKeyHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

// ParseBasicAuth Tests

func TestParseBasicAuthValid(t *testing.T) {
	username, password, ok := ParseBasicAuth("Basic YWRtaW46cGFzc3dvcmQ=") // admin:password

	if !ok {
		t.Error("Expected parsing to succeed")
	}
	if username != "admin" {
		t.Errorf("Expected username 'admin', got %q", username)
	}
	if password != "password" {
		t.Errorf("Expected password 'password', got %q", password)
	}
}

func TestParseBasicAuthInvalid(t *testing.T) {
	_, _, ok := ParseBasicAuth("Invalid")

	if ok {
		t.Error("Expected parsing to fail")
	}
}

func TestParseBasicAuthNoBasicPrefix(t *testing.T) {
	_, _, ok := ParseBasicAuth("Bearer token")

	if ok {
		t.Error("Expected parsing to fail")
	}
}

// Calendar Handler Tests

func TestCalendarHandlerListSuccess(t *testing.T) {
	db := createMockStorage(t)
	h := NewCalendarHandler(db)

	// Create a test calendar for unit tests
	calendarStore := calendars.New(db)
	err := calendarStore.Create("cal-1", "user-1", "Test Calendar", "Test Desc", "#ff0000", "UTC", false)
	if err != nil {
		t.Fatalf("Failed to create test calendar: %v", err)
	}

	req := httptest.NewRequest("GET", "/calendars", nil)
	w := httptest.NewRecorder()

	h.ListCalendarsHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var calendars []map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&calendars); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(calendars) != 1 {
		t.Errorf("Expected 1 calendar, got %d", len(calendars))
	}
}

func TestCalendarHandlerGetSuccess(t *testing.T) {
	db := createMockStorage(t)
	h := NewCalendarHandler(db)

	// Create test calendar
	calendarStore := calendars.New(db)
	err := calendarStore.Create("cal-1", "user-1", "Test Calendar", "Test Description", "#ff0000", "UTC", false)
	if err != nil {
		t.Fatalf("Failed to create test calendar: %v", err)
	}

	req := httptest.NewRequest("GET", "/calendars?uid=cal-1", nil)
	w := httptest.NewRecorder()

	h.GetCalendarHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var cal map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&cal); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if cal["uid"] != "cal-1" {
		t.Errorf("Expected UID 'cal-1', got %v", cal["uid"])
	}
}

func TestCalendarHandlerCreateSuccess(t *testing.T) {
	db := createMockStorage(t)
	h := NewCalendarHandler(db)

	body := bytes.NewBufferString(`{"display_name":"Test Calendar","description":"Test Desc","color":"#ff0000","timezone":"UTC","is_public":false}`)
	req := httptest.NewRequest("POST", "/calendars", body)
	w := httptest.NewRecorder()

	h.CreateCalendarHandler(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	var cal map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&cal); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if cal["display_name"] != "Test Calendar" {
		t.Errorf("Expected display_name 'Test Calendar', got %v", cal["display_name"])
	}
}

func TestCalendarHandlerCreateMissingDisplayName(t *testing.T) {
	db := createMockStorage(t)
	h := NewCalendarHandler(db)

	body := bytes.NewBufferString(`{"description":"Test Desc"}`)
	req := httptest.NewRequest("POST", "/calendars", body)
	w := httptest.NewRecorder()

	h.CreateCalendarHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestCalendarHandlerCreateInvalidTimezone(t *testing.T) {
	db := createMockStorage(t)
	h := NewCalendarHandler(db)

	body := bytes.NewBufferString(`{"display_name":"Test Calendar","timezone":"nonexistent-tz"}`)
	req := httptest.NewRequest("POST", "/calendars", body)
	w := httptest.NewRecorder()

	h.CreateCalendarHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestCalendarHandlerCreateDisplayNameTooLong(t *testing.T) {
	db := createMockStorage(t)
	h := NewCalendarHandler(db)

	longName := ""
	for i := 0; i < 201; i++ {
		longName += "a"
	}

	body := bytes.NewBufferString(`{"display_name":"` + longName + `"}`)
	req := httptest.NewRequest("POST", "/calendars", body)
	w := httptest.NewRecorder()

	h.CreateCalendarHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestCalendarHandlerUpdateSuccess(t *testing.T) {
	db := createMockStorage(t)
	h := NewCalendarHandler(db)

	createTestCalendar(t, db)

	body := bytes.NewBufferString(`{"display_name":"Updated Calendar"}`)
	req := httptest.NewRequest("PUT", "/calendars?uid=cal-1", body)
	w := httptest.NewRecorder()

	h.UpdateCalendarHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestCalendarHandlerDeleteSuccess(t *testing.T) {
	db := createMockStorage(t)
	h := NewCalendarHandler(db)

	createTestCalendar(t, db)

	req := httptest.NewRequest("DELETE", "/calendars?uid=cal-1", nil)
	w := httptest.NewRecorder()

	h.DeleteCalendarHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestCalendarHandlerGrantShare(t *testing.T) {
	db := createMockStorage(t)
	h := NewCalendarHandler(db)

	createTestCalendar(t, db)

	body := bytes.NewBufferString(`{"user_id":"user-123","role":"read"}`)
	req := httptest.NewRequest("POST", "/calendars?uid=cal-1/shares", body)
	w := httptest.NewRecorder()

	h.GrantShareHandler(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	var share map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&share); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if share["user_id"] != "user-123" {
		t.Errorf("Expected user_id 'user-123', got %v", share["user_id"])
	}
}

func TestCalendarHandlerGrantShareInvalidRole(t *testing.T) {
	db := createMockStorage(t)
	h := NewCalendarHandler(db)

	body := bytes.NewBufferString(`{"user_id":"user-123","role":"invalid"}`)
	req := httptest.NewRequest("POST", "/calendars/cal-1/shares", body)
	w := httptest.NewRecorder()

	h.GrantShareHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestCalendarHandlerRevokeShare(t *testing.T) {
	db := createMockStorage(t)
	h := NewCalendarHandler(db)

	createTestCalendar(t, db)

	req := httptest.NewRequest("DELETE", "/calendars?uid=cal-1&share_id=1", nil)
	w := httptest.NewRecorder()

	h.RevokeShareHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestCalendarHandlerExportSuccess(t *testing.T) {
	db := createMockStorage(t)
	h := NewCalendarHandler(db)

	req := httptest.NewRequest("GET", "/calendars/cal-1/export", nil)
	w := httptest.NewRecorder()

	h.ExportCalendarHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "text/calendar; charset=utf-8" {
		t.Errorf("Expected Content-Type 'text/calendar; charset=utf-8', got %q", contentType)
	}
}

func TestCalendarHandlerGetStats(t *testing.T) {
	db := createMockStorage(t)
	h := NewCalendarHandler(db)

	createTestCalendar(t, db)

	req := httptest.NewRequest("GET", "/calendars?uid=cal-1&endpoint=stats", nil)
	w := httptest.NewRecorder()

	h.GetCalendarStatsHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var stats map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&stats); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if stats["calendar_uid"] != "cal-1" {
		t.Errorf("Expected calendar_uid 'cal-1', got %v", stats["calendar_uid"])
	}
}

// Event Handler Tests

func TestEventHandlerListSuccess(t *testing.T) {
	db := createMockStorage(t)
	h := NewEventHandler(db)

	createTestCalendar(t, db)
	createTestEvent(t, db)

	req := httptest.NewRequest("GET", "/events?calendar_uid=cal-1&start=2026-03-25T00:00:00Z&end=2026-03-25T23:59:59Z", nil)
	w := httptest.NewRecorder()

	h.ListEventsHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var events []map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&events); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(events))
	}
}

func TestEventHandlerGetSuccess(t *testing.T) {
	db := createMockStorage(t)
	h := NewEventHandler(db)

	createTestCalendar(t, db)
	createTestEvent(t, db)

	req := httptest.NewRequest("GET", "/events?uid=event-1", nil)
	w := httptest.NewRecorder()

	h.GetEventHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var event map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&event); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if event["uid"] != "event-1" {
		t.Errorf("Expected UID 'event-1', got %v", event["uid"])
	}
}

func TestEventHandlerCreateSuccess(t *testing.T) {
	db := createMockStorage(t)
	h := NewEventHandler(db)

	body := bytes.NewBufferString(`{"summary":"Test Event","description":"Test Desc","dtstart":"2026-03-25T14:00:00Z","dtend":"2026-03-25T15:00:00Z"}`)
	req := httptest.NewRequest("POST", "/events", body)
	w := httptest.NewRecorder()

	h.CreateEventHandler(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	var event map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&event); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if event["summary"] != "Test Event" {
		t.Errorf("Expected summary 'Test Event', got %v", event["summary"])
	}
}

func TestEventHandlerCreateMissingSummary(t *testing.T) {
	db := createMockStorage(t)
	h := NewEventHandler(db)

	body := bytes.NewBufferString(`{"dtstart":"2026-03-25T14:00:00Z"}`)
	req := httptest.NewRequest("POST", "/events", body)
	w := httptest.NewRecorder()

	h.CreateEventHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestEventHandlerCreateMissingDTStart(t *testing.T) {
	db := createMockStorage(t)
	h := NewEventHandler(db)

	body := bytes.NewBufferString(`{"summary":"Test Event"}`)
	req := httptest.NewRequest("POST", "/events", body)
	w := httptest.NewRecorder()

	h.CreateEventHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestEventHandlerUpdateSuccess(t *testing.T) {
	db := createMockStorage(t)
	h := NewEventHandler(db)

	createTestCalendar(t, db)
	createTestEvent(t, db)

	body := bytes.NewBufferString(`{"summary":"Updated Event"}`)
	req := httptest.NewRequest("PUT", "/events?uid=event-1", body)
	w := httptest.NewRecorder()

	h.UpdateEventHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestEventHandlerDeleteSuccess(t *testing.T) {
	db := createMockStorage(t)
	h := NewEventHandler(db)

	createTestCalendar(t, db)
	createTestEvent(t, db)

	req := httptest.NewRequest("DELETE", "/events?uid=event-1", nil)
	w := httptest.NewRecorder()

	h.DeleteEventHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestEventHandlerExportSuccess(t *testing.T) {
	db := createMockStorage(t)
	h := NewEventHandler(db)

	req := httptest.NewRequest("GET", "/events/event-1/export", nil)
	w := httptest.NewRecorder()

	h.ExportEventHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "text/calendar; charset=utf-8" {
		t.Errorf("Expected Content-Type 'text/calendar; charset=utf-8', got %q", contentType)
	}
}

// Error Handler Tests

func TestValidationErrorWrite(t *testing.T) {
	// Test that ValidationError writes correct status
	err := errors.ValidationError("field", "invalid value")

	// Create a test response writer
	w := httptest.NewRecorder()
	err.Write(w, http.StatusBadRequest)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestUnauthorizedWrite(t *testing.T) {
	err := errors.Unauthorized("invalid credentials")

	w := httptest.NewRecorder()
	err.Write(w, http.StatusUnauthorized)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}
