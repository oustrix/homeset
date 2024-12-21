package users

import (
	"context"
	"errors"
	"fmt"

	"github.com/guregu/null/zero"
	"github.com/oustrix/homeset/internal/domain"
	"github.com/oustrix/homeset/internal/models"
	"github.com/oustrix/homeset/internal/store"
	"github.com/oustrix/homeset/internal/store/dto"
)

var (
	ErrUserNotFound = domain.Error{
		Description: "user not found",
	}
)

//go:generate go run github.com/gojuno/minimock/v3/cmd/minimock -g -s _mock.go -i GetUserRepository -o mocks

// GetUserParams used to provide data for GetUser.
type GetUserParams struct {
	// User ID.
	ID zero.Int
	// User name.
	Username zero.String
}

// GetUserResult used to output data from GetUser.
type GetUserResult struct {
	// Found user.
	User models.User
}

// GetUser used to get a user by filter.
// Possible errors:
// - ErrUserNotFound: cannot find user with given filters.
type GetUser func(ctx context.Context, params GetUserParams) (GetUserResult, error)

// GetUserRepository describes data source for GetUser usecase.
type GetUserRepository interface {
	GetUser(ctx context.Context, input dto.GetUserInput) (models.User, error)
}

type getUser struct {
	storage GetUserRepository
}

// GetUserConfig used to provide data for NewGetUser.
type GetUserConfig struct {
	Storage GetUserRepository
}

// NewGetUser creates a new usecase GetUser.
func NewGetUser(config GetUserConfig) GetUser {
	uc := getUser{
		storage: config.Storage,
	}

	return uc.handle
}

func (uc *getUser) handle(ctx context.Context, params GetUserParams) (GetUserResult, error) {
	input := uc.buildGetUserInput(params)

	user, err := uc.storage.GetUser(ctx, input)
	switch {
	case err == nil:
		// OK
	case errors.Is(err, store.ErrNotFound):
		return GetUserResult{}, ErrUserNotFound
	default:
		return GetUserResult{}, fmt.Errorf("storage.GetUser: %w", err)
	}

	return uc.buildResult(user), nil
}

func (uc *getUser) buildGetUserInput(params GetUserParams) dto.GetUserInput {
	return dto.GetUserInput{
		IDEq:       params.ID,
		UsernameEq: params.Username,
	}
}

func (uc *getUser) buildResult(user models.User) GetUserResult {
	return GetUserResult{
		User: user,
	}
}
