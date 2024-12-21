package users_test

import (
	"context"
	"errors"
	"math/rand/v2"
	"testing"

	"github.com/guregu/null/zero"

	"github.com/gojuno/minimock/v3"
	"github.com/google/uuid"
	"github.com/oustrix/homeset/internal/domain/users"
	"github.com/oustrix/homeset/internal/domain/users/mocks"
	"github.com/oustrix/homeset/internal/models"
	"github.com/oustrix/homeset/internal/store"
	"github.com/oustrix/homeset/internal/store/dto"
	"github.com/stretchr/testify/suite"
)

type getUserSuite struct {
	suite.Suite
}

func TestGetUser(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(getUserSuite))
}

func (s *getUserSuite) TestUC_OK() {
	ctrl := minimock.NewController(s.T())

	user := models.User{
		ID:           rand.Int64(),
		Username:     uuid.NewString(),
		PasswordHash: uuid.NewString(),
	}

	storage := mocks.NewGetUserRepositoryMock(ctrl)
	storage.GetUserMock.
		Expect(minimock.AnyContext, dto.GetUserInput{
			IDEq:       zero.IntFrom(user.ID),
			UsernameEq: zero.StringFrom(user.Username),
		}).
		Return(user, nil)

	handle := users.NewGetUser(users.GetUserConfig{
		Storage: storage,
	})

	result, err := handle(context.Background(), users.GetUserParams{
		ID:       zero.IntFrom(user.ID),
		Username: zero.StringFrom(user.Username),
	})
	s.Require().NoError(err)
	s.Require().NotEmpty(result)

	s.Require().EqualValues(user, result.User)
}

func (s *getUserSuite) TestUC_Error_NotFound() {
	ctrl := minimock.NewController(s.T())

	id := zero.IntFrom(rand.Int64())

	storage := mocks.NewGetUserRepositoryMock(ctrl)
	storage.GetUserMock.
		Expect(minimock.AnyContext, dto.GetUserInput{
			IDEq: id,
		}).
		Return(models.User{}, store.ErrNotFound)

	handle := users.NewGetUser(users.GetUserConfig{
		Storage: storage,
	})

	result, err := handle(context.Background(), users.GetUserParams{
		ID: id,
	})
	s.Require().ErrorIs(err, users.ErrUserNotFound)
	s.Require().Empty(result)
}

func (s *getUserSuite) TestUC_Error_Unexpected() {
	ctrl := minimock.NewController(s.T())

	storage := mocks.NewGetUserRepositoryMock(ctrl)
	storage.GetUserMock.
		Expect(minimock.AnyContext, dto.GetUserInput{}).
		Return(models.User{}, errors.New("unexpected error"))

	handle := users.NewGetUser(users.GetUserConfig{
		Storage: storage,
	})

	result, err := handle(context.Background(), users.GetUserParams{})
	s.Require().EqualError(err, "storage.GetUser: unexpected error")
	s.Require().Empty(result)
}
