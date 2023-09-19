package usecase

import "pheet-fiber-backend/middleware/repository"


type ImiddlewareUsecase interface {
	FindAccessToken(userId, accessToken string) bool
}

type middlewareUsecase struct {
	middleRepo repository.ImiddlewareRepository
}

func NewMiddlewareUsecase(middleRepo repository.ImiddlewareRepository) ImiddlewareUsecase {
	return middlewareUsecase{middleRepo: middleRepo}
}

func (u middlewareUsecase) FindAccessToken(userId, accessToken string) bool {
	return u.middleRepo.FindAccessToken(userId, accessToken)
}