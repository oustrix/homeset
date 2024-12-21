package testsqlite

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	_ "github.com/oustrix/homeset/migrations"
	"github.com/oustrix/homeset/pkg/sqlite"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
	*sqlite.SQLite
}

func (s *Suite) SetupTest() {
	sqliteDB, err := sqlite.New(sqlite.Config{
		DSN: fmt.Sprintf("file:%s?mode=memory", uuid.NewString()),
	})
	s.Require().NoError(err)
	s.Require().NotEmpty(sqliteDB)

	provider, err := goose.NewProvider(
		goose.DialectSQLite3,
		sqliteDB.DB,
		nil,
	)
	s.Require().NoError(err)
	s.Require().NotNil(provider)

	_, err = provider.Up(context.Background())
	s.Require().NoError(err)

	s.SQLite = sqliteDB
}

func (s *Suite) TearDownTest() {
	s.Require().NoError(s.DB.Close())
}
