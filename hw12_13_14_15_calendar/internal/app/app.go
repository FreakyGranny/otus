package app

import (
	"context"

	"github.com/FreakyGranny/otus/hw12_13_14_15_calendar/internal/storage"
)

// App calendar instance.
type App struct {
	storage Storage
}

// Storage sql storage.
type Storage interface {
	GetEvent(id int64) (*storage.Event, error)
	GetEventList() ([]*storage.Event, error)
	CreateEvent(e *storage.Event) error
	UpdateEvent(e *storage.Event) error
	DeleteEvent(id int64) error
	Close() error
}

// New returns new app.
func New(storage Storage) *App {
	return &App{
		storage: storage,
	}
}

// CreateEvent creates new calendar event.
func (a *App) CreateEvent(ctx context.Context, title string) error {
	return a.storage.CreateEvent(
		&storage.Event{
			Title: title,
		},
	)
}
