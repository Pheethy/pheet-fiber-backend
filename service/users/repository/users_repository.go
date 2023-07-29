package repository

import (
	"pheet-fiber-backend/service/users"

	"github.com/jmoiron/sqlx"
)

type usersRepository struct {
	psqlDB *sqlx.DB
}

func NewUsersRepository(psqlDB *sqlx.DB) users.IUsersRepository {
	return usersRepository{
		psqlDB: psqlDB,
	}
}