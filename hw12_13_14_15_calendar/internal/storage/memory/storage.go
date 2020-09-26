package memorystorage

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/FreakyGranny/otus/hw12_13_14_15_calendar/internal/storage"
)

var (
	errEventNotFound   = errors.New("event not found")
	errContextCanceled = errors.New("context canceled")
)

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
func (s *Storage) CreateEvent(ctx context.Context, e *storage.Event) error {
	select {
	case <-ctx.Done():
		return errContextCanceled
	default:
		s.mu.Lock()
		defer s.mu.Unlock()
		s.seq++
		e.ID = s.seq
		s.values[e.ID] = *e

		return nil
	}
}

// GetEvent returns event by id.
func (s *Storage) GetEvent(ctx context.Context, id int64) (*storage.Event, error) {
	select {
	case <-ctx.Done():
		return nil, errContextCanceled
	default:
		s.mu.RLock()
		defer s.mu.RUnlock()
		e, ok := s.values[id]
		if !ok {
			return nil, errEventNotFound
		}

		return &e, nil
	}
}

// GetEventList returns list of events.
func (s *Storage) GetEventList(ctx context.Context, date time.Time, period time.Duration) ([]*storage.Event, error) {
	result := make([]*storage.Event, 0)
	select {
	case <-ctx.Done():
		return result, errContextCanceled
	default:
		s.mu.RLock()
		defer s.mu.RUnlock()
		for _, e := range s.values {
			x := e
			if e.StartDate.After(date) && e.StartDate.Before(date.Add(period)) {
				result = append(result, &x)
			}
		}

		return result, nil
	}
}

// UpdateEvent updates event.
func (s *Storage) UpdateEvent(ctx context.Context, e *storage.Event) error {
	select {
	case <-ctx.Done():
		return errContextCanceled
	default:
		s.mu.Lock()
		defer s.mu.Unlock()
		_, ok := s.values[e.ID]
		if !ok {
			return errEventNotFound
		}
		s.values[e.ID] = *e

		return nil
	}
}

// DeleteEvent deletes event.
func (s *Storage) DeleteEvent(ctx context.Context, id int64) error {
	select {
	case <-ctx.Done():
		return errContextCanceled
	default:
		s.mu.Lock()
		defer s.mu.Unlock()
		delete(s.values, id)

		return nil
	}
}

// DeleteOldEvents deletes events older then given delta.
func (s *Storage) DeleteOldEvents(ctx context.Context, d int) error {
	return nil
}

// GetEventForNotification returns list of events.
func (s *Storage) GetEventForNotification(ctx context.Context, i time.Duration) ([]*storage.Event, error) {
	return make([]*storage.Event, 0), nil
}
