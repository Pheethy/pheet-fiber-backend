package handlers

import (
	"pheet-fiber-backend/config"
	"pheet-fiber-backend/service/users"
)

type usersHandlers struct {
	cfg     config.Iconfig
	usersUs users.IUsersUsecase
}

func NewUsersHandler(cfg config.Iconfig, usersUs users.IUsersUsecase) users.IUsersHandlers {
	return usersHandlers{
		cfg:     cfg,
		usersUs: usersUs,
	}
}
