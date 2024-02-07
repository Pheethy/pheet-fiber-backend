package database

import (
	"context"
	"pheet-fiber-backend/config"

	"github.com/Pheethy/psql"
	"github.com/Pheethy/sqlx"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/opentracing/opentracing-go"

	_util_tracing "pheet-fiber-backend/service/utils/opentracing"
)

const (
	PGX  = "pgx"
	SQLX = "sqlx"
)

func DBConnect(ctx context.Context, cfg config.IDbConfig) *sqlx.DB {
	/* init tracing*/
	tracer, closer := _util_tracing.Init("flavorparser")
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	/* connect */
	psqlClient := getPostgresClient(cfg.Url(), tracer)

	db := psqlClient.GetClient()
	db.DB.SetMaxOpenConns(cfg.MaxConns())

	return db
}

func getPostgresClient(conn string, tracing opentracing.Tracer) *psql.Client {
	client, err := psql.NewPsqlWithTracingConnection(conn, tracing)
	if err != nil {
		panic(err)
	}

	return client
}
