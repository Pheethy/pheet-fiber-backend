package repository

import (
	"pheet-fiber-backend/models"

	"github.com/jmoiron/sqlx"
)

type ImiddlewareRepository interface {
	FindAccessToken(userId, accessToken string) bool
	FindRole() ([]*models.Role, error)
}

type middlewareRepository struct {
	db *sqlx.DB
}

func NewMiddlewareRepository(db *sqlx.DB) ImiddlewareRepository {
	return middlewareRepository{db: db}
}

func (r middlewareRepository) FindAccessToken(userId, accessToken string) bool {
	sql := `
		SELECT
			(CASE WHEN COUNT(*) = 1 THEN TRUE ELSE FALSE END)
		FROM
			"oauth"
		WHERE
			"user_id" = $1
		AND
			"access_token" = $2;
	`
	var check bool
	if err := r.db.Get(&check, sql, userId, accessToken); err != nil {
		return false
	}

	return true
}

func (r middlewareRepository) FindRole() ([]*models.Role, error) {
	sql := `
		SELECT
			"id",
			"title"
		FROM
			"roles"
		ORDER BY
			"id" DESC;
	`
	var roles = make([]*models.Role, 0)
	if err := r.db.Select(&roles, sql); err != nil {
		return nil, err
	}

	return roles, nil
}
