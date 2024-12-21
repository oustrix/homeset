package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/Masterminds/squirrel"
	_ "github.com/mattn/go-sqlite3"
)

// SQLite is wrapper around sql.DB sqlite3 driver.
type SQLite struct {
	*sql.DB

	Builder squirrel.StatementBuilderType
}

// Config used to provide data for New.
type Config struct {
	DSN string
}

// New creates a new SQLite.
func New(config Config) (*SQLite, error) {
	sqlite := &SQLite{
		Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}

	sqlDB, err := sql.Open("sqlite3", config.DSN)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	err = sqlDB.Ping()
	if err != nil {
		return nil, fmt.Errorf("sqlDB.Ping: %w", err)
	}

	sqlite.DB = sqlDB

	return sqlite, nil
}

func (s *SQLite) Close(ctx context.Context) {
	err := s.DB.Close()
	if err != nil {
		slog.ErrorContext(
			ctx,
			"failed to close sqlite3 connection",
			"error", err.Error(),
		)
	}
}
