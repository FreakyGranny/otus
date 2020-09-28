package sqlstorage

import (
	"context"
	"time"

	"github.com/FreakyGranny/otus/hw12_13_14_15_calendar/internal/storage"
	_ "github.com/jackc/pgx/stdlib" //nolint
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

// Storage sql storage implementation.
type Storage struct {
	db *sqlx.DB
}

// New returns new sql storage.
func New(dsn string) *Storage {
	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("failed to load driver")
	}

	return &Storage{db: db}
}

// Close close db connection.
func (s *Storage) Close() error {
	return s.db.Close()
}

// CreateEvent creates new event.
func (s *Storage) CreateEvent(ctx context.Context, e *storage.Event) error {
	query := `INSERT INTO events(owner_id, title, descr, start_date, end_date, notify_before)
						  VALUES(:owner_id, :title, :descr, :start_date, :end_date, :notify_before) RETURNING id`
	rows, err := s.db.NamedQueryContext(ctx, query, e)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&e.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetEvent returns event by id.
func (s *Storage) GetEvent(ctx context.Context, id int64) (*storage.Event, error) {
	e := &storage.Event{}
	err := s.db.QueryRowxContext(ctx, "SELECT id, owner_id, title, descr, start_date, end_date, notify_before FROM events WHERE id = $1", id).StructScan(e)

	return e, err
}

// GetEventList returns list of events.
func (s *Storage) GetEventList(ctx context.Context, date time.Time, period time.Duration) ([]*storage.Event, error) {
	res := make([]*storage.Event, 0)
	query := `SELECT id, owner_id, title, descr, start_date, end_date, notify_before FROM events
					WHERE start_date
					BETWEEN $1 AND $1 + make_interval(hours => $2)
	`
	rows, err := s.db.QueryxContext(ctx, query, date, int(period.Hours()))
	if err != nil {
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		var e storage.Event
		if err := rows.StructScan(&e); err != nil {
			return res, err
		}
		res = append(res, &e)
	}

	return res, err
}

// UpdateEvent updates event.
func (s *Storage) UpdateEvent(ctx context.Context, e *storage.Event) error {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Error().
			Err(err).
			Msg("unable to create transaction")
	}
	dbEvent := &storage.Event{}
	err = tx.QueryRowxContext(ctx, "SELECT id, owner_id, title, descr, start_date, end_date, notify_before FROM events WHERE id = $1", e.ID).StructScan(dbEvent)
	if err != nil {
		return err
	}

	if len(e.Title) == 0 {
		e.Title = dbEvent.Title
	}
	if e.StartDate.Before(time.Unix(0, 1)) {
		e.StartDate = dbEvent.StartDate
	}
	if e.EndDate.Before(time.Unix(0, 1)) {
		e.EndDate = dbEvent.EndDate
	}
	if e.OwnerID == 0 {
		e.OwnerID = dbEvent.OwnerID
	}
	if len(e.Descr) == 0 {
		e.Descr = dbEvent.Descr
	}
	if e.NotifyBefore == 0 {
		e.NotifyBefore = dbEvent.NotifyBefore
	}

	query := `UPDATE events SET owner_id=:owner_id, title=:title, descr=:descr, start_date=:start_date, end_date=:end_date, notify_before=:notify_before WHERE id=:id;`
	_, err = tx.NamedExecContext(ctx, query, e)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// DeleteEvent deletes event.
func (s *Storage) DeleteEvent(ctx context.Context, id int64) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM events WHERE id = $1`, id)

	return err
}

// DeleteOldEvents deletes events older then given delta.
func (s *Storage) DeleteOldEvents(ctx context.Context, d int) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM events WHERE start_date + make_interval(days => $1) < now()`, d)

	return err
}

// GetEventForNotification returns list of events.
func (s *Storage) GetEventForNotification(ctx context.Context, i time.Duration) ([]*storage.Event, error) {
	res := make([]*storage.Event, 0)
	query := `SELECT id, owner_id, title, descr, start_date, end_date, notify_before FROM events
					WHERE notify_before > 0 
					AND (start_date - make_interval(secs => notify_before)) 
					BETWEEN now() AND now() + make_interval(secs => $1)`
	rows, err := s.db.QueryxContext(ctx, query, i.Seconds())
	if err != nil {
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		var e storage.Event
		if err := rows.StructScan(&e); err != nil {
			return res, err
		}
		res = append(res, &e)
	}

	return res, err
}
