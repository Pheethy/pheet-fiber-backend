package users

import "pheet-fiber-backend/models"

type IUsersRepository interface {
	InsertUser(userReq *models.UserRegisterReq, isAdmin bool) (*models.UserPassport, error)
}