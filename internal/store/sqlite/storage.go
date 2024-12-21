package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/sqlscan"
	sqliteDriver "github.com/oustrix/homeset/pkg/sqlite"
)

// Storage is a implementation of store.Storage with sqlite3 database.
type Storage struct {
	*sqliteDriver.SQLite
}

// NewStorage creates a new Storage.
func NewStorage(sqliteDB *sqliteDriver.SQLite) *Storage {
	return &Storage{sqliteDB}
}

// Close closes database connection.
func (s *Storage) Close(ctx context.Context) {
	s.SQLite.Close(ctx)
}

func getx[T any](ctx context.Context, db *sql.DB, query sq.Sqlizer) (T, error) {
	dest := new(T)

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return *dest, fmt.Errorf("query.ToSql: %w", err)
	}

	err = sqlscan.Get(ctx, db, dest, sqlQuery, args...)
	if err != nil {
		return *dest, handleSQLiteError(fmt.Errorf("sqlscan.Get: %w", err))
	}

	return *dest, nil
}
