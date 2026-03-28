package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/audstanley/david/app"
	"github.com/audstanley/david/app/api/handlers"
	"github.com/audstanley/david/app/storage"
	"github.com/audstanley/david/app/storage/users"
	"github.com/audstanley/david/app/test"
)

// TestServer wraps a test server with storage and handlers
type TestServer struct {
	Storage         *storage.Storage
	Dir             string
	AuthHandler     *handlers.AuthHandler
	CalendarHandler *handlers.CalendarHandler
	EventHandler    *handlers.EventHandler
}

// SetupTestServer creates a test server with real storage
func SetupTestServer(t *testing.T) *TestServer {
	storageInstance, dir := test.CreateTempStorage(t)

	// Create a test user
	userStore := users.New(storageInstance)
	hash := app.GenHash([]byte("password"))
	err := userStore.Create("user-1", "admin", hash, "admin@test.com", "Admin User", "admin")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	return &TestServer{
		Storage:         storageInstance,
		Dir:             dir,
		AuthHandler:     handlers.NewAuthHandler("test-secret-key"),
		CalendarHandler: handlers.NewCalendarHandler(),
		EventHandler:    handlers.NewEventHandler(),
	}
}

// CleanupTestServer cleans up test resources
func CleanupTestServer(t *testing.T, ts *TestServer) {
	if err := test.CleanupStorage(ts.Dir); err != nil {
		t.Logf("Failed to cleanup storage: %v", err)
	}
}

// Auth Handler Integration Tests

func TestAuthHandlerLoginSuccess(t *testing.T) {
	ts := SetupTestServer(t)
	defer CleanupTestServer(t, ts)

	body := bytes.NewBufferString(`{"username":"admin","password":"password"}`)
	req := httptest.NewRequest("POST", "/login", body)
	w := httptest.NewRecorder()

	ts.AuthHandler.LoginHandler(w, req)

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

	user, ok := resp["user"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected user in response")
	}
	if user["username"] != "admin" {
		t.Errorf("Expected username 'admin', got %v", user["username"])
	}
}

