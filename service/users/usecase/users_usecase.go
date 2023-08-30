package usecase

import (
	"pheet-fiber-backend/config"
	"pheet-fiber-backend/models"
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

func (u usersUsecase) InsertCustomer(userReq *models.UserRegisterReq) (*models.UserPassport, error) {
	// Hashing a password
	if err := userReq.BcryptHashing(); err != nil {
		return nil, err
	}
	
	return u.usersRepo.InsertUser(userReq, false)
}