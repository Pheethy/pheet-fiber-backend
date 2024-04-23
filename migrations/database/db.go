package database

import (
	"context"
	"io"
	"pheet-fiber-backend/config"

	"github.com/Pheethy/psql"
	"github.com/Pheethy/sqlx"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/opentracing/opentracing-go"
	// _util_tracing "pheet-fiber-backend/service/utils/opentracing"
)

const (
	PGX  = "pgx"
	SQLX = "sqlx"
)

func DBConnect(ctx context.Context, cfg config.IDbConfig) (*sqlx.DB, io.Closer) {
	/* init tracing*/
	// tracer, closer := _util_tracing.Init("flavorparser")
	// opentracing.SetGlobalTracer(tracer)

	/* connect */
	psqlClient := getPostgresClient(cfg.Url(), nil)

	db := psqlClient.GetClient()
	db.DB.SetMaxOpenConns(cfg.MaxConns())

	return db, nil
}

func getPostgresClient(conn string, tracing opentracing.Tracer) *psql.Client {
	client, err := psql.NewPsqlWithTracingConnection(conn, tracing)
	if err != nil {
		panic(err)
	}

	return client
}
