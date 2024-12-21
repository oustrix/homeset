package store

import (
	"context"

	"github.com/oustrix/homeset/internal/models"
	"github.com/oustrix/homeset/internal/store/dto"
)

// Storage describes application data source.
type Storage interface {
	GetUser(ctx context.Context, input dto.GetUserInput) (models.User, error)
	CreateUser(ctx context.Context, input dto.CreateUserInput) (models.User, error)

	Close(ctx context.Context)
}
