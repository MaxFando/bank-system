package postgres_test

import (
	"context"
	"github.com/MaxFando/bank-system/pkg/sqlext/internal/postgres"
	"github.com/MaxFando/bank-system/pkg/sqlext/tests/containers"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type DBSuite struct {
	suite.Suite
	ctx         context.Context
	db          *postgres.DB
	pgContainer *containers.PostgresContainer
}

func (suite *DBSuite) SetupSuite() {
	suite.ctx = context.Background()

	pgContainer, err := containers.CreatePostgresContainer(suite.ctx)
	suite.NoError(err)
	suite.pgContainer = pgContainer

	db, err := sqlx.ConnectContext(suite.ctx, "postgres", pgContainer.ConnectionString)
	suite.NoError(err)
	suite.db = postgres.NewDB(db)

	err = db.PingContext(suite.ctx)
	suite.NoError(err)

	_, err = db.ExecContext(suite.ctx, "CREATE TABLE users (id serial PRIMARY KEY, name VARCHAR(50), age INT);")
	suite.NoError(err)

	_, err = db.ExecContext(suite.ctx, "INSERT INTO users (name, age) VALUES ('Alice', 20), ('Bob', 25), ('Charlie', 30);")
	suite.NoError(err)
}

func (suite *DBSuite) TearDownSuite() {
	err := suite.pgContainer.Terminate(suite.ctx)
	suite.NoError(err)
}

func (suite *DBSuite) TestGet() {
	var row struct {
		ID   int    `db:"id"`
		Name string `db:"name"`
	}

	query := `select id, name from users where name = :name`
	query, args, err := suite.db.BindNamed(query, map[string]interface{}{"name": "Alice"})
	suite.NoError(err)

	err = suite.db.Get(suite.ctx, &row, query, args...)
	suite.NoError(err)

	suite.Equal(1, row.ID)
	suite.Equal("Alice", row.Name)
}

func (suite *DBSuite) TestSelect() {
	type row struct {
		ID   int    `db:"id"`
		Name string `db:"name"`
	}

	rows := make([]row, 0, 3)
	query := `select id, name from users`
	err := suite.db.Select(suite.ctx, &rows, query)
	suite.NoError(err)

	suite.Len(rows, 3)
}

func (suite *DBSuite) TestQuery() {
	query := `select id, name from users`
	rows, err := suite.db.Query(suite.ctx, query)
	suite.NoError(err)

	defer rows.Close()

	var count int
	for rows.Next() {
		count++
	}

	suite.Equal(3, count)
}

func (suite *DBSuite) TestCondition() {
	query := `select id, name from users where name = :name and age = :age limit :limit`

	query, args, err := suite.db.BindNamed(query, map[string]interface{}{"name": "Alice", "age": 20, "limit": 1})
	suite.NoError(err)

	rows, err := suite.db.Query(suite.ctx, query, args...)
	suite.NoError(err)

	defer rows.Close()

	var count int
	for rows.Next() {
		count++
	}

	suite.Equal(1, count)
}

func (suite *DBSuite) TestWithTxRollback() {
	err := suite.db.WithTx(suite.ctx, func(ctx context.Context) error {
		query := `insert into users (name) values (:name)`
		_, _ = suite.db.Exec(ctx, query, map[string]interface{}{"name": "Eve"})
		return assert.AnError
	})
	suite.ErrorIs(err, assert.AnError)

	var count int
	err = suite.db.Get(suite.ctx, &count, "select count(*) from users")
	suite.NoError(err)
	suite.Equal(3, count)
}

func TestDBSuite(t *testing.T) {
	suite.Run(t, new(DBSuite))
}
