package app

import (
	"context"
	"errors"
	"time"

	"github.com/FreakyGranny/otus/hw12_13_14_15_calendar/internal/storage"
)

var (
	// ErrEventFieldWrong required params are wrong.
	ErrEventFieldWrong = errors.New("some fields are wrong")
	// ErrEventIDZero id must be non zero value.
	ErrEventIDZero = errors.New("id must be positive")
	// ErrDateIsNotMonday date is not start of week.
	ErrDateIsNotMonday = errors.New("date is not monday")
	// ErrDateIsNotFirstDay date is not start of month.
	ErrDateIsNotFirstDay = errors.New("date is not first day of month")
)

// Application business logic.
type Application interface {
	GetEvent(ctx context.Context, id int64) (*storage.Event, error)
	GetEventForDay(ctx context.Context, date time.Time) ([]*storage.Event, error)
	GetEventForWeek(ctx context.Context, date time.Time) ([]*storage.Event, error)
	GetEventForMonth(ctx context.Context, date time.Time) ([]*storage.Event, error)
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

// GetEventForDay returns list of events for day.
func (a *App) GetEventForDay(ctx context.Context, date time.Time) ([]*storage.Event, error) {
	return a.storage.GetEventList(ctx, date.Truncate(24*time.Hour), time.Hour*24)
}

// GetEventForWeek returns list of events for week.
func (a *App) GetEventForWeek(ctx context.Context, date time.Time) ([]*storage.Event, error) {
	if date.Weekday() != time.Monday {
		return nil, ErrDateIsNotMonday
	}

	return a.storage.GetEventList(ctx, date.Truncate(24*time.Hour), time.Hour*24*7)
}

// GetEventForMonth returns list of events for month.
func (a *App) GetEventForMonth(ctx context.Context, date time.Time) ([]*storage.Event, error) {
	if date.Day() != 1 {
		return nil, ErrDateIsNotFirstDay
	}
	dayCount := 30
	if date.Month() == 2 {
		dayCount = 28
	}
	t := time.Date(date.Year(), date.Month(), dayCount, 0, 0, 0, 0, time.UTC).AddDate(0, 1, 0)
	if t.Day() != 1 {
		dayCount++
	}

	return a.storage.GetEventList(ctx, date.Truncate(24*time.Hour), time.Hour*time.Duration(24*dayCount))
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
