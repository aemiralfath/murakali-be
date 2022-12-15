package postgre

import (
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"murakali/config"
)

func NewPG(config *config.Config) (*sql.DB, error) {
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		config.Postgres.PostgresqlUser,
		config.Postgres.PostgresqlPassword,
		config.Postgres.PostgresqlHost,
		config.Postgres.PostgresqlDbname,
	)
	fmt.Println(connString)

	connectionCfg, err := pgx.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config %w", err)
	}

	connStr := stdlib.RegisterConnConfig(connectionCfg)
	db, err := sql.Open(config.Postgres.PgDriver, connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return db, err
	}

	return db, nil
}