func TestAuthHandlerLoginInvalidCredentials(t *testing.T) {
	ts := SetupTestServer(t)
	defer CleanupTestServer(t, ts)

	body := bytes.NewBufferString(`{"username":"admin","password":"wrongpassword"}`)
	req := httptest.NewRequest("POST", "/login", body)
	w := httptest.NewRecorder()

	ts.AuthHandler.LoginHandler(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestAuthHandlerLoginInvalidJSON(t *testing.T) {
	ts := SetupTestServer(t)
	defer CleanupTestServer(t, ts)

	body := bytes.NewBufferString(`invalid json`)
	req := httptest.NewRequest("POST", "/login", body)
	w := httptest.NewRecorder()

	ts.AuthHandler.LoginHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestAuthHandlerLoginMissingFields(t *testing.T) {
	ts := SetupTestServer(t)
	defer CleanupTestServer(t, ts)

	body := bytes.NewBufferString(`{"username":"admin"}`)
	req := httptest.NewRequest("POST", "/login", body)
	w := httptest.NewRecorder()

	ts.AuthHandler.LoginHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestAuthHandlerRefreshSuccess(t *testing.T) {
	ts := SetupTestServer(t)
	defer CleanupTestServer(t, ts)

	body := bytes.NewBufferString(`{"refresh_token":"mock_refresh_test_20260324120000"}`)
	req := httptest.NewRequest("POST", "/refresh", body)
	w := httptest.NewRecorder()

	ts.AuthHandler.RefreshHandler(w, req)

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
	ts := SetupTestServer(t)
	defer CleanupTestServer(t, ts)

	body := bytes.NewBufferString(`{"refresh_token":""}`)
	req := httptest.NewRequest("POST", "/refresh", body)
	w := httptest.NewRecorder()

	ts.AuthHandler.RefreshHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestAuthHandlerLogoutSuccess(t *testing.T) {
	ts := SetupTestServer(t)
	defer CleanupTestServer(t, ts)

	req := httptest.NewRequest("POST", "/logout", nil)
	w := httptest.NewRecorder()

	ts.AuthHandler.LogoutHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestAuthHandlerVerifyValidToken(t *testing.T) {
	ts := SetupTestServer(t)
	defer CleanupTestServer(t, ts)

	req := httptest.NewRequest("GET", "/verify", nil)
	req.Header.Set("Authorization", "Bearer mock_user-1_user_20260324120000")
	w := httptest.NewRecorder()

	ts.AuthHandler.VerifyHandler(w, req)

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
	ts := SetupTestServer(t)
	defer CleanupTestServer(t, ts)

	req := httptest.NewRequest("GET", "/verify", nil)
	req.Header.Set("Authorization", "Bearer invalid_token")
	w := httptest.NewRecorder()

	ts.AuthHandler.VerifyHandler(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestAuthHandlerCreateAPIKey(t *testing.T) {
	ts := SetupTestServer(t)
	defer CleanupTestServer(t, ts)

	body := bytes.NewBufferString(`{"name":"Test Key","permissions":["calendar.read"],"expiry_days":30}`)
	req := httptest.NewRequest("POST", "/api-keys", body)
	w := httptest.NewRecorder()

	ts.AuthHandler.CreateAPIKeyHandler(w, req)

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
	ts := SetupTestServer(t)
	defer CleanupTestServer(t, ts)

	body := bytes.NewBufferString(`{"name":"","permissions":["calendar.read"],"expiry_days":30}`)
	req := httptest.NewRequest("POST", "/api-keys", body)
	w := httptest.NewRecorder()

	ts.AuthHandler.CreateAPIKeyHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestAuthHandlerCreateAPIKeyInvalidExpiry(t *testing.T) {
	ts := SetupTestServer(t)
	defer CleanupTestServer(t, ts)

	body := bytes.NewBufferString(`{"name":"Test Key","permissions":["calendar.read"],"expiry_days":0}`)
	req := httptest.NewRequest("POST", "/api-keys", body)
	w := httptest.NewRecorder()

	ts.AuthHandler.CreateAPIKeyHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

// Calendar Handler Integration Tests

func TestCalendarHandlerListSuccess(t *testing.T) {
	ts := SetupTestServer(t)
	defer CleanupTestServer(t, ts)

	req := httptest.NewRequest("GET", "/calendars", nil)
	w := httptest.NewRecorder()

	ts.CalendarHandler.ListCalendarsHandler(w, req)

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

func TestCalendarHandlerCreateSuccess(t *testing.T) {
	ts := SetupTestServer(t)
	defer CleanupTestServer(t, ts)

	body := bytes.NewBufferString(`{"display_name":"Test Calendar","description":"Test Desc","color":"#ff0000","timezone":"UTC","is_public":false}`)
	req := httptest.NewRequest("POST", "/calendars", body)
	w := httptest.NewRecorder()

	ts.CalendarHandler.CreateCalendarHandler(w, req)

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
	ts := SetupTestServer(t)
	defer CleanupTestServer(t, ts)

	body := bytes.NewBufferString(`{"description":"Test Desc"}`)
	req := httptest.NewRequest("POST", "/calendars", body)
	w := httptest.NewRecorder()

	ts.CalendarHandler.CreateCalendarHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestCalendarHandlerCreateInvalidTimezone(t *testing.T) {
	ts := SetupTestServer(t)
	defer CleanupTestServer(t, ts)

	body := bytes.NewBufferString(`{"display_name":"Test Calendar","timezone":"nonexistent-tz"}`)
	req := httptest.NewRequest("POST", "/calendars", body)
	w := httptest.NewRecorder()

	ts.CalendarHandler.CreateCalendarHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestCalendarHandlerCreateDisplayNameTooLong(t *testing.T) {
	ts := SetupTestServer(t)
	defer CleanupTestServer(t, ts)

	longName := ""
	for i := 0; i < 201; i++ {
		longName += "a"
	}

	body := bytes.NewBufferString(`{"display_name":"` + longName + `"}`)
	req := httptest.NewRequest("POST", "/calendars", body)
	w := httptest.NewRecorder()

	ts.CalendarHandler.CreateCalendarHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestCalendarHandlerUpdateSuccess(t *testing.T) {
	ts := SetupTestServer(t)
	defer CleanupTestServer(t, ts)

	body := bytes.NewBufferString(`{"display_name":"Updated Calendar"}`)
	req := httptest.NewRequest("PUT", "/calendars/cal-1", body)
	w := httptest.NewRecorder()

	ts.CalendarHandler.UpdateCalendarHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestCalendarHandlerDeleteSuccess(t *testing.T) {
	ts := SetupTestServer(t)
	defer CleanupTestServer(t, ts)

	req := httptest.NewRequest("DELETE", "/calendars/cal-1", nil)
	w := httptest.NewRecorder()

	ts.CalendarHandler.DeleteCalendarHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestCalendarHandlerGrantShare(t *testing.T) {
	ts := SetupTestServer(t)
	defer CleanupTestServer(t, ts)

	body := bytes.NewBufferString(`{"user_id":"user-123","role":"read"}`)
	req := httptest.NewRequest("POST", "/calendars/cal-1/shares", body)
	w := httptest.NewRecorder()

	ts.CalendarHandler.GrantShareHandler(w, req)

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
	ts := SetupTestServer(t)
	defer CleanupTestServer(t, ts)

	body := bytes.NewBufferString(`{"user_id":"user-123","role":"invalid"}`)
	req := httptest.NewRequest("POST", "/calendars/cal-1/shares", body)
	w := httptest.NewRecorder()

	ts.CalendarHandler.GrantShareHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestCalendarHandlerExportSuccess(t *testing.T) {
	ts := SetupTestServer(t)
	defer CleanupTestServer(t, ts)

	req := httptest.NewRequest("GET", "/calendars/cal-1/export", nil)
	w := httptest.NewRecorder()

	ts.CalendarHandler.ExportCalendarHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "text/calendar; charset=utf-8" {
		t.Errorf("Expected Content-Type 'text/calendar; charset=utf-8', got %q", contentType)
	}
}

func TestCalendarHandlerGetStats(t *testing.T) {
	ts := SetupTestServer(t)
	defer CleanupTestServer(t, ts)

	req := httptest.NewRequest("GET", "/calendars/cal-1/stats", nil)
	w := httptest.NewRecorder()

	ts.CalendarHandler.GetCalendarStatsHandler(w, req)

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

// Event Handler Integration Tests

func TestEventHandlerListSuccess(t *testing.T) {
	ts := SetupTestServer(t)
	defer CleanupTestServer(t, ts)

	req := httptest.NewRequest("GET", "/events", nil)
	w := httptest.NewRecorder()

	ts.EventHandler.ListEventsHandler(w, req)

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

func TestEventHandlerCreateSuccess(t *testing.T) {
	ts := SetupTestServer(t)
	defer CleanupTestServer(t, ts)

	body := bytes.NewBufferString(`{"summary":"Test Event","description":"Test Desc","dtstart":"2026-03-25T14:00:00Z","dtend":"2026-03-25T15:00:00Z"}`)
	req := httptest.NewRequest("POST", "/events", body)
	w := httptest.NewRecorder()

	ts.EventHandler.CreateEventHandler(w, req)

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
	ts := SetupTestServer(t)
	defer CleanupTestServer(t, ts)

	body := bytes.NewBufferString(`{"dtstart":"2026-03-25T14:00:00Z"}`)
	req := httptest.NewRequest("POST", "/events", body)
	w := httptest.NewRecorder()

	ts.EventHandler.CreateEventHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestEventHandlerCreateMissingDTStart(t *testing.T) {
	ts := SetupTestServer(t)
	defer CleanupTestServer(t, ts)

	body := bytes.NewBufferString(`{"summary":"Test Event"}`)
	req := httptest.NewRequest("POST", "/events", body)
	w := httptest.NewRecorder()

	ts.EventHandler.CreateEventHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestEventHandlerUpdateSuccess(t *testing.T) {
	ts := SetupTestServer(t)
	defer CleanupTestServer(t, ts)

	body := bytes.NewBufferString(`{"summary":"Updated Event"}`)
	req := httptest.NewRequest("PUT", "/events/event-1", body)
	w := httptest.NewRecorder()

	ts.EventHandler.UpdateEventHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestEventHandlerDeleteSuccess(t *testing.T) {
	ts := SetupTestServer(t)
	defer CleanupTestServer(t, ts)

	req := httptest.NewRequest("DELETE", "/events/event-1", nil)
	w := httptest.NewRecorder()

	ts.EventHandler.DeleteEventHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestEventHandlerExportSuccess(t *testing.T) {
	ts := SetupTestServer(t)
	defer CleanupTestServer(t, ts)

	req := httptest.NewRequest("GET", "/events/event-1/export", nil)
	w := httptest.NewRecorder()

	ts.EventHandler.ExportEventHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "text/calendar; charset=utf-8" {
		t.Errorf("Expected Content-Type 'text/calendar; charset=utf-8', got %q", contentType)
	}
}

// ParseBasicAuth Tests

func TestParseBasicAuthValid(t *testing.T) {
	username, password, ok := handlers.ParseBasicAuth("Basic YWRtaW46cGFzc3dvcmQ=") // admin:password

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
	_, _, ok := handlers.ParseBasicAuth("Invalid")

	if ok {
		t.Error("Expected parsing to fail")
	}
}

func TestParseBasicAuthNoBasicPrefix(t *testing.T) {
	_, _, ok := handlers.ParseBasicAuth("Bearer token")

	if ok {
		t.Error("Expected parsing to fail")
	}
}
