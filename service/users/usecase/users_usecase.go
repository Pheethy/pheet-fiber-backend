package usecase

import (
	"context"
	"errors"
	"pheet-fiber-backend/auth"
	"pheet-fiber-backend/config"
	"pheet-fiber-backend/constants"
	"pheet-fiber-backend/models"
	"pheet-fiber-backend/polymor"
	"pheet-fiber-backend/service/users"

	"github.com/gofrs/uuid"
	"github.com/opentracing/opentracing-go"
)

type usersUsecase struct {
	cfg       config.Iconfig
	usersRepo users.IUsersRepository
}

func NewUsersUsecase(cfg config.Iconfig, usersRepo users.IUsersRepository) users.IUsersUsecase {
	return usersUsecase{
		cfg:       cfg,
		usersRepo: usersRepo,
	}
}

func (u usersUsecase) InsertUser(ctx context.Context, userReq *models.User) (*models.UserPassport, error) {
	/* Hashing Password */
	if err := userReq.BcryptHashing(); err != nil {
		return nil, err
	}

	return u.usersRepo.InsertUser(ctx, userReq, false)
}

func (u usersUsecase) InsertAdmin(ctx context.Context, userReq *models.User) (*models.UserPassport, error) {
	/* Hashing Password */
	if err := userReq.BcryptHashing(); err != nil {
		return nil, err
	}

	return u.usersRepo.InsertUser(ctx, userReq, true)
}

func (u usersUsecase) FetchUserProfile(ctx context.Context, userId *uuid.UUID) (*models.UserProfile, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FetchUserProfile")
	defer span.Finish()

	user, err := u.usersRepo.FetchOneUserById(ctx, userId)
	if err != nil {
		return nil, err
	}
	userProfile := &models.UserProfile{
		Id:        user.Id,
		Username:  user.Username,
		Email:     user.Email,
		RoleId:    user.RoleId,
		RoleTitle: user.RoleTitle,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	return userProfile, nil
}

func (u usersUsecase) GetPassport(ctx context.Context, req *models.User) (*models.UserPassport, error) {
	/* Find User */
	user, err := u.usersRepo.FetchOneUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	/* Check Password */
	if ok := user.ComparePassword(req); !ok {
		return nil, errors.New(constants.ERROR_PASSWORD_NOT_MATCH)
	}
	/* Sign Token */
	accessToken, err := auth.NewAuth(constants.Access, u.cfg.Jwt(), &models.UserClaims{
		Id:     user.Id,
		RoleId: user.RoleId,
	})
	if err != nil {
		return nil, err
	}
	refreshToken, err := auth.NewAuth(constants.Refresh, u.cfg.Jwt(), &models.UserClaims{
		Id:     user.Id,
		RoleId: user.RoleId,
	})
	if err != nil {
		return nil, err
	}
	/* Prepare OAuth */
	reqOAuth := new(models.OAuth)
	polymor.SetDefult(reqOAuth)
	/* Create Instant UserPassport */
	userPass := &models.UserPassport{
		User: &models.User{
			Id:        user.Id,
			Username:  user.Username,
			Email:     user.Email,
			RoleId:    user.RoleId,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		Token: &models.UserToken{
			Id:           reqOAuth.Id,
			AccessToken:  accessToken.SignToken(),
			RefreshToken: refreshToken.SignToken(),
		},
	}
	/* Stamp OAuth */
	reqOAuth.SetToken(userPass)
	if err := u.usersRepo.UpsertOAuth(ctx, reqOAuth); err != nil {
		return nil, err
	}

	return userPass, nil
}

func (u usersUsecase) RefreshPassport(ctx context.Context, req *models.UserRefreshCredential) (*models.UserPassport, error) {
	/* Parse Token */
	claims, err := auth.ParseToken(u.cfg.Jwt(), req.RefreshToken)
	if err != nil {
		return nil, err
	}
	/* Find Oauth*/
	oauth, err := u.usersRepo.FetchOAuthByRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}
	/* Find User */
	user, err := u.usersRepo.FetchOneUserById(ctx, oauth.UserId)
	if err != nil {
		return nil, err
	}
	/* Create New Claims */
	newClaims := &models.UserClaims{
		Id:     user.Id,
		RoleId: user.RoleId,
	}
	/* Sign Token */
	accessToken, err := auth.NewAuth(constants.Access, u.cfg.Jwt(), newClaims)
	if err != nil {
		return nil, err
	}
	refreshToken := auth.RepeatClaims(u.cfg.Jwt(), newClaims, int(claims.ExpiresAt.Unix()))
	/* Create UserPassport */
	userPass := &models.UserPassport{
		User: &models.User{
			Id:        user.Id,
			Username:  user.Username,
			Email:     user.Email,
			RoleId:    user.RoleId,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		Token: &models.UserToken{
			Id:           oauth.Id,
			AccessToken:  accessToken.SignToken(),
			RefreshToken: refreshToken,
		},
	}
	/* Update Oauth */
	oauth.SetToken(userPass)
	if err := u.usersRepo.UpsertOAuth(ctx, oauth); err != nil {
		return nil, err
	}

	return userPass, nil
}

func (u usersUsecase) FetchOneOauth(ctx context.Context, refreshToken string) (*models.OAuth, error) {
	return u.usersRepo.FetchOAuthByRefreshToken(ctx, refreshToken)
}

func (u usersUsecase) DeleteOAuth(ctx context.Context, oauthId *uuid.UUID) error {
	return u.usersRepo.DeleteOAuth(ctx, oauthId)
}
