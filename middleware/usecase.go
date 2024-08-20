package middleware

import (
	"context"
	"pheet-fiber-backend/models"

	"github.com/gofrs/uuid"
)

type IMiddlewareUsecase interface {
	FindAccessToken(ctx context.Context, userId *uuid.UUID, accessToken string) bool
	FetchRoles(ctx context.Context) ([]*models.Roles, error)
}
