package storage

import "time"

// Event calendar event.
type Event struct {
	ID           int64     `db:"id"`
	Title        string    `db:"title"`
	StartDate    time.Time `db:"start_date"`
	EndDate      time.Time `db:"end_date"`
	Descr        string    `db:"descr"`
	OwnerID      int64     `db:"owner_id"`
	NotifyBefore int64     `db:"notify_before"`
}
