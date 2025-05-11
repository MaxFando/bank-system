package sqlext

import (
	"context"
	"fmt"
	"github.com/MaxFando/bank-system/pkg/sqlext/internal/postgres"
	"time"

	"github.com/jackc/pgx/v5"
	pgxstd "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type PostgresDB struct {
	*postgres.DB
	conn *sqlx.DB
}

func NewPostgresDB(ctx context.Context, dsn string, opts ...ConnOption) (*PostgresDB, error) {
	cfg := &config{
		maxOpenConns: 2,
		maxIdleConns: 2,

		connLifeTime: 60 * time.Minute,
		connIdleTime: 30 * time.Minute,

		tracingEnabled: true,
	}

	for _, opt := range opts {
		if err := opt(cfg); err != nil {
			return nil, fmt.Errorf("не удалось применить параметр подключения: %w", err)
		}
	}

	connConfig, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("ошибка при разборе DSN: %w", err)
	}

	connConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	dsn = pgxstd.RegisterConnConfig(connConfig)
	dbConn, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("подключение к базе данных: %w", err)
	}

	dbConn.SetMaxIdleConns(cfg.maxIdleConns)
	dbConn.SetMaxOpenConns(cfg.maxOpenConns)

	dbConn.SetConnMaxLifetime(cfg.connLifeTime)
	dbConn.SetConnMaxIdleTime(cfg.connIdleTime)

	if err := dbConn.PingContext(ctx); err != nil {
		_ = dbConn.Close()
		return nil, fmt.Errorf("ошибка при пинге базы данных: %w", err)
	}

	return &PostgresDB{
		DB:   postgres.NewDB(dbConn),
		conn: dbConn,
	}, nil
}

func (db *PostgresDB) Close() error {
	return db.conn.Close()
}
