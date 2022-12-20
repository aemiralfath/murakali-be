package postgre

import (
	"context"
	"database/sql"

	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"murakali/config"
	"murakali/pkg/logger"
)

type myQueryTracer struct {
	log *zap.SugaredLogger
}

func (tracer *myQueryTracer) TraceQueryStart(ctx context.Context, pg *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	tracer.log.Infow("Executing command start", "sql", data.SQL, "args", data.Args)
	return ctx
}

func (tracer *myQueryTracer) TraceQueryEnd(_ context.Context, _ *pgx.Conn, data pgx.TraceQueryEndData) {
	tracer.log.Infow("Executing command end", "sql err", data.Err, "args tag", data.CommandTag)
}

func NewPG(cfg *config.Config, log logger.Logger) (*sql.DB, error) {
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",

		cfg.Postgres.PostgresqlUser,

		cfg.Postgres.PostgresqlPassword,

		cfg.Postgres.PostgresqlHost,

		cfg.Postgres.PostgresqlDbname,
	)

	connectionCfg, err := pgx.ParseConfig(connString)

	if err != nil {
		return nil, fmt.Errorf("failed to parse cfg %w", err)
	}

	connectionCfg.Tracer = &myQueryTracer{
		log: log.GetZapLogger(),
	}

	connStr := stdlib.RegisterConnConfig(connectionCfg)

	db, err := sql.Open(cfg.Postgres.PgDriver, connStr)

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return db, err
	}

	return db, nil
}
