package sqlite

import (
	"context"
	"fmt"

	"github.com/oustrix/homeset/internal/models"
	"github.com/oustrix/homeset/internal/store"
	"github.com/oustrix/homeset/internal/store/dto"
)

// CreateUser creates new user.
func (s *Storage) CreateUser(ctx context.Context, input dto.CreateUserInput) (models.User, error) {
	fields := map[string]interface{}{
		"username":      input.Username,
		"password_hash": input.PasswordHash,
	}

	query := s.Builder.
		Insert(store.UsersTable).
		SetMap(fields).
		Suffix(fmt.Sprintf("RETURNING %s", store.UsersTableColumns))

	return getx[models.User](ctx, s.DB, query)
}
