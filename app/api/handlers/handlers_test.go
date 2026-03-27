package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/audstanley/david/app/api/errors"
)

// Auth Handler Tests

func TestAuthHandlerLoginSuccess(t *testing.T) {
	h := NewAuthHandler("test-secret")

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
	h := NewAuthHandler("test-secret")

	body := bytes.NewBufferString(`{"username":"admin","password":"wrongpassword"}`)
	req := httptest.NewRequest("POST", "/login", body)
	w := httptest.NewRecorder()

	h.LoginHandler(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestAuthHandlerLoginInvalidJSON(t *testing.T) {
	h := NewAuthHandler("test-secret")

	body := bytes.NewBufferString(`invalid json`)
	req := httptest.NewRequest("POST", "/login", body)
	w := httptest.NewRecorder()

	h.LoginHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestAuthHandlerLoginMissingFields(t *testing.T) {
	h := NewAuthHandler("test-secret")

	body := bytes.NewBufferString(`{"username":"admin"}`)
	req := httptest.NewRequest("POST", "/login", body)
	w := httptest.NewRecorder()

	h.LoginHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestAuthHandlerRefreshSuccess(t *testing.T) {
	h := NewAuthHandler("test-secret")

	body := bytes.NewBufferString(`{"refresh_token":"mock_refresh_test_20260324120000"}`)
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
	h := NewAuthHandler("test-secret")

	body := bytes.NewBufferString(`{"refresh_token":""}`)
	req := httptest.NewRequest("POST", "/refresh", body)
	w := httptest.NewRecorder()

	h.RefreshHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestAuthHandlerLogoutSuccess(t *testing.T) {
	h := NewAuthHandler("test-secret")

	req := httptest.NewRequest("POST", "/logout", nil)
	w := httptest.NewRecorder()

	h.LogoutHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestAuthHandlerVerifyValidToken(t *testing.T) {
	h := NewAuthHandler("test-secret")

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

	if !resp["valid"].(bool) {
		t.Error("Token should be valid")
	}
}

func TestAuthHandlerVerifyInvalidToken(t *testing.T) {
	h := NewAuthHandler("test-secret")

	req := httptest.NewRequest("GET", "/verify", nil)
	req.Header.Set("Authorization", "Bearer invalid_token")
	w := httptest.NewRecorder()

	h.VerifyHandler(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestAuthHandlerVerifyNoToken(t *testing.T) {
	h := NewAuthHandler("test-secret")

	req := httptest.NewRequest("GET", "/verify", nil)
	w := httptest.NewRecorder()

	h.VerifyHandler(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestAuthHandlerCreateAPIKey(t *testing.T) {
	h := NewAuthHandler("test-secret")

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
	h := NewAuthHandler("test-secret")

	body := bytes.NewBufferString(`{"name":"","permissions":["calendar.read"],"expiry_days":30}`)
	req := httptest.NewRequest("POST", "/api-keys", body)
	w := httptest.NewRecorder()

	h.CreateAPIKeyHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestAuthHandlerCreateAPIKeyInvalidExpiry(t *testing.T) {
	h := NewAuthHandler("test-secret")

	body := bytes.NewBufferString(`{"name":"Test Key","permissions":["calendar.read"],"expiry_days":0}`)
	req := httptest.NewRequest("POST", "/api-keys", body)
	w := httptest.NewRecorder()

	h.CreateAPIKeyHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestAuthHandlerListAPIKeys(t *testing.T) {
	h := NewAuthHandler("test-secret")

	req := httptest.NewRequest("GET", "/api-keys", nil)
	w := httptest.NewRecorder()

	h.ListAPIKeyHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestAuthHandlerRevokeAPIKey(t *testing.T) {
	h := NewAuthHandler("test-secret")

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
	h := NewCalendarHandler()

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
	h := NewCalendarHandler()

	req := httptest.NewRequest("GET", "/calendars/cal-1", nil)
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
	h := NewCalendarHandler()

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
	h := NewCalendarHandler()

	body := bytes.NewBufferString(`{"description":"Test Desc"}`)
	req := httptest.NewRequest("POST", "/calendars", body)
	w := httptest.NewRecorder()

	h.CreateCalendarHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestCalendarHandlerCreateInvalidTimezone(t *testing.T) {
	h := NewCalendarHandler()

	body := bytes.NewBufferString(`{"display_name":"Test Calendar","timezone":"nonexistent-tz"}`)
	req := httptest.NewRequest("POST", "/calendars", body)
	w := httptest.NewRecorder()

	h.CreateCalendarHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestCalendarHandlerCreateDisplayNameTooLong(t *testing.T) {
	h := NewCalendarHandler()

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
	h := NewCalendarHandler()

	req := httptest.NewRequest("PUT", "/calendars/cal-1", nil)
	w := httptest.NewRecorder()

	h.UpdateCalendarHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestCalendarHandlerDeleteSuccess(t *testing.T) {
	h := NewCalendarHandler()

	req := httptest.NewRequest("DELETE", "/calendars/cal-1", nil)
	w := httptest.NewRecorder()

	h.DeleteCalendarHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestCalendarHandlerGrantShare(t *testing.T) {
	h := NewCalendarHandler()

	body := bytes.NewBufferString(`{"user_id":"user-123","role":"read"}`)
	req := httptest.NewRequest("POST", "/calendars/cal-1/shares", body)
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
	h := NewCalendarHandler()

	body := bytes.NewBufferString(`{"user_id":"user-123","role":"invalid"}`)
	req := httptest.NewRequest("POST", "/calendars/cal-1/shares", body)
	w := httptest.NewRecorder()

	h.GrantShareHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestCalendarHandlerRevokeShare(t *testing.T) {
	h := NewCalendarHandler()

	req := httptest.NewRequest("DELETE", "/calendars/cal-1/shares/1", nil)
	w := httptest.NewRecorder()

	h.RevokeShareHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestCalendarHandlerExportSuccess(t *testing.T) {
	h := NewCalendarHandler()

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
	h := NewCalendarHandler()

	req := httptest.NewRequest("GET", "/calendars/cal-1/stats", nil)
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
	h := NewEventHandler()

	req := httptest.NewRequest("GET", "/events", nil)
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
	h := NewEventHandler()

	req := httptest.NewRequest("GET", "/events/event-1", nil)
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
	h := NewEventHandler()

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
	h := NewEventHandler()

	body := bytes.NewBufferString(`{"dtstart":"2026-03-25T14:00:00Z"}`)
	req := httptest.NewRequest("POST", "/events", body)
	w := httptest.NewRecorder()

	h.CreateEventHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestEventHandlerCreateMissingDTStart(t *testing.T) {
	h := NewEventHandler()

	body := bytes.NewBufferString(`{"summary":"Test Event"}`)
	req := httptest.NewRequest("POST", "/events", body)
	w := httptest.NewRecorder()

	h.CreateEventHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestEventHandlerUpdateSuccess(t *testing.T) {
	h := NewEventHandler()

	req := httptest.NewRequest("PUT", "/events/event-1", nil)
	w := httptest.NewRecorder()

	h.UpdateEventHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestEventHandlerDeleteSuccess(t *testing.T) {
	h := NewEventHandler()

	req := httptest.NewRequest("DELETE", "/events/event-1", nil)
	w := httptest.NewRecorder()

	h.DeleteEventHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestEventHandlerExportSuccess(t *testing.T) {
	h := NewEventHandler()

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
