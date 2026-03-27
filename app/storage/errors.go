// Package storage provides LevelDB-based storage for iCalendar data
package storage

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/syndtr/goleveldb/leveldb"
	leveldbErrors "github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/util"
)

// StorageError represents a storage error
type StorageError struct {
	Code    string
	Message string
	Err     error
}

func (e *StorageError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *StorageError) Unwrap() error {
	return e.Err
}

// Error codes
var (
	ErrNotFound      = &StorageError{Code: "NOT_FOUND", Message: "entity not found"}
	ErrAlreadyExists = &StorageError{Code: "ALREADY_EXISTS", Message: "entity already exists"}
	ErrInvalidInput  = &StorageError{Code: "INVALID_INPUT", Message: "invalid input data"}
	ErrDatabaseOpen  = &StorageError{Code: "DATABASE_OPEN", Message: "failed to open database"}
	ErrDatabaseClose = &StorageError{Code: "DATABASE_CLOSE", Message: "failed to close database"}
	ErrEncode        = &StorageError{Code: "ENCODE", Message: "encoding error"}
	ErrDecode        = &StorageError{Code: "DECODE", Message: "decoding error"}
	ErrPermission    = &StorageError{Code: "PERMISSION", Message: "permission denied"}
	ErrConflict      = &StorageError{Code: "CONFLICT", Message: "conflict detected"}
	ErrRateLimit     = &StorageError{Code: "RATE_LIMIT", Message: "rate limit exceeded"}
	ErrLimitReached  = &StorageError{Code: "LIMIT_REACHED", Message: "limit reached"}
)

// IsNotFound checks if an error is a not found error
func IsNotFound(err error) bool {
	var storageErr *StorageError
	if errors.As(err, &storageErr) {
		return storageErr.Code == "NOT_FOUND"
	}
	return errors.Is(err, leveldbErrors.ErrNotFound)
}

// IsAlreadyExists checks if an error is an already exists error
func IsAlreadyExists(err error) bool {
	var storageErr *StorageError
	if errors.As(err, &storageErr) {
		return storageErr.Code == "ALREADY_EXISTS"
	}
	return false
}

// Storage provides the base storage operations
type Storage struct {
	db *leveldb.DB
}

// NewStorage creates a new storage instance
func NewStorage(path string) (*Storage, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, &StorageError{Code: "DATABASE_OPEN", Message: "failed to open database", Err: err}
	}
	return &Storage{db: db}, nil
}

// Close closes the database
func (s *Storage) Close() error {
	if err := s.db.Close(); err != nil {
		return &StorageError{Code: "DATABASE_CLOSE", Message: "failed to close database", Err: err}
	}
	return nil
}

// Encode serializes data to JSON
func Encode(v interface{}) ([]byte, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, &StorageError{Code: "ENCODE", Message: "encoding error", Err: err}
	}
	return data, nil
}

// Decode deserializes JSON to data
func Decode(data []byte, v interface{}) error {
	if err := json.Unmarshal(data, v); err != nil {
		return &StorageError{Code: "DECODE", Message: "decoding error", Err: err}
	}
	return nil
}

// Get retrieves a value by key
func (s *Storage) Get(key []byte) ([]byte, error) {
	data, err := s.db.Get(key, nil)
	if err != nil {
		if errors.Is(err, leveldb.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, &StorageError{Code: "GET", Message: "failed to get value", Err: err}
	}
	return data, nil
}

// Put stores a value by key
func (s *Storage) Put(key, value []byte) error {
	err := s.db.Put(key, value, nil)
	if err != nil {
		return &StorageError{Code: "PUT", Message: "failed to put value", Err: err}
	}
	return nil
}

// Delete removes a value by key
func (s *Storage) Delete(key []byte) error {
	err := s.db.Delete(key, nil)
	if err != nil {
		return &StorageError{Code: "DELETE", Message: "failed to delete value", Err: err}
	}
	return nil
}

// Exists checks if a key exists
func (s *Storage) Exists(key []byte) (bool, error) {
	_, err := s.db.Get(key, nil)
	if err != nil {
		if errors.Is(err, leveldb.ErrNotFound) {
			return false, nil
		}
		return false, &StorageError{Code: "EXISTS", Message: "failed to check existence", Err: err}
	}
	return true, nil
}

// Iterate iterates over keys with a prefix
func (s *Storage) Iterate(prefix []byte, fn func(key, value []byte) error) error {
	iter := s.db.NewIterator(util.BytesPrefix(prefix), nil)
	defer iter.Release()

	for iter.Next() {
		key := iter.Key()
		value := iter.Value()
		if err := fn(key, value); err != nil {
			return err
		}
	}

	if err := iter.Error(); err != nil {
		return &StorageError{Code: "ITERATE", Message: "iteration error", Err: err}
	}

	return nil
}

// Batch puts multiple key-value pairs atomically
func (s *Storage) Batch(pairs [][2][]byte) error {
	batch := &leveldb.Batch{}
	for _, pair := range pairs {
		if pair[1] == nil {
			batch.Delete(pair[0])
		} else {
			batch.Put(pair[0], pair[1])
		}
	}
	err := s.db.Write(batch, nil)
	if err != nil {
		return &StorageError{Code: "BATCH", Message: "failed to write batch", Err: err}
	}
	return nil
}

// CompareAndSwap atomically compares and swaps a value
func (s *Storage) CompareAndSwap(key, expected, new []byte) (bool, error) {
	batch := &leveldb.Batch{}

	current, err := s.db.Get(key, nil)
	if err != nil && !errors.Is(err, leveldbErrors.ErrNotFound) {
		return false, &StorageError{Code: "CAS", Message: "failed to get current value", Err: err}
	}

	if !bytes.Equal(current, expected) {
		return false, nil
	}

	batch.Put(key, new)
	if err := s.db.Write(batch, nil); err != nil {
		return false, &StorageError{Code: "CAS", Message: "failed to write batch", Err: err}
	}

	return true, nil
}
