package sqlite

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/oustrix/homeset/internal/models"
	"github.com/oustrix/homeset/internal/store"
	"github.com/oustrix/homeset/internal/store/dto"
)

// GetUser finds user with given filters.
func (s *Storage) GetUser(ctx context.Context, input dto.GetUserInput) (models.User, error) {
	query := s.Builder.
		Select(store.UsersTableColumns).
		From(store.UsersTable)

	if input.IDEq.Valid {
		query = query.Where(sq.Eq{"id": input.IDEq.Int64})
	}

	if input.UsernameEq.Valid {
		query = query.Where(sq.Eq{"username": input.UsernameEq.String})
	}

	return getx[models.User](ctx, s.DB, query)
}
