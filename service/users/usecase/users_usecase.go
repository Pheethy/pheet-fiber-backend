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
	cfg       config.Iconfig
	usersRepo users.IUsersRepository
}

func NewUsersUsecase(cfg config.Iconfig, usersRepo users.IUsersRepository) users.IUsersUsecase {
	return &usersUsecase{
		cfg:       cfg,
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

func (u usersUsecase) InsertAdmin(userReq *models.UserRegisterReq) (*models.UserPassport, error) {
	// Hashing a password
	if err := userReq.BcryptHashing(); err != nil {
		return nil, err
	}

	return u.usersRepo.InsertUser(userReq, true)
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
		Id:     user.Id,
		RoleId: user.RoleId,
	})
	if err != nil {
		return nil, fmt.Errorf("")
	}
	// Refresh Token
	refresh, err := _auth_service.NewAuthService(constants.Refresh, u.cfg.Jwt(), &models.UserClaims{
		Id:     user.Id,
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

func (u usersUsecase) FetchUserProfile(ctx context.Context, userId string) (*models.Users, error) {
	return u.usersRepo.FetchUserProfile(ctx, userId)
}

func (u usersUsecase) RefreshPassport(ctx context.Context, req *models.UserRefreshCredential) (*models.UserPassport, error) {

	claims, err := _auth_service.ParseToken(u.cfg.Jwt(), req.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("parse token failed: %v", err)
	}

	oAuth, err := u.usersRepo.FetchOneOauth(ctx, req.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("FetchOneOauth Failed: %v", err)
	}

	user, err := u.usersRepo.FetchUserProfile(ctx, oAuth.UserId)
	if err != nil {
		return nil, fmt.Errorf("fetch user profile failed: %v", err)
	}

	newClaims := &models.UserClaims{
		Id:     user.Id,
		RoleId: user.RoleId,
	}

	accessToken, err := _auth_service.NewAuthService(constants.Access, u.cfg.Jwt(), newClaims)
	if err != nil {
		return nil, fmt.Errorf("new claims failed: %v", err)
	}

	refreshToken := _auth_service.RepeatToken(u.cfg.Jwt(), newClaims, claims.ExpiresAt.Unix())

	passport := &models.UserPassport{
		User:  user,
		Token: &models.UserToken{
			Id:           oAuth.Id.String(),
			AccessToken:  accessToken.SignToken(),
			RefreshToken: refreshToken,
		},
	}

	if err := u.usersRepo.UpdateOauth(ctx, passport.Token); err != nil {
		return nil, fmt.Errorf("update oauth failed: %v", err)
	}

	return passport, nil
}

func (u usersUsecase) DeleteOauth(ctx context.Context, oId string) error {
	return u.usersRepo.DeleteOauth(ctx, oId)
}