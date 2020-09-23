package queue

import (
	"time"
)

// Message Queue unit.
type Message struct {
	ContentType string
	Body        []byte
}

// Notification notification for calendar event.
type Notification struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	StartDate time.Time `json:"date"`
	OwnerID   int64     `json:"user"`
}
