package usecase

import "pheet-fiber-backend/middleware/repository"


type ImiddlewareUsecase interface {

}

type middlewareUsecase struct {
	middleRepo repository.ImiddlewareRepository
}

func NewMiddlewareUsecase(middleRepo repository.ImiddlewareRepository) ImiddlewareUsecase {
	return middlewareUsecase{middleRepo: middleRepo}
}