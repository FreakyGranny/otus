package storage

import (
	"context"
	"time"
)

//go:generate mockgen -source=$GOFILE -destination=../mocks/storage_mock.go -package=mocks .

// Storage storage provider.
type Storage interface {
	GetEvent(ctx context.Context, id int64) (*Event, error)
	GetEventList(ctx context.Context) ([]*Event, error)
	CreateEvent(ctx context.Context, e *Event) error
	UpdateEvent(ctx context.Context, e *Event) error
	DeleteEvent(ctx context.Context, id int64) error
	Close() error
	DeleteOldEvents(ctx context.Context, d int) error
	GetEventForNotification(ctx context.Context, i time.Duration) ([]*Event, error)
}

// Event calendar event.
type Event struct {
	ID           int64     `db:"id" json:"id"`
	Title        string    `db:"title" json:"title"`
	StartDate    time.Time `db:"start_date" json:"start_date"`
	EndDate      time.Time `db:"end_date" json:"end_date"`
	Descr        string    `db:"descr" json:"descr"`
	OwnerID      int64     `db:"owner_id" json:"owner_id"`
	NotifyBefore int64     `db:"notify_before" json:"notify_before"`
}
