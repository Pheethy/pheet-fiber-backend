package users

import (
	"context"
	"pheet-fiber-backend/models"

	"github.com/gofrs/uuid"
)

type IUsersUsecase interface {
	FetchOneOauth(ctx context.Context, refreshToken string) (*models.OAuth, error)
	FetchUserProfile(ctx context.Context, userId *uuid.UUID) (*models.UserProfile, error)
	InsertUser(ctx context.Context, userReq *models.User) (*models.UserPassport, error)
	InsertAdmin(ctx context.Context, userReq *models.User) (*models.UserPassport, error)
	GetPassport(ctx context.Context, req *models.User) (*models.UserPassport, error)
	RefreshPassport(ctx context.Context, req *models.UserRefreshCredential) (*models.UserPassport, error)
	DeleteOAuth(ctx context.Context, oauthId *uuid.UUID) error
}

