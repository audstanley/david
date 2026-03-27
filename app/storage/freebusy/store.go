package freebusy

import (
	"encoding/json"
	"fmt"

	"github.com/audstanley/david/app/storage"
)

const (
	keyPrefixFreeBusy = "freebusy:"
)

// Store implements FreeBusyStore
type Store struct {
	db *storage.Storage
}

// New creates a new FreeBusy store
func New(db *storage.Storage) *Store {
	return &Store{db: db}
}

// Store stores a free/busy block
func (s *Store) Store(fb *storage.FreeBusyBlock) error {
	data, err := storage.Encode(fb)
	if err != nil {
		return fmt.Errorf("failed to encode free/busy block: %w", err)
	}

	pairs := [][2][]byte{
		{[]byte(keyPrefixFreeBusy + fb.UID), data},
	}

	return s.db.Batch(pairs)
}

// Get retrieves a free/busy block by UID
func (s *Store) Get(uid string) (*storage.FreeBusyBlock, error) {
	data, err := s.db.Get([]byte(keyPrefixFreeBusy + uid))
	if err != nil {
		return nil, err
	}

	var fb storage.FreeBusyBlock
	if err := storage.Decode(data, &fb); err != nil {
		return nil, fmt.Errorf("failed to decode free/busy block: %w", err)
	}

	return &fb, nil
}

// List retrieves all free/busy blocks
func (s *Store) List() ([]*storage.FreeBusyBlock, error) {
	var blocks []*storage.FreeBusyBlock

	err := s.db.Iterate([]byte(keyPrefixFreeBusy), func(key, value []byte) error {
		var fb storage.FreeBusyBlock
		if err := storage.Decode(value, &fb); err != nil {
			return err
		}
		blocks = append(blocks, &fb)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return blocks, nil
}

// Delete deletes a free/busy block by UID
func (s *Store) Delete(uid string) error {
	_, err := s.Get(uid)
	if err != nil {
		return err
	}

	return s.db.Delete([]byte(keyPrefixFreeBusy + uid))
}

// Encode serializes a slice of free/busy blocks
func encodeFreeBusyBlocks(blocks []storage.FreeBusyBlock) (string, error) {
	data, err := json.Marshal(blocks)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// DecodeFreeBusyBlocks deserializes a slice of free/busy blocks
func decodeFreeBusyBlocks(s string) ([]storage.FreeBusyBlock, error) {
	var blocks []storage.FreeBusyBlock
	if err := json.Unmarshal([]byte(s), &blocks); err != nil {
		return nil, err
	}
	return blocks, nil
}
