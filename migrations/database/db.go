package database

import (
	"context"
	"log"
	"pheet-fiber-backend/config"

	"github.com/BlackMocca/sqlx"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	PGX = "pgx"
	SQLX = "sqlx"
)

func DBConnect(ctx context.Context, cfg config.IDbConfig) *sqlx.DB {
	//connect
	db, err := sqlx.ConnectContext(ctx, PGX, cfg.Url())
	if err != nil {
		log.Fatalf("connect to database failed: %v", err)
	}
	
	db.DB.SetMaxOpenConns(cfg.MaxConns())

	return db
}