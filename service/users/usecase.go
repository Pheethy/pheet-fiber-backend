package users

import (
	"context"
	"pheet-fiber-backend/models"
)

type IUsersUsecase interface {
	InsertAdmin(userReq *models.UserRegisterReq) (*models.UserPassport, error)
	InsertCustomer(userReq *models.UserRegisterReq) (*models.UserPassport, error)
	GetPassport(ctx context.Context, req *models.UserCredential) (*models.UserPassport, error)
	FetchUserProfile(ctx context.Context, userId string) (*models.Users, error)
	RefreshPassport(ctx context.Context, req *models.UserRefreshCredential) (*models.UserPassport, error)
	DeleteOauth(ctx context.Context, oId string) error
}