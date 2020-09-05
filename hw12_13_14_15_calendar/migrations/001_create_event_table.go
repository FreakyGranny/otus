package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(Up001, Down001)
}

// Up001 up migration.
func Up001(tx *sql.Tx) error {
	_, err := tx.Exec(
		`CREATE TABLE events (
			id bigserial NOT NULL primary key,
			title varchar NOT NULL,
			start_date timestamp NOT NULL,
			end_date timestamp NOT NULL,
			descr varchar,
			owner_id int NOT NULL,
			notify_before NOT NULL DEFAULT 0);
	`)
	if err != nil {
		return err
	}

	return nil
}

// Down001 down migration.
func Down001(tx *sql.Tx) error {
	_, err := tx.Exec("DROP TABLE events;")
	if err != nil {
		return err
	}

	return nil
}
