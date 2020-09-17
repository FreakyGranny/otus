package app

import (
	"context"
	"errors"
	"time"

	"github.com/FreakyGranny/otus/hw12_13_14_15_calendar/internal/storage"
)

// ErrEventFieldWrong required params are wrong.
var ErrEventFieldWrong = errors.New("some fields are wrong")

// ErrEventIDZero id must be non zero value.
var ErrEventIDZero = errors.New("id must be positive")

// Application business logic.
type Application interface {
	GetEvent(ctx context.Context, id int64) (*storage.Event, error)
	GetEventList(ctx context.Context) ([]*storage.Event, error)
	CreateEvent(ctx context.Context, title string, startDate time.Time, endDate time.Time, ownerID int64, descr string, notifyBefore int64) (*storage.Event, error)
	UpdateEvent(ctx context.Context, id int64, title string, startDate time.Time, endDate time.Time, ownerID int64, descr string, notifyBefore int64) (*storage.Event, error)
	DeleteEvent(ctx context.Context, id int64) error
}

// App calendar instance.
type App struct {
	storage storage.Storage
}

// New returns new app.
func New(storage storage.Storage) *App {
	return &App{
		storage: storage,
	}
}

// GetEvent returns event by id.
func (a *App) GetEvent(ctx context.Context, id int64) (*storage.Event, error) {
	if id <= 0 {
		return nil, ErrEventIDZero
	}

	return a.storage.GetEvent(ctx, id)
}

// GetEventList returns list of events.
func (a *App) GetEventList(ctx context.Context) ([]*storage.Event, error) {
	return a.storage.GetEventList(ctx)
}

// CreateEvent creates new calendar event.
func (a *App) CreateEvent(ctx context.Context, title string, startDate time.Time, endDate time.Time, ownerID int64, descr string, notifyBefore int64) (*storage.Event, error) {
	if len(title) == 0 || ownerID <= 0 || startDate.After(endDate) {
		return nil, ErrEventFieldWrong
	}

	e := &storage.Event{
		Title:        title,
		StartDate:    startDate,
		EndDate:      endDate,
		OwnerID:      ownerID,
		Descr:        descr,
		NotifyBefore: notifyBefore,
	}

	err := a.storage.CreateEvent(ctx, e)
	if err != nil {
		return nil, err
	}

	return e, nil
}

// UpdateEvent updates event.
func (a *App) UpdateEvent(ctx context.Context, id int64, title string, startDate time.Time, endDate time.Time, ownerID int64, descr string, notifyBefore int64) (*storage.Event, error) {
	e := &storage.Event{
		ID:           id,
		Title:        title,
		StartDate:    startDate,
		EndDate:      endDate,
		OwnerID:      ownerID,
		Descr:        descr,
		NotifyBefore: notifyBefore,
	}
	err := a.storage.UpdateEvent(ctx, e)
	if err != nil {
		return nil, err
	}

	return e, nil
}

// DeleteEvent deletes event.
func (a *App) DeleteEvent(ctx context.Context, id int64) error {
	if id <= 0 {
		return ErrEventIDZero
	}

	return a.storage.DeleteEvent(ctx, id)
}
