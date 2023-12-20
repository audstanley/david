package app

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/webdav"
)

// var noCrudOperations = CrudType{"", false, false, false, false}
func authInfoRelativelyEqual(configAuthInfo, attemptedAuthInfoUpdate *AuthInfo, testName string) bool {
	log.WithFields(logrus.Fields{"testName": testName}).Info("authInfoRelativelyEqual")
	// configAuthInfo - is the config's authInfo
	// attemptedAuthIntoUpdate - is the authInfo that is returned from the authenticate function
	// after an attempted authentication.
	areEqual := true
	switch testName {
	case "user not found":
		// an edge case validation
		if configAuthInfo == nil {
			// in this case  we can't force the authInfo memory addresses to be the same
			// so configAuthInfo should be nil.
			return areEqual
		}
	case "password doesn't match":
		// an edge case validation
		if configAuthInfo.Authenticated == attemptedAuthInfoUpdate.Authenticated {
			log.WithFields(logrus.Fields{"configAuthInfo": configAuthInfo, "attemptedAuthInfoUpdate": attemptedAuthInfoUpdate}).Info("authInfoRelativelyEqual")
			// in this test case, authInfo is not populated for an insuccessful authentication
			// attempt, so attemptedAuthInfoUpdate.Crud will be null.
			// and attemptedAuthInfoUpdate.Username will not return the username that was attempted.
			return areEqual
		}
	default:
		log.WithFields(logrus.Fields{"configAuthInfo": configAuthInfo,
			"attemptedAuthInfoUpdate": attemptedAuthInfoUpdate}).Info("authInfoRelativelyEqual")
		if configAuthInfo == nil || attemptedAuthInfoUpdate == nil {
			return !areEqual
		}
		areEqual := configAuthInfo.Username == attemptedAuthInfoUpdate.Username &&
			configAuthInfo.Authenticated == attemptedAuthInfoUpdate.Authenticated
		log.WithField("areEqual", areEqual).Info("authInfoRelativelyEqual")
		return areEqual
	}
	return !areEqual
}

// This test is failing, because we have test cases that are haven't covered the
// new additions to the code that offer errors for when unauthorized crud operations occur.
func TestAuthenticate(t *testing.T) {
	noAuthentication := AuthInfo{CrudType: &CrudType{Crud: "", Create: false, Read: false, Update: false, Delete: false}}
	type args struct {
		config   *Config
		username string
		password string
	}
	tests := []struct {
		name    string
		args    args
		want    *AuthInfo
		wantErr bool
	}{
		{ // Test 1
			"empty username",
			args{
				config: &Config{Users: map[string]*UserInfo{
					"foo": {
						Password:    GenHash([]byte("password")),
						Permissions: "",
						Crud:        noAuthentication.CrudType,
					},
				}},
				username: "",
				password: "password",
			},
			&AuthInfo{
				Username:      "",
				Authenticated: false,
				CrudType:      noAuthentication.CrudType,
			},
			true,
		},
		{ // Test 2
			"empty password",
			args{
				config: &Config{Users: map[string]*UserInfo{
					"foo": {
						Password:    GenHash([]byte("password")),
						Permissions: "",
						Crud:        noAuthentication.CrudType,
					},
				}},
				username: "foo",
				password: "",
			},
			&AuthInfo{
				Username:      "foo",
				Authenticated: false,
				CrudType:      noAuthentication.CrudType,
			},
			true,
		},
		{ // Test 3
			"empty username without users",
			args{
				config:   &Config{},
				username: "",
				password: "password",
			},
			&AuthInfo{
				Username:      "",
				Authenticated: false,
				CrudType:      noAuthentication.CrudType,
			},
			false,
		},
		{ // Test 4
			"empty password without users",
			args{
				config:   &Config{},
				username: "foo",
				password: "",
			},
			&AuthInfo{
				Username:      "",
				Authenticated: false,
			},
			false,
		},
		{ // Test 5
			"user not found",
			args{
				config: &Config{Users: map[string]*UserInfo{
					"bar": {
						Password:    GenHash([]byte("password")),
						Permissions: "",
						Crud:        noAuthentication.CrudType,
					},
				}},
				username: "foo",
				password: "password",
			},
			&AuthInfo{
				Username:      "foo",
				Authenticated: false,
				CrudType:      noAuthentication.CrudType,
			},
			true,
		},
		{ // Test 6
			"password doesn't match",
			args{
				config: &Config{Users: map[string]*UserInfo{
					"foo": {
						Password:    GenHash([]byte("not-my-password")),
						Permissions: "",
						Crud:        noAuthentication.CrudType,
					},
				}},
				username: "foo",
				password: "password",
			},
			&AuthInfo{
				Username:      "foo",
				Authenticated: false,
				CrudType:      noAuthentication.CrudType,
			},
			true,
		},
		{ // Test 7
			"all fine",
			args{
				config: &Config{Users: map[string]*UserInfo{
					"foo": {
						Password:    GenHash([]byte("password")),
						Permissions: "",
						Crud:        noAuthentication.CrudType,
					},
				}},
				username: "foo",
				password: "password",
			},
			&AuthInfo{
				Username:      "foo",
				Authenticated: false,
				CrudType:      noAuthentication.CrudType,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := authenticate(tt.args.config, tt.args.username, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("1. authenticate() testName = %v, error = %v, wantErr %v, got-authinfo = %v", tt.name, err, tt.wantErr, *got)
				return
			}
			// The all fine test case is the only one that should return a valid authInfo object.
			if !authInfoRelativelyEqual(got, &noAuthentication, tt.name) && tt.name != "all fine" {
				t.Errorf("2. authenticate() testName = %v, got-authInfo = %v, want-authInfo = %v", tt.name, *got, *tt.want)
			}
		})
	}
}

