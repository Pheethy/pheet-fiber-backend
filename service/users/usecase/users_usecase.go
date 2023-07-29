package usecase

import (
	"pheet-fiber-backend/config"
	"pheet-fiber-backend/service/users"
)

type usersUsecase struct {
	cfg config.Iconfig
	usersRepo users.IUsersRepository
}

func NewUsersUsecase(cfg config.Iconfig, usersRepo users.IUsersRepository) users.IUsersUsecase {
	return usersUsecase{
		cfg: cfg,
		usersRepo: usersRepo,
	}
}