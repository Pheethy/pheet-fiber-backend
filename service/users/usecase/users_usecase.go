package usecase

import (
	"context"
	"fmt"
	"pheet-fiber-backend/config"
	"pheet-fiber-backend/constants"
	"pheet-fiber-backend/models"
	"pheet-fiber-backend/service/users"

	_auth_service "pheet-fiber-backend/auth/service"

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

	//Access Token
	access, err := _auth_service.NewAuthService(constants.Access, u.cfg.Jwt(), &models.UserClaims{
		Id: user.Id,
		RoleId: user.RoleId,
	})
	if err != nil {
		return nil, fmt.Errorf("")
	}
	// Refresh Token
	refresh, err := _auth_service.NewAuthService(constants.Refresh, u.cfg.Jwt(), &models.UserClaims{
		Id: user.Id,
		RoleId: user.RoleId,
	})
	if err != nil {
		return nil, fmt.Errorf("")
	}

	//Cast user to userPassport
	userPass := &models.UserPassport{
		User: &models.Users{
			Id:       user.Id,
			Email:    user.Email,
			UserName: user.Username,
			RoleId:   user.RoleId,
		},
		Token: &models.UserToken{
			AccessToken:  access.SignToken(),
			RefreshToken: refresh.SignToken(),
		},
	}

	if err := u.usersRepo.InsertOauth(ctx, userPass); err != nil {
		return nil, err
	}
	
	return userPass, nil
}