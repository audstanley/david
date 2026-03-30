package users

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/audstanley/david/app"
	"github.com/audstanley/david/app/storage"
	"golang.org/x/crypto/bcrypt"
)

const (
	keyPrefixUser       = "user:"
	keyPrefixUserByUID  = "uid:"
	keyPrefixUserByName = "name:"
)

// Store implements UserStore
type Store struct {
	db *storage.Storage
}

// New creates a new user store
func New(db *storage.Storage) *Store {
	return &Store{db: db}
}

// Create creates a new user
func (s *Store) Create(userID, username, passwordHash, email, displayName string, role string) error {
	// Check if username already exists
	exists, err := s.db.Exists([]byte(keyPrefixUserByName + username))
	if err != nil {
		return fmt.Errorf("failed to check username: %w", err)
	}
	if exists {
		return storage.ErrAlreadyExists
	}

	user := &storage.User{
		ID:           userID,
		Username:     username,
		PasswordHash: passwordHash,
		Email:        email,
		DisplayName:  displayName,
		Role:         role,
		APIKeys:      []string{},
		Created:      now(),
		Updated:      now(),
	}

	data, err := storage.Encode(user)
	if err != nil {
		return fmt.Errorf("failed to encode user: %w", err)
	}

	pairs := [][2][]byte{
		{[]byte(keyPrefixUser + userID), data},
		{[]byte(keyPrefixUserByUID + userID), []byte(username)},
		{[]byte(keyPrefixUserByName + username), []byte(userID)},
	}

	return s.db.Batch(pairs)
}

// GetByID retrieves a user by ID
func (s *Store) GetByID(userID string) (*storage.User, error) {
	data, err := s.db.Get([]byte(keyPrefixUser + userID))
	if err != nil {
		return nil, err
	}

	var user storage.User
	if err := storage.Decode(data, &user); err != nil {
		return nil, fmt.Errorf("failed to decode user: %w", err)
	}

	return &user, nil
}

// GetByUsername retrieves a user by username
func (s *Store) GetByUsername(username string) (*storage.User, error) {
	userIDData, err := s.db.Get([]byte(keyPrefixUserByName + username))
	if err != nil {
		return nil, err
	}

	return s.GetByID(string(userIDData))
}

// Update updates a user
func (s *Store) Update(userID string, updates map[string]interface{}) error {
	user, err := s.GetByID(userID)
	if err != nil {
		return err
	}

	// Store old username before applying updates
	oldName := user.Username

	// Apply updates
	if v, ok := updates["username"]; ok {
		user.Username = v.(string)
	}
	if v, ok := updates["email"]; ok {
		user.Email = v.(string)
	}
	if v, ok := updates["displayName"]; ok {
		user.DisplayName = v.(string)
	}
	if v, ok := updates["role"]; ok {
		user.Role = v.(string)
	}
	if v, ok := updates["passwordHash"]; ok {
		user.PasswordHash = v.(string)
	}

	user.Updated = now()

	data, err := storage.Encode(user)
	if err != nil {
		return fmt.Errorf("failed to encode user: %w", err)
	}

	pairs := [][2][]byte{
		{[]byte(keyPrefixUser + userID), data},
	}

	if user.Username != oldName {
		// Remove old name index and create new one
		pairs = append(pairs, [][2][]byte{
			{[]byte(keyPrefixUserByName + oldName), nil},
			{[]byte(keyPrefixUserByName + user.Username), []byte(userID)},
		}...)
	} else if _, ok := updates["username"]; ok {
		// Username was updated but is the same value, update the index anyway
		pairs = append(pairs, [][2][]byte{
			{[]byte(keyPrefixUserByName + user.Username), []byte(userID)},
		}...)
	}

	return s.db.Batch(pairs)
}

// Delete deletes a user
func (s *Store) Delete(userID string) error {
	user, err := s.GetByID(userID)
	if err != nil {
		return err
	}

	pairs := [][2][]byte{
		{[]byte(keyPrefixUser + userID), nil},
		{[]byte(keyPrefixUserByUID + userID), nil},
		{[]byte(keyPrefixUserByName + user.Username), nil},
	}

	return s.db.Batch(pairs)
}

// List lists users with pagination
func (s *Store) List(limit, offset int) ([]*storage.User, error) {
	if limit <= 0 {
		limit = 100
	}
	if limit > 1000 {
		limit = 1000
	}

	var users []*storage.User
	skipped := 0
	returned := 0

	err := s.db.Iterate([]byte(keyPrefixUser), func(key, value []byte) error {
		if skipped < offset {
			skipped++
			return nil
		}

		if returned >= limit {
			return storage.ErrLimitReached
		}

		var user storage.User
		if err := storage.Decode(value, &user); err != nil {
			return err
		}
		users = append(users, &user)
		returned++
		return nil
	})

	if err == storage.ErrLimitReached {
		return users, nil
	}
	if err != nil {
		return nil, err
	}

	return users, nil
}

// VerifyPassword verifies a password against the hash
func (s *Store) VerifyPassword(passwordHash, password string) (bool, error) {
	// Direct comparison using bcrypt.CompareHashAndPassword
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		return false, err
	}
	return true, nil
}

// VerifyPasswordHash verifies a password using the app package
func (s *Store) VerifyPasswordHash(password, hash string) error {
	return app.VerifyPasswordHash(password, hash)
}

// HashPassword hashes a password
func (s *Store) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// now returns the current time
func now() time.Time {
	return time.Now().UTC()
}

// Encode serializes a slice of strings
func encodeStringSlice(items []string) (string, error) {
	data, err := json.Marshal(items)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// DecodeStringSlice deserializes a string slice
func decodeStringSlice(s string) ([]string, error) {
	var items []string
	if err := json.Unmarshal([]byte(s), &items); err != nil {
		return nil, err
	}
	return items, nil
}
