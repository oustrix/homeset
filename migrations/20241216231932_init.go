package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upInit, downInit)
}

func upInit(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(
		ctx,
		`
		CREATE TABLE IF NOT EXISTS users (
            username      TEXT    PRIMARY KEY,
            password_hash TEXT    NOT NULL
         );
		`,
	)
	if err != nil {
		return err
	}

	return nil
}

func downInit(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(
		ctx,
		`
		DROP TABLE IF EXISTS users;
		`,
	)
	if err != nil {
		return err
	}

	return nil
}
