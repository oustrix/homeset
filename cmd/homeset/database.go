package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/oustrix/homeset/internal/store"
	sqliteStore "github.com/oustrix/homeset/internal/store/sqlite"
	"github.com/oustrix/homeset/pkg/sqlite"
	"github.com/pressly/goose/v3"

	_ "github.com/oustrix/homeset/migrations"
)

type storageConfig struct {
	DBMSName string

	sqliteConfig sqliteConfig
}

type sqliteConfig struct {
	DSN string
}

func resolveStorage(ctx context.Context, config storageConfig) (store.Storage, error) {
	switch config.DBMSName {
	case "sqlite":
		return createSQLiteStorage(config.sqliteConfig)
	}

	slog.ErrorContext(
		ctx,
		"failed to resolve storage",
		"dbms", config.DBMSName,
	)

	return nil, errors.New("failed to resolve storage")
}

func createSQLiteStorage(config sqliteConfig) (*sqliteStore.Storage, error) {
	const op = "createSQLiteStorage"

	sqliteDB, err := sqlite.New(sqlite.Config{
		DSN: config.DSN,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: sqlite.New: %w", op, err)
	}

	err = goose.SetDialect(string(goose.DialectSQLite3))
	if err != nil {
		return nil, fmt.Errorf("%s: goose.SetDialect: %w", op, err)
	}

	err = migrate(sqliteDB.DB)
	if err != nil {
		return nil, fmt.Errorf("%s: migrate: %w", op, err)
	}

	sqliteStorage := sqliteStore.NewStorage(sqliteDB)

	return sqliteStorage, nil
}

func migrate(db *sql.DB) error {
	return goose.Up(db, "migrations")
}
