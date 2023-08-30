package users

import "pheet-fiber-backend/models"

type IUsersUsecase interface {
	InsertCustomer(userReq *models.UserRegisterReq) (*models.UserPassport, error)
}