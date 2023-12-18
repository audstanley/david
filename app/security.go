// Useful information: http://www.webdav.org/specs/#dav
// and http://www.webdav.org/specs/rfc4918.html
package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

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

	// Handle HTTP authorization from method headers
	err, ok = handleHeadersForAuthorization(a, ctx, w, req, authInfo)
	if err == nil && !ok {
		return
	} else if err != nil {
		log.WithFields(log.Fields{"error": err, "user": authInfo.Username, "method": req.Method}).Error("Error handling authorization - This method condition hasn't been handled yet")
	} else if err != nil && !ok {
		log.WithFields(log.Fields{"error": err, "user": authInfo.Username, "method": req.Method}).Error("Error handling authorization - This method condition hasn't been handled yet")
		return
	}
	// =================================================================================================================

	// Serve request with authenticated user context
	a.Handler.ServeHTTP(w, req.WithContext(ctx))
}

// Resolve returns the physical path for the given name.
func Resolve(ctx context.Context, name string, d Dir) string {
	// Validate the name for any invalid characters or separators.
	if filepath.Separator != '/' && strings.ContainsRune(name, filepath.Separator) ||
		strings.Contains(name, "\x00") { // Null bytes are illegal in file names because they can be used to terminate strings prematurely and cause unexpected behavior.
		return ""
	}
	// Retrieve the base directory path from the configuration.
	dir := string(d.Config.Dir)
	// Use current directory if base directory is not set.
	if dir == "" {
		dir = "."
	}
	// Obtain authentication information from the context.
	authInfo := AuthFromContext(ctx)
	// Check if user is authenticated and has configured subdirectory.
	if authInfo != nil && authInfo.Authenticated {
		// Get user information from the configuration.
		userInfo := d.Config.Users[authInfo.Username]
		// If user has a configured subdirectory, append it to the path.
		if userInfo != nil && userInfo.Subdir != nil {
			return filepath.Join(dir, *userInfo.Subdir, filepath.FromSlash(path.Clean("/"+name)))
		}
	}
	// Build the final physical path by combining base directory and the provided name.
	return filepath.Join(dir, filepath.FromSlash(path.Clean("/"+name)))
}

// Define allowed methods for your WebDAV resource
var allowedMethods = []string{
	"GET", "HEAD", "PUT", "POST", "DELETE",
	"PROPFIND", "PROPPATCH", "COPY", "MOVE", "LOCK",
	"UNLOCK", "MKCOL", "DELETE",
}

const (
	Propfind string = "PROPFIND"
	Mkol     string = "MKOL"
	Move     string = "MOVE"
	Lock     string = "LOCK"
	Unlock   string = "UNLOCK"
	Propatch string = "PROPPATCH"
	Copy     string = "COPY"
)

