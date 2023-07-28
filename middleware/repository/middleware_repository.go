package repository

import "github.com/jmoiron/sqlx"

type ImiddlewareRepository interface {

}

type middlewareRepository struct {
	db *sqlx.DB
}

func NewMiddlewareRepository(db *sqlx.DB) ImiddlewareRepository {
	return middlewareRepository{db: db}
} 