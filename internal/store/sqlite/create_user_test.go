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

type createUserIntegrationSuite struct {
	testsqlite.Suite

	storage *sqlite.Storage
}

func TestCreateUserIntegration(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(createUserIntegrationSuite))
}

func (s *createUserIntegrationSuite) SetupTest() {
	s.Suite.SetupTest()

	s.storage = sqlite.NewStorage(s.SQLite)
}

func (s *createUserIntegrationSuite) Test_OK() {
	ctx := context.Background()

	var (
		username     = uuid.NewString()
		passwordHash = uuid.NewString()
	)

	createdUser, err := s.storage.CreateUser(ctx, dto.CreateUserInput{
		Username:     username,
		PasswordHash: passwordHash,
	})
	s.Require().NoError(err)
	s.Require().NotEmpty(createdUser)

	foundUser, err := s.storage.GetUser(ctx, dto.GetUserInput{
		IDEq: zero.IntFrom(createdUser.ID),
	})
	s.Require().NoError(err)
	s.Require().NotEmpty(foundUser)

	s.Require().Equal(username, foundUser.Username)
}

func (s *createUserIntegrationSuite) Test_Error_NotUniqueUsername() {
	ctx := context.Background()

	firstUser, err := s.storage.CreateUser(ctx, dto.CreateUserInput{
		Username: uuid.NewString(),
	})
	s.Require().NoError(err)
	s.Require().NotEmpty(firstUser)

	secondUser, err := s.storage.CreateUser(ctx, dto.CreateUserInput{
		Username: firstUser.Username,
	})
	s.Require().ErrorIs(err, store.ErrUniqueViolation)
	s.Require().Empty(secondUser)
}
