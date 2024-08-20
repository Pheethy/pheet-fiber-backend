package users

import (
	"context"
	"pheet-fiber-backend/models"

	"github.com/gofrs/uuid"
)

type IUsersRepository interface {
	InsertUser(ctx context.Context, userReq *models.User, isAdmin bool) (*models.UserPassport, error)
	UpsertOAuth(ctx context.Context, req *models.OAuth) error
	FetchOneUserById(ctx context.Context, userId *uuid.UUID) (*models.User, error)
	FetchOneUserByEmail(ctx context.Context, email string) (*models.User, error)
	FetchOAuthByRefreshToken(ctx context.Context, refreshToken string) (*models.OAuth, error)
	DeleteOAuth(ctx context.Context, oauthId *uuid.UUID) error
}

