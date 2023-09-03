package users

import (
	"context"
	"pheet-fiber-backend/models"
)

type IUsersUsecase interface {
	InsertCustomer(userReq *models.UserRegisterReq) (*models.UserPassport, error)
	GetPassport(ctx context.Context, req *models.UserCredential) (*models.UserPassport, error)
}