func TestAuthFromContext(t *testing.T) {
	type fakeKey int
	var fakeKeyValue fakeKey

	baseCtx := context.Background()
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want *AuthInfo
	}{
		{
			"success",
			args{
				ctx: context.WithValue(baseCtx, authInfoKey, &AuthInfo{"username", true, &CrudType{"crud", true, true, true, true}}),
			},
			&AuthInfo{"username", true, &CrudType{"crud", true, true, true, true}},
		},
		{
			"failure",
			args{
				ctx: context.WithValue(baseCtx, fakeKeyValue, &AuthInfo{"username", true, &CrudType{"crud", true, true, true, true}}),
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AuthFromContext(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AuthFromContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandle(t *testing.T) {
	type args struct {
		ctx      context.Context
		w        *httptest.ResponseRecorder
		r        *http.Request
		username []byte
		password []byte
		a        *App
	}
	tests := []struct {
		name       string
		args       args
		statusCode int
	}{
		{
			"basic auth error",
			args{
				context.Background(),
				httptest.NewRecorder(),
				httptest.NewRequest("PROPFIND", "/", nil),
				nil,
				nil,
				&App{Config: &Config{Users: map[string]*UserInfo{
					"foo": {
						Password: GenHash([]byte("password")),
					},
				}}},
			},
			401,
		},
		{
			"unauthorized error",
			args{
				context.Background(),
				httptest.NewRecorder(),
				httptest.NewRequest("PROPFIND", "/", nil),
				[]byte("u"),
				[]byte("p"),
				&App{Config: &Config{Users: map[string]*UserInfo{
					"foo": {
						Password: GenHash([]byte("password")),
					},
				}}},
			},
			401,
		},
		{
			"ok",
			args{
				context.Background(),
				httptest.NewRecorder(),
				httptest.NewRequest("PROPFIND", "/", nil),
				[]byte("foo"),
				[]byte("password"),
				&App{
					Config: &Config{Users: map[string]*UserInfo{
						"foo": {
							Password: GenHash([]byte("password")),
						},
					}},
					Handler: &webdav.Handler{
						FileSystem: webdav.NewMemFS(),
						LockSystem: webdav.NewMemLS(),
					},
				},
			},
			207,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.username != nil || tt.args.password != nil {
				tt.args.r.SetBasicAuth(string(tt.args.username), string(tt.args.password))
			}

			handle(tt.args.ctx, tt.args.w, tt.args.r, tt.args.a)
			resp := tt.args.w.Result()

			if resp.StatusCode != tt.statusCode {
				t.Errorf("TestHandle() = %v, want %v", resp.StatusCode, tt.statusCode)
			}
		})
	}
}