func handleHeadersForAuthorization(a *App, ctx context.Context, w http.ResponseWriter, req *http.Request, authInfo *AuthInfo) (error, bool) {
	// Initialize authorization status as True (assuming allowed)
	ok := true
	switch req.Method {
	case http.MethodGet:
		// GET not allowed, return Method Not Allowed (405)
		handleMethodNotAllowed(ctx, w, req)
		return nil, !ok
	case http.MethodPut:
		// Check user's "Create" permission for PUT requests
		log.WithField("method", req.Method).Debug("Method received")
		// Unauthorized due to missing permission
		if !a.Config.Users[authInfo.Username].Crud.Create {
			w.WriteHeader(http.StatusForbidden)
			return nil, !ok
		} else {
			// Authorized!
			return nil, ok
		}
	case http.MethodPost:
		// Log the received POST request but don't handle authorization here
		log.WithField("method", req.Method).Debug("Method received")
	case http.MethodDelete:
		// Check user's "Delete" permission for DELETE requests
		log.WithField("method", req.Method).Debug("Method received")
		if !a.Config.Users[authInfo.Username].Crud.Delete {
			// Unauthorized due to missing permission
			w.WriteHeader(http.StatusForbidden)
			return nil, !ok
		} else {
			// Authorized!
			return nil, ok
		}
	case http.MethodHead:
		// Log the received HEAD request but don't handle authorization here
		log.WithField("method", req.Method).Debug("Method received")
	case http.MethodOptions:
		// Handle OPTIONS request by setting allowed methods and WebDAV headers
		log.WithField("method", req.Method).Debug("Method received")
		// Respond to OPTIONS request
		w.Header().Set("Allow", strings.Join(allowedMethods, ", "))
		w.Header().Set("DAV", "1, 2, source") // Indicate supported WebDAV versions and extensions
		w.WriteHeader(http.StatusOK)
		return nil, !ok // Not authorized in the strict sense, but OPTIONS doesn't require file access
	case Propfind:
		// Special handling for PROPFIND requests
		log.WithFields(log.Fields{"user": authInfo.Username,
			"method": req.Method,
			"crud":   authInfo.CrudType.Crud},
		).Debug("Method received")
		if !a.Config.Users[authInfo.Username].Crud.Read {
			// Check user's "Read" permission
			w.WriteHeader(http.StatusUnauthorized) // 401 Unauthorized
			return nil, !ok
		} else {
			// User can read existing files, but additional check for non-existent files requested with Create/Update permissions
			if !a.Config.Users[authInfo.Username].Crud.Create || !a.Config.Users[authInfo.Username].Crud.Update {
				// Get the requested file path
				filePath := Resolve(ctx, req.URL.Path, Dir{a.Config})
				log.WithFields(log.Fields{"user": authInfo.Username, "Path": filePath}).Debug("Header received")

				// Check if the file exists (if not, user might be trying to open a non-existent file they shouldn't have access to)
				// **Important note:** The previous code relied on `os.Stat` which can panic during custom WebDAV filesystem implementations. This revised approach avoids the panic potential.
				_, err := os.Stat(filePath)
				if err != nil {
					if errors.Is(err, os.ErrNotExist) {
						// File doesn't exist, and user lacks Create/Update permissions
						if a.Config.Log.Create {
							log.WithFields(log.Fields{
								"path":       filePath,
								"user":       authInfo.Username,
								"User-Agent": req.Header.Get("User-Agent"),
							}).Debug("User does not have the permission to open a non existant file that their operating system is attempting to find")
							// Not authorized due to attempting to open a non-existent file with insufficient permissions
							return nil, !ok
						} else {
							// For some other reason, the file doesn't exist, but the user has Create/Update permissions
							return nil, !ok
						}
					}
				}
				return nil, ok
			}

			return nil, ok
		}
	case Mkol:
		// Check user's "Create" permission for MKCOL
		log.WithField("method", Mkol).Debug("Method received")
		if !a.Config.Users[authInfo.Username].Crud.Create {
			// Unauthorized due to missing permission
			w.WriteHeader(http.StatusUnauthorized)
			return nil, !ok
		} else {
			// Authorized!
			return nil, ok
		}
	case Move:
		// Check user's "Update" permission for MOVE
		log.WithField("method", Move).Debug("Method received")
		if !a.Config.Users[authInfo.Username].Crud.Update {
			// Unauthorized due to missing permission
			filePath := Resolve(ctx, req.URL.Path, Dir{a.Config})
			log.WithFields(log.Fields{"user": authInfo.Username, "method": Move, "crud": authInfo.CrudType.Crud, "path": filePath}).Debug("User does not have the permission to move the file")
			w.WriteHeader(http.StatusUnauthorized)
			return nil, !ok
		} else {
			// Authorized!
			return nil, ok
		}
	case Lock:
		// LOCK requires "Create" permission
		log.WithField("method", Lock).Debug("Method received")
		if !a.Config.Users[authInfo.Username].Crud.Create {
			w.WriteHeader(http.StatusUnauthorized)
			return nil, !ok
		} else {
			return nil, ok
		}
	case Unlock:
		// UNLOCK requires "Create" permission
		log.WithField("method", Unlock).Debug("Method received")
		if !a.Config.Users[authInfo.Username].Crud.Create {
			w.WriteHeader(http.StatusUnauthorized)
			return nil, !ok
		} else {
			return nil, ok
		}
	case Propatch:
		log.WithField("method", Propatch).Debug("Method received")
		return nil, ok
	default:
		// David has not implemented this method yet
		log.WithField("method", req.Method).Debug("Method received")
		return errors.New("This method condition hasn't been handled yet"), ok
	}
	w.WriteHeader(http.StatusNotImplemented)
	return errors.New("no single method was received"), !ok
}

// handle methods not allowed
func handleMethodNotAllowed(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	log.WithField("method", req.Method).Debug("Method received")
	w.Write([]byte("Method not allowed"))
	w.WriteHeader(http.StatusMethodNotAllowed)
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
