package sqlstorage

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/FreakyGranny/otus/hw12_13_14_15_calendar/internal/storage"
	_ "github.com/jackc/pgx/stdlib" //nolint
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func sslModeToString(sslEnable bool) string {
	if sslEnable {
		return "enable"
	}

	return "disable"
}

// BuildDsn build dsn string from params.
func BuildDsn(host string, port int, user string, password string, dbName string, sslMode bool) string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		host,
		strconv.Itoa(port),
		user,
		dbName,
		password,
		sslModeToString(sslMode),
	)
}

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
func (s *Storage) CreateEvent(e *storage.Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

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
func (s *Storage) GetEvent(id int64) (*storage.Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	e := &storage.Event{}
	err := s.db.QueryRowxContext(ctx, "SELECT id, owner_id, title, descr, start_date, end_date, notify_before FROM events WHERE id = $1", id).StructScan(e)

	return e, err
}

// GetEventList returns list of events.
func (s *Storage) GetEventList() ([]*storage.Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res := make([]*storage.Event, 0)
	rows, err := s.db.QueryxContext(ctx, "SELECT id, owner_id, title, descr, start_date, end_date, notify_before FROM events")
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
	log.Info().Msgf("%d", len(res))

	return res, err
}

// UpdateEvent updates event.
func (s *Storage) UpdateEvent(e *storage.Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	query := `UPDATE events SET owner_id=:owner_id, title=:title, descr=:descr, start_date=:start_date, end_date=:end_date, notify_before=:notify_before WHERE id=:id;`
	_, err := s.db.NamedExecContext(ctx, query, e)

	return err
}

// DeleteEvent deletes event.
func (s *Storage) DeleteEvent(id int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := s.db.ExecContext(ctx, `DELETE FROM events WHERE id = $1`, id)

	return err
}
