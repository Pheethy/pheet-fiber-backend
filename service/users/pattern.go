package users

import "pheet-fiber-backend/models"

type IUsersPattern interface {
	Customer() (IUsersPattern, error)
	Admin() (IUsersPattern, error)
	Result() (*models.UserPassport, error)
}