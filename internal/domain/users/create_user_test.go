package users_test

import (
	"context"
	"errors"
	"math/rand/v2"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/google/uuid"
	"github.com/oustrix/homeset/internal/domain/users"
	"github.com/oustrix/homeset/internal/domain/users/mocks"
	"github.com/oustrix/homeset/internal/models"
	"github.com/oustrix/homeset/internal/store"
	"github.com/oustrix/homeset/internal/store/dto"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
)

type createUserSuite struct {
	suite.Suite
}

func TestCreateUser(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(createUserSuite))
}

func (s *createUserSuite) TestUC_OK() {
	ctrl := minimock.NewController(s.T())

	var (
		username        = uuid.NewString()
		password        = uuid.NewString()
		passwordHash, _ = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	)

	storage := mocks.NewCreateUserRepositoryMock(ctrl)
	storage.CreateUserMock.
		Set(func(_ context.Context, input dto.CreateUserInput) (u1 models.User, err error) {
			s.Require().Equal(username, input.Username)

			compareErr := bcrypt.CompareHashAndPassword([]byte(input.PasswordHash), []byte(password))
			s.Require().NoError(compareErr)

			return models.User{
				ID:           rand.Int64(),
				Username:     username,
				PasswordHash: string(passwordHash),
			}, nil
		})

	handle := users.NewCreateUser(users.CreateUserConfig{
		Storage: storage,
	})

	result, err := handle(context.Background(), users.CreateUserParams{
		Username: username,
		Password: password,
	})
	s.Require().NoError(err)
	s.Require().NotEmpty(result)

	s.Require().Equal(username, result.User.Username)
	s.Require().NoError(bcrypt.CompareHashAndPassword([]byte(result.User.PasswordHash), []byte(password)))
}

func (s *createUserSuite) TestUC_Err_UserAlreadyExists() {
	ctrl := minimock.NewController(s.T())

	storage := mocks.NewCreateUserRepositoryMock(ctrl)
	storage.CreateUserMock.
		Return(models.User{}, store.ErrUniqueViolation)

	handle := users.NewCreateUser(users.CreateUserConfig{
		Storage: storage,
	})

	result, err := handle(context.Background(), users.CreateUserParams{
		Username: uuid.NewString(),
		Password: uuid.NewString(),
	})
	s.Require().ErrorIs(err, users.ErrUserAlreadyExists)
	s.Require().Empty(result)
}

func (s *createUserSuite) TestUC_Err_Unexpected() {
	ctrl := minimock.NewController(s.T())

	storage := mocks.NewCreateUserRepositoryMock(ctrl)
	storage.CreateUserMock.
		Return(models.User{}, errors.New("unexpected error"))

	handle := users.NewCreateUser(users.CreateUserConfig{
		Storage: storage,
	})

	result, err := handle(context.Background(), users.CreateUserParams{
		Username: uuid.NewString(),
		Password: uuid.NewString(),
	})
	s.Require().EqualError(err, "storage.CreateUser: unexpected error")
	s.Require().Empty(result)
}
