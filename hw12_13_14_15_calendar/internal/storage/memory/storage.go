package memorystorage

import (
	"errors"
	"sync"

	"github.com/FreakyGranny/otus/hw12_13_14_15_calendar/internal/storage"
)

var errEventNotFound = errors.New("event not found")

// Storage memory storage.
type Storage struct {
	mu     sync.RWMutex
	seq    int64
	values map[int64]storage.Event
}

// New returns new memory storage.
func New() *Storage {
	return &Storage{
		values: make(map[int64]storage.Event),
	}
}

// Close closes storage.
func (s *Storage) Close() error {
	return nil
}

// CreateEvent creates new event.
func (s *Storage) CreateEvent(e *storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.seq++
	e.ID = s.seq
	s.values[e.ID] = *e

	return nil
}

// GetEvent returns event by id.
func (s *Storage) GetEvent(id int64) (*storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.values[id]
	if !ok {
		return nil, errEventNotFound
	}

	return &e, nil
}

// GetEventList returns list of events.
func (s *Storage) GetEventList() ([]*storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]*storage.Event, 0)
	for _, e := range s.values {
		x := e
		result = append(result, &x)
	}

	return result, nil
}

// UpdateEvent updates event.
func (s *Storage) UpdateEvent(e *storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.values[e.ID]
	if !ok {
		return errEventNotFound
	}
	s.values[e.ID] = *e

	return nil
}

// DeleteEvent deletes event.
func (s *Storage) DeleteEvent(id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.values, id)

	return nil
}
