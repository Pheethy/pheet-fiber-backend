package users

import (
	"context"
	"pheet-fiber-backend/models"
)

type IUsersRepository interface {
	InsertUser(userReq *models.UserRegisterReq, isAdmin bool) (*models.UserPassport, error)
	FindOneUserByEmail(ctx context.Context, email string) (*models.UserCredentialCheck, error)
}