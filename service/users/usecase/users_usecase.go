package usecase

import (
	"context"
	"fmt"
	"pheet-fiber-backend/config"
	"pheet-fiber-backend/models"
	"pheet-fiber-backend/service/users"

	"golang.org/x/crypto/bcrypt"
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

func (u usersUsecase) GetPassport(ctx context.Context, req *models.UserCredential) (*models.UserPassport, error) {
	//Fetch User
	user, err := u.usersRepo.FindOneUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	//Check Password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("password is invalid: %v", err)
	}

	//Cast user to userPassport
	userPass := &models.UserPassport{
		User: &models.Users{
			Id:       user.Id,
			Email:    user.Email,
			UserName: user.Username,
			RoleId:   user.RoleId,
		},
		Token: nil,
	}

	return userPass, nil
}