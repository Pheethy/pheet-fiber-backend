package users

import (
	"context"
	"pheet-fiber-backend/models"
)

type IUsersRepository interface {
	InsertUser(userReq *models.UserRegisterReq, isAdmin bool) (*models.UserPassport, error)
	InsertOauth(ctx context.Context, req *models.UserPassport) error
	FindOneUserByEmail(ctx context.Context, email string) (*models.UserCredentialCheck, error)
	FetchUserProfile(ctx context.Context, id string) (*models.Users, error)
	FetchOneOauth(ctx context.Context, reToken string) (*models.Oauth, error)
	UpdateOauth(ctx context.Context, req *models.UserToken) error
	DeleteOauth(ctx context.Context, oId string) error
}