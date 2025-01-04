package sqlite_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/guregu/null/zero"
	"github.com/oustrix/homeset/internal/store"
	"github.com/oustrix/homeset/internal/store/dto"
	"github.com/oustrix/homeset/internal/store/sqlite"
	"github.com/oustrix/homeset/test/testsqlite"
	"github.com/stretchr/testify/suite"
)

type getUserIntegrationSuite struct {
	testsqlite.Suite

	storage *sqlite.Storage
}

func TestGetUserIntegration(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(getUserIntegrationSuite))
}

func (s *getUserIntegrationSuite) SetupTest() {
	s.Suite.SetupTest()

	s.storage = sqlite.NewStorage(s.SQLite)
}

func (s *getUserIntegrationSuite) TestSQLite_OK() {
	ctx := context.Background()

	createdUser, err := s.storage.CreateUser(ctx, dto.CreateUserInput{
		Username:     uuid.NewString(),
		PasswordHash: uuid.NewString(),
	})
	s.Require().NoError(err)
	s.Require().NotEmpty(createdUser)

	foundUser, err := s.storage.GetUser(ctx, dto.GetUserInput{
		UsernameEq: zero.StringFrom(createdUser.Username),
	})
	s.Require().NoError(err)
	s.Require().NotEmpty(foundUser)

	s.Require().EqualValues(createdUser, foundUser)
}

func (s *getUserIntegrationSuite) TestSQLite_OK_ByUsername() {
	ctx := context.Background()

	createdUser, err := s.storage.CreateUser(ctx, dto.CreateUserInput{
		Username:     uuid.NewString(),
		PasswordHash: uuid.NewString(),
	})
	s.Require().NoError(err)
	s.Require().NotEmpty(createdUser)

	foundUser, err := s.storage.GetUser(ctx, dto.GetUserInput{
		UsernameEq: zero.StringFrom(createdUser.Username),
	})
	s.Require().NoError(err)
	s.Require().NotEmpty(foundUser)

	s.Require().EqualValues(createdUser, foundUser)
}

func (s *getUserIntegrationSuite) TestSQLite_Error_NotFound() {
	user, err := s.storage.GetUser(context.Background(), dto.GetUserInput{
		UsernameEq: zero.StringFrom(uuid.NewString()),
	})
	s.Require().ErrorIs(err, store.ErrNotFound)
	s.Require().Empty(user)
}
