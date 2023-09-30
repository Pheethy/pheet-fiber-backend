package usecase

import (
	"pheet-fiber-backend/middleware/repository"
	"pheet-fiber-backend/models"
)


type ImiddlewareUsecase interface {
	FindAccessToken(userId, accessToken string) bool
	FindRole() ([]*models.Role, error)
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

func (u middlewareUsecase) FindRole() ([]*models.Role, error) {
	return u.middleRepo.FindRole()
}