package sqlite

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/oustrix/homeset/internal/store"
)

func handleSQLiteError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return store.ErrNotFound
	}

	if strings.Contains(err.Error(), "UNIQUE constraint failed") {
		return store.ErrUniqueViolation
	}

	return err
}
