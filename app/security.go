package app

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type contextKey int

var authInfoKey contextKey

// AuthInfo holds the username and authentication status
type AuthInfo struct {
	Username      string
	Authenticated bool
	CrudType      *CrudType
}

// authWebdavHandlerFunc is a type definition which holds a context and application reference to
// match the AuthWebdavHandler interface.
type authWebdavHandlerFunc func(c context.Context, w http.ResponseWriter, r *http.Request, a *App)

// ServeHTTP simply calls the AuthWebdavHandlerFunc with given parameters
func (f authWebdavHandlerFunc) ServeHTTP(c context.Context, w http.ResponseWriter, r *http.Request, a *App) {
	f(c, w, r, a)
}

// NewBasicAuthWebdavHandler creates a new http handler with basic auth features.
// The handler will use the application config for user and password lookups.
func NewBasicAuthWebdavHandler(a *App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		handlerFunc := authWebdavHandlerFunc(handle)
		handlerFunc.ServeHTTP(ctx, w, r, a)
	})
}

var testCrudType = CrudType{"", false, false, false, false}

// authenticate validates the provided username and password against the configured users and returns an AuthInfo object.
func authenticate(cfg *Config, username, password string) (*AuthInfo, error) {

	// Perform authentication only if required
	if !cfg.AuthenticationNeeded() {
		return &AuthInfo{Username: "", Authenticated: false, CrudType: &testCrudType}, nil
	}

	// Validate username and password presence
	if username == "" || password == "" {
		return &AuthInfo{Username: username, Authenticated: false, CrudType: &testCrudType}, errors.New("username not found or password empty")
	}

	// Retrieve user information from configuration
	user := cfg.Users[username]
	crud := cfg.Users[username].Crud

	if user == nil {
		return &AuthInfo{Username: username, Authenticated: false, CrudType: &testCrudType}, errors.New("user not found")
	}
	// Verify provided password against stored hash
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return &AuthInfo{Username: username, Authenticated: false, CrudType: &testCrudType}, errors.New("Password doesn't match")
	}

	// Return successful authentication information
	return &AuthInfo{Username: username, Authenticated: true, CrudType: crud}, nil
}

// AuthFromContext returns information about the authentication state of the current user.
func AuthFromContext(ctx context.Context) *AuthInfo {
	// Attempt to retrieve the AuthInfo object from the context
	info, ok := ctx.Value(authInfoKey).(*AuthInfo)
	if !ok {
		// Return nil if AuthInfo is not present
		return nil
	}
	// Return the retrieved AuthInfo object
	return info
}

func handle(ctx context.Context, w http.ResponseWriter, req *http.Request, a *App) {
	// CORS preflight request handling
	if req.Method == "OPTIONS" {
		// Allow preflight requests from configured origins and with valid headers
		if a.Config.Cors.Origin == req.Header.Get("Origin") &&
			req.Header.Get("Access-Control-Request-Method") != "" &&
			req.Header.Get("Access-Control-Request-Headers") != "" {
			w.WriteHeader(http.StatusNoContent) // Send empty response
			return
		}
	}

	// Authentication bypass for systems without users
	if !a.Config.AuthenticationNeeded() {
		a.Handler.ServeHTTP(w, req.WithContext(ctx))
		return
	}
	// Extract username and password from HTTP Basic Auth header
	username, password, ok := httpAuth(req, a.Config)
	if !ok {
		// Respond with Unauthorized status and optional realm
		SayUnauthorized(w, a.Config.Realm)
		return
	}
	// Authenticate user credentials
	authInfo, err := authenticate(a.Config, username, password)
	// Log failed login attempt with user and IP address
	if err != nil {
		ipAddr := req.Header.Get("X-Forwarded-For")
		if len(ipAddr) == 0 {
			remoteAddr := req.RemoteAddr
			lastIndex := strings.LastIndex(remoteAddr, ":")
			if lastIndex != -1 {
				ipAddr = remoteAddr[:lastIndex]
			} else {
				ipAddr = remoteAddr
			}
		}
		log.WithField("user", username).WithField("address", ipAddr).WithError(err).Warn("User failed to login")
	}
	// Check if user is authenticated and authorized
	if !authInfo.Authenticated || !authInfo.CrudType.Read {
		// Respond with Unauthorized status and optional realm
		SayUnauthorized(w, a.Config.Realm)
		return
	}
	// Add authentication information to context
	ctx = context.WithValue(ctx, authInfoKey, authInfo)
	// Serve request with authenticated user context
	a.Handler.ServeHTTP(w, req.WithContext(ctx))
}

func httpAuth(r *http.Request, config *Config) (string, string, bool) {
	if config.AuthenticationNeeded() {
		username, password, ok := r.BasicAuth()
		return username, password, ok
	}

	return "", "", true
}

func SayUnauthorized(w http.ResponseWriter, realm string) {
	w.Header().Set("WWW-Authenticate", "Basic realm="+realm)
	w.WriteHeader(http.StatusUnauthorized)
	_, err := w.Write([]byte(fmt.Sprintf("%d %s", http.StatusUnauthorized, "Unauthorized")))

	if err != nil {
		log.WithError(err).Error("Error sending unauthorized response")
	}
}

// GenHash generates a bcrypt hashed password string
func GenHash(password []byte) string {
	pw, err := bcrypt.GenerateFromPassword(password, 10)
	if err != nil {
		log.Fatal(err)
	}

	return string(pw)
}
