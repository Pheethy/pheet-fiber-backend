package usecase

import (
	"context"
	"pheet-fiber-backend/middleware"
	"pheet-fiber-backend/models"

	"github.com/gofrs/uuid"
)

type middlewareUsecase struct {
	middleRepo middleware.IMiddlewareRepository
}

func NewMiddlewareUsecase(middleRepo middleware.IMiddlewareRepository) middleware.IMiddlewareUsecase {
	return middlewareUsecase{
		middleRepo: middleRepo,
	}
}

func (m middlewareUsecase) FindAccessToken(ctx context.Context, userId *uuid.UUID, accessToken string) bool {
	return m.middleRepo.FindAccessToken(ctx, userId, accessToken)
}

func (m middlewareUsecase) FetchRoles(ctx context.Context) ([]*models.Roles, error) {
	return m.middleRepo.FetchRoles(ctx)
}

