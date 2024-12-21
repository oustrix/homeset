package users

import (
	"context"
	"errors"
	"fmt"

	"github.com/oustrix/homeset/internal/domain"
	"github.com/oustrix/homeset/internal/store"
	"golang.org/x/crypto/bcrypt"

	"github.com/oustrix/homeset/internal/models"
	"github.com/oustrix/homeset/internal/store/dto"
)

var ErrUserAlreadyExists = domain.Error{
	Description: "user already exists",
}

//go:generate go run github.com/gojuno/minimock/v3/cmd/minimock -g -s _mock.go -i CreateUserRepository -o mocks

// CreateUserParams used to provide data for CreateUser.
type CreateUserParams struct {
	// New user name.
	Username string
	// New user password.
	Password string
}

// CreateUserResult used to output data from CreateUser.
type CreateUserResult struct {
	// Created user.
	User models.User
}

// CreateUser creates a new user.
type CreateUser func(ctx context.Context, params CreateUserParams) (CreateUserResult, error)

// CreateUserRepository describes data source for CreateUser.
type CreateUserRepository interface {
	CreateUser(ctx context.Context, input dto.CreateUserInput) (models.User, error)
}

// CreateUserConfig used to provide data for NewCreateUser.
type CreateUserConfig struct {
	Storage CreateUserRepository
}

type createUser struct {
	storage CreateUserRepository
}

// NewCreateUser creates a new CreateUser usecase.
func NewCreateUser(config CreateUserConfig) CreateUser {
	uc := createUser{
		storage: config.Storage,
	}

	return uc.handle
}

func (uc *createUser) handle(ctx context.Context, params CreateUserParams) (CreateUserResult, error) {
	passwordHash, err := uc.hashUserPassword(params.Password)
	if err != nil {
		return CreateUserResult{}, fmt.Errorf("uc.hashUserPassword: %w", err)
	}

	input := uc.buildCreateUserInput(params.Username, passwordHash)

	user, err := uc.storage.CreateUser(ctx, input)
	switch {
	case err == nil:
		// OK
	case errors.Is(err, store.ErrUniqueViolation):
		return CreateUserResult{}, ErrUserAlreadyExists
	default:
		return CreateUserResult{}, fmt.Errorf("storage.CreateUser: %w", err)
	}

	return uc.buildResult(user), nil
}

func (uc *createUser) hashUserPassword(password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("bcrypt.GenerateFromPassword: %w", err)
	}

	return string(passwordHash), nil
}

func (uc *createUser) buildCreateUserInput(username, passwordHash string) dto.CreateUserInput {
	return dto.CreateUserInput{
		Username:     username,
		PasswordHash: passwordHash,
	}
}

func (uc *createUser) buildResult(user models.User) CreateUserResult {
	return CreateUserResult{
		User: user,
	}
}
