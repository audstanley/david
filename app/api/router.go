package api

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"github.com/audstanley/david/app/api/handlers"
)

// Config holds API configuration
type Config struct {
	Host         string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// Router manages HTTP routes
type Router struct {
	mux  *mux.Router
	cors *cors.Cors
}

// NewRouter creates a new API router
func NewRouter(config Config) *Router {
	r := &Router{
		mux: mux.NewRouter(),
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type", "X-Requested-With"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	r.cors = c

	return r
}

// SetupRoutes configures all routes
func (r *Router) SetupRoutes() {
	// Health check
	r.mux.HandleFunc("/health", r.healthHandler).Methods("GET")
	r.mux.HandleFunc("/status", r.statusHandler).Methods("GET")

	// API routes
	api := r.mux.PathPrefix("/api").Subrouter()

	// Auth routes
	auth := handlers.NewAuthHandler("secret-key", nil)
	api.HandleFunc("/auth/login", auth.LoginHandler).Methods("POST")
	api.HandleFunc("/auth/refresh", auth.RefreshHandler).Methods("POST")
	api.HandleFunc("/auth/logout", auth.LogoutHandler).Methods("POST")
	api.HandleFunc("/auth/verify", auth.VerifyHandler).Methods("GET")
	api.HandleFunc("/auth/api-key", auth.CreateAPIKeyHandler).Methods("POST")
	api.HandleFunc("/auth/api-key", auth.ListAPIKeyHandler).Methods("GET")
	api.HandleFunc("/auth/api-key/{id}", auth.RevokeAPIKeyHandler).Methods("DELETE")

	// Calendar routes
	calHandler := handlers.NewCalendarHandler(nil)
	api.HandleFunc("/calendars", calHandler.ListCalendarsHandler).Methods("GET")
	api.HandleFunc("/calendars/{uid}", calHandler.GetCalendarHandler).Methods("GET")
	api.HandleFunc("/calendars", calHandler.CreateCalendarHandler).Methods("POST")
	api.HandleFunc("/calendars/{uid}", calHandler.UpdateCalendarHandler).Methods("PUT")
	api.HandleFunc("/calendars/{uid}", calHandler.DeleteCalendarHandler).Methods("DELETE")
	api.HandleFunc("/calendars/{uid}/shares", calHandler.ListSharesHandler).Methods("GET")
	api.HandleFunc("/calendars/{uid}/shares", calHandler.GrantShareHandler).Methods("POST")
	api.HandleFunc("/calendars/{uid}/shares/{id}", calHandler.RevokeShareHandler).Methods("DELETE")
	api.HandleFunc("/calendars/{uid}/export.ics", calHandler.ExportCalendarHandler).Methods("GET")
	api.HandleFunc("/calendars/{uid}/stats", calHandler.GetCalendarStatsHandler).Methods("GET")

	// Event routes
	eventHandler := handlers.NewEventHandler(nil)
	api.HandleFunc("/events", eventHandler.ListEventsHandler).Methods("GET")
	api.HandleFunc("/events/{uid}", eventHandler.GetEventHandler).Methods("GET")
	api.HandleFunc("/events", eventHandler.CreateEventHandler).Methods("POST")
	api.HandleFunc("/events/{uid}", eventHandler.UpdateEventHandler).Methods("PUT")
	api.HandleFunc("/events/{uid}", eventHandler.DeleteEventHandler).Methods("DELETE")
	api.HandleFunc("/events/{uid}/export.ics", eventHandler.ExportEventHandler).Methods("GET")
}

func (r *Router) healthHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy"}`))
}

func (r *Router) statusHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"operational","version":"0.1.0"}`))
}

// ServeHTTP implements http.Handler
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.cors.Handler(r.mux).ServeHTTP(w, req)
}
