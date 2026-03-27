package todos

import (
	"fmt"
	"time"

	"github.com/audstanley/david/app/storage"
)

const (
	keyPrefixTodo      = "todo:"
	keyPrefixTodoByUID = "uid:"
	keyPrefixTodoByCal = "cal:"
)

// Store implements TodoStore
type Store struct {
	db *storage.Storage
}

// New creates a new todo store
func New(db *storage.Storage) *Store {
	return &Store{db: db}
}

// Create creates a new todo
func (s *Store) Create(todo *storage.Todo) error {
	// Check if todo UID already exists
	exists, err := s.db.Exists([]byte(keyPrefixTodoByUID + todo.UID))
	if err != nil {
		return fmt.Errorf("failed to check todo UID: %w", err)
	}
	if exists {
		return storage.ErrAlreadyExists
	}

	todo.Sequence = 1
	todo.Created = now()
	todo.LastModified = now()

	data, err := storage.Encode(todo)
	if err != nil {
		return fmt.Errorf("failed to encode todo: %w", err)
	}

	pairs := [][2][]byte{
		{[]byte(keyPrefixTodo + todo.UID), data},
		{[]byte(keyPrefixTodoByUID + todo.UID), []byte(todo.CalendarUID)},
		{[]byte(keyPrefixTodoByCal + todo.CalendarUID + ":" + todo.UID), data},
	}

	return s.db.Batch(pairs)
}

// GetByUID retrieves a todo by UID
func (s *Store) GetByUID(todoUID string) (*storage.Todo, error) {
	data, err := s.db.Get([]byte(keyPrefixTodo + todoUID))
	if err != nil {
		return nil, err
	}

	var todo storage.Todo
	if err := storage.Decode(data, &todo); err != nil {
		return nil, fmt.Errorf("failed to decode todo: %w", err)
	}

	return &todo, nil
}

// Update updates a todo
func (s *Store) Update(todoUID string, todo *storage.Todo, sequence int) error {
	existing, err := s.GetByUID(todoUID)
	if err != nil {
		return err
	}

	// Update sequence if needed
	if sequence != 0 && sequence > existing.Sequence {
		todo.Sequence = sequence
	}

	todo.LastModified = now()

	data, err := storage.Encode(todo)
	if err != nil {
		return fmt.Errorf("failed to encode todo: %w", err)
	}

	pairs := [][2][]byte{
		{[]byte(keyPrefixTodo + todoUID), data},
		{[]byte(keyPrefixTodoByCal + todo.CalendarUID + ":" + todoUID), data},
	}

	return s.db.Batch(pairs)
}

// Delete deletes a todo
func (s *Store) Delete(todoUID string) error {
	todo, err := s.GetByUID(todoUID)
	if err != nil {
		return err
	}

	pairs := [][2][]byte{
		{[]byte(keyPrefixTodo + todoUID), nil},
		{[]byte(keyPrefixTodoByUID + todoUID), nil},
		{[]byte(keyPrefixTodoByCal + todo.CalendarUID + ":" + todoUID), nil},
	}

	return s.db.Batch(pairs)
}

// List lists todos for a calendar with pagination
func (s *Store) List(calendarUID string, limit, offset int) ([]*storage.Todo, error) {
	if limit <= 0 {
		limit = 100
	}
	if limit > 1000 {
		limit = 1000
	}

	var todos []*storage.Todo
	count := 0

	err := s.db.Iterate([]byte(keyPrefixTodo), func(key, value []byte) error {
		var todo storage.Todo
		if err := storage.Decode(value, &todo); err != nil {
			return err
		}

		if todo.CalendarUID == calendarUID {
			if count < offset {
				count++
				return nil
			}
			todos = append(todos, &todo)
			count++
			if len(todos) >= limit {
				return storage.ErrLimitReached
			}
		}
		return nil
	})

	if err != nil && err != storage.ErrLimitReached {
		return nil, err
	}

	return todos, nil
}

// now returns the current time
func now() time.Time {
	return time.Now().UTC()
}